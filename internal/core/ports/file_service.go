package ports

import "github.com/pasDamola/file-tracker/internal/core/domain"

type FileService interface {
	LogModification(modification domain.FileModification) error
	GetLogs() []domain.FileModification
}
