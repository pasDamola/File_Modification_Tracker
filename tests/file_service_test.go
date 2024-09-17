package tests

import (
	"testing"

	"github.com/pasDamola/file-tracker/internal/core/domain"
)

func TestFileService(t *testing.T) {
	fileService := &MockFileService{}

	modification := domain.FileModification{Path: "/Users/oyincode/Desktop/test_file/hello_world.txt", Timestamp: "06/01/02 03:04 PM"}
	err := fileService.LogModification(modification)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	logs := fileService.GetLogs()
	if len(logs) != 1 {
		t.Fatalf("Expected 1 log entry, got %d", len(logs))
	}

	if logs[0] != modification {
		t.Fatalf("Expected log entry to be %v, got %v", modification, logs[0])
	}
}
