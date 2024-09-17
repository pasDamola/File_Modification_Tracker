package http

import (
	"encoding/json"
	"net/http"
	"strings"
	"sync"

	"github.com/pasDamola/file-tracker/internal/core/ports"
	"github.com/pasDamola/file-tracker/internal/core/services"
)

var (
	mu sync.Mutex
)

func StartAPI(fileService ports.FileService, daemon *services.Daemon) {
	http.HandleFunc("/execute", func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Commands []string `json:"commands"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		for _, command := range req.Commands {
			if !isReadOnlyCommand(command) {
				http.Error(w, "Command not allowed", http.StatusForbidden)
				return
			}
			daemon.QueueCommand(command)
		}
		w.WriteHeader(http.StatusAccepted)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		if !daemon.IsWorkerThreadRunning() || !daemon.IsTimerThreadRunning() {
			http.Error(w, "Service is not healthy", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Service is running"))
	})

	http.HandleFunc("/logs", func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(fileService.GetLogs())
	})

	http.ListenAndServe(":8080", nil)
}

func isReadOnlyCommand(command string) bool {
	// Check if the command contains any write or delete operations
	return !strings.Contains(command, ">") &&
		!strings.Contains(command, ">>") &&
		!strings.Contains(command, "rm") &&
		!strings.Contains(command, "mv") &&
		!strings.Contains(command, "cp")
}
