package tests

import (
	"testing"
	"time"

	"github.com/pasDamola/file-tracker/config"
	"github.com/pasDamola/file-tracker/internal/adapters/osquery"
	"github.com/pasDamola/file-tracker/internal/core/domain"
	"github.com/pasDamola/file-tracker/internal/core/services"
)

type MockFileService struct {
	logs []domain.FileModification
}

func (m *MockFileService) LogModification(modification domain.FileModification) error {
	m.logs = append(m.logs, modification)
	return nil
}

func (m *MockFileService) GetLogs() []domain.FileModification {
	return m.logs
}

func TestDaemon(t *testing.T) {
	cfg := config.Config{
		Directory:  "/Users/oyincode/Desktop/test_file",
		Frequency:  1,
		SocketPath: "/Users/oyincode/.osquery/shell.em",
	}
	mockFileService := &MockFileService{}
	osqueryClient := &osquery.OsqueryAdapter{} // Mock or real implementation

	d := services.NewDaemon(cfg, osqueryClient, mockFileService)

	d.Start()
	time.Sleep(2 * time.Second)

	if !d.IsWorkerThreadRunning() {
		t.Errorf("Expected worker thread to be running")
	}

	if !d.IsTimerThreadRunning() {
		t.Errorf("Expected timer thread to be running")
	}

	d.Stop()

	if d.IsWorkerThreadRunning() {
		t.Errorf("Expected worker thread to be stopped")
	}

	if d.IsTimerThreadRunning() {
		t.Errorf("Expected timer thread to be stopped")
	}
}
