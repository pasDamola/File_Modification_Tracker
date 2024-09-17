package services

import (
	"github.com/pasDamola/file-tracker/internal/core/domain"
	"github.com/pasDamola/file-tracker/internal/core/ports"
)

type fileService struct {
	logs []domain.FileModification // In-memory store for simplicity; replace with a persistent store in production.
}

func NewFileService() ports.FileService {
	return &fileService{
		logs: []domain.FileModification{},
	}
}

func (fs *fileService) LogModification(modification domain.FileModification) error {
	fs.logs = append(fs.logs, modification)
	return nil
}

func (fs *fileService) GetLogs() []domain.FileModification {
	return fs.logs
}
