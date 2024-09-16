package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/pasDamola/file-tracker/daemon"
)

var (
	mu sync.Mutex
)

func StartAPI(d *daemon.Daemon) {
	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		type executeRequest struct {
			Commands []string `json:"commands"`
		}
		var req executeRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		for _, command := range req.Commands {
			if !isReadOnlyCommand(command) {
				http.Error(w, "Command not allowed", http.StatusForbidden)
				return
			}
			d.WorkerQueue <- command
		}
		w.WriteHeader(http.StatusAccepted)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		healthCheckHandler(w, r, d)
	})

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		logsHandler(w, r, d)
	})

	http.ListenAndServe(":8080", nil)
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request, d *daemon.Daemon) {
	mu.Lock()
	defer mu.Unlock()

	select {
	case <-d.WorkerThreadDone:
		http.Error(w, "Worker thread not running", http.StatusInternalServerError)
		return
	default:
	}

	select {
	case <-d.TimerThreadDone:
		http.Error(w, "Timer thread not running", http.StatusInternalServerError)
		return
	default:
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Service is running"))
}

func logsHandler(w http.ResponseWriter, r *http.Request, d *daemon.Daemon) {
	mu.Lock()
	defer mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d.Logs)
}

func isReadOnlyCommand(command string) bool {
	// Check if the command contains any write or delete operations
	return !strings.Contains(command, ">") &&
		!strings.Contains(command, ">>") &&
		!strings.Contains(command, "rm") &&
		!strings.Contains(command, "mv") &&
		!strings.Contains(command, "cp")
}
