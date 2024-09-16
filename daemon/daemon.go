package daemon

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"time"

	"github.com/osquery/osquery-go"
	"github.com/pasDamola/file-tracker/config"
)

type Daemon struct {
	WorkerQueue      chan string
	cfg              config.Config
	client           *osquery.ExtensionManagerClient
	WorkerThreadDone chan struct{}
	TimerThreadDone  chan struct{}
	Logs             []string
}

func NewDaemon(cfg config.Config) *Daemon {

	client, err := osquery.NewClient(cfg.SocketPath, 5*time.Second)
	if err != nil {
		log.Fatalf("Failed to create osquery client: %v", err)
	}

	return &Daemon{
		WorkerQueue: make(chan string),
		cfg:         cfg,
		client:      client,
	}
}

func (d *Daemon) Start() {
	go d.workerThread()
	go d.timerThread()
}

func (d *Daemon) workerThread() {
	defer close(d.WorkerThreadDone)
	for command := range d.WorkerQueue {
		log.Printf("Executing command: %s", command)

		cmd := exec.Command("sh", "-c", command)
		var out bytes.Buffer
		cmd.Stdout = &out
		var stderr bytes.Buffer
		cmd.Stderr = &stderr

		if err := cmd.Run(); err != nil {
			log.Printf("Error executing command: %s, stderr: %s", err, stderr.String())
			continue
		}

		log.Printf("Command output: %s", out.String())
	}
}

func (d *Daemon) timerThread() {
	defer close(d.TimerThreadDone)
	ticker := time.NewTicker(time.Duration(d.cfg.Frequency) * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		log.Println("Fetching file modification stats...")

		query := "SELECT path, ctime FROM file WHERE directory = '" + d.cfg.Directory + "';"

		results, err := d.client.Query(query)
		if err != nil {
			log.Printf("Error querying osquery: %v", err)
			continue
		}

		for _, result := range results.Response {
			ctimeUnix, err := strconv.ParseInt(result["ctime"], 10, 64)
			if err != nil {
				log.Printf("Error parsing ctime: %v", err)
				continue
			}

			t := time.Unix(ctimeUnix, 0)

			// Format time as "yy/mm/dd hh:mm AM/PM for easier understanding"
			formattedTime := t.Format("06/01/02 03:04 PM")
			logText := fmt.Sprintf("File modified: %s at time: %v", result["path"], formattedTime)
			d.Logs = append(d.Logs, logText)
			log.Printf("File modified: %s at time: %v", result["path"], formattedTime)
		}
	}
}
