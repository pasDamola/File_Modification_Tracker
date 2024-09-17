package services

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"strconv"
	"sync"
	"time"

	"github.com/pasDamola/file-tracker/config"
	"github.com/pasDamola/file-tracker/internal/adapters/osquery"
	"github.com/pasDamola/file-tracker/internal/core/domain"
	"github.com/pasDamola/file-tracker/internal/core/ports"
)

type Daemon struct {
	WorkerQueue     chan string
	cfg             config.Config
	client          *osquery.OsqueryAdapter
	fileService     ports.FileService
	wg              sync.WaitGroup
	isWorkerRunning bool
	isTimerRunning  bool
	mu              sync.Mutex
	stopCh          chan struct{}
}

func NewDaemon(cfg config.Config, queryClient *osquery.OsqueryAdapter, fileService ports.FileService) *Daemon {
	return &Daemon{
		WorkerQueue: make(chan string),
		cfg:         cfg,
		client:      queryClient,
		fileService: fileService,
		stopCh:      make(chan struct{}),
	}

}

func (d *Daemon) Start() {
	d.isWorkerRunning = true
	d.isTimerRunning = true

	d.wg.Add(2)
	go func() {
		defer d.wg.Done()
		d.workerThread()
	}()

	go func() {
		defer d.wg.Done()
		d.timerThread()
	}()
}

func (d *Daemon) Stop() {
	close(d.stopCh)

	d.wg.Wait()

	close(d.WorkerQueue)

	d.isTimerRunning = false
	d.isWorkerRunning = false
}

func (d *Daemon) QueueCommand(command string) {
	d.WorkerQueue <- command
}

func (d *Daemon) IsWorkerThreadRunning() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.isWorkerRunning
}

func (d *Daemon) IsTimerThreadRunning() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return d.isTimerRunning
}

func (d *Daemon) workerThread() {
	for {
		select {
		case command := <-d.WorkerQueue:
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
		case <-d.stopCh:
			return
		}
	}
}

func (d *Daemon) timerThread() {
	ticker := time.NewTicker(time.Duration(d.cfg.Frequency) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			log.Println("Fetching file modification stats...")

			query := "SELECT path, ctime FROM file WHERE directory = '" + d.cfg.Directory + "';"
			results, err := d.client.Client.Query(query)
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
				formattedTime := t.Format("06/01/02 03:04 PM")
				logText := fmt.Sprintf("File modified: %s at time: %v", result["path"], formattedTime)

				modification := domain.FileModification{Path: result["path"], Timestamp: formattedTime}
				d.fileService.LogModification(modification)

				log.Println(logText)
			}
		case <-d.stopCh:
			return
		}
	}
}
