package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/pasDamola/file-tracker/daemon"
)

func TestIsReadOnlyCommand(t *testing.T) {
	tests := []struct {
		command  string
		expected bool
	}{
		{"ls -l", true},
		{"cat file.txt", true},
		{"rm -rf /", false},
		{"mv file.txt newfile.txt", false},
		{"echo Hello", true},
		{"> output.txt", false},
	}

	for _, tt := range tests {
		t.Run(tt.command, func(t *testing.T) {
			result := isReadOnlyCommand(tt.command)
			if result != tt.expected {
				t.Errorf("isReadOnlyCommand(%q) = %v; want %v", tt.command, result, tt.expected)
			}
		})
	}
}

func TestLogFormatting(t *testing.T) {
	timestamp := int64(1726486988)
	expectedFormat := "24/09/16 12:43 PM"

	tm := time.Unix(timestamp, 0)
	formattedTime := tm.Format("06/01/02 03:04 PM")

	if formattedTime != expectedFormat {
		t.Errorf("Expected formatted time %s but got %s", expectedFormat, formattedTime)
	}
}

func TestHealthCheckEndpoint(t *testing.T) {
	d := &daemon.Daemon{
		WorkerQueue:      make(chan string),
		Logs:             []string{},
		WorkerThreadDone: make(chan struct{}),
		TimerThreadDone:  make(chan struct{}),
	}

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		healthCheckHandler(w, r, d)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}
}

func TestLogsHandler(t *testing.T) {
	d := &daemon.Daemon{
		WorkerQueue: make(chan string),
		Logs:        []string{"Log entry 1", "Log entry 2"},
	}

	req, err := http.NewRequest("GET", "/logs", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logsHandler(w, r, d)
	})

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var logsResponse []string
	json.Unmarshal(rr.Body.Bytes(), &logsResponse)

	if len(logsResponse) != len(d.Logs) {
		t.Errorf("Expected %d logs but got %d", len(d.Logs), len(logsResponse))
	}
}
