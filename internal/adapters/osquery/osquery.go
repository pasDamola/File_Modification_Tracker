package osquery

import (
	"time"

	"github.com/osquery/osquery-go"
)

type OsqueryAdapter struct {
	Client *osquery.ExtensionManagerClient
}

func NewOsqueryAdapter(socketPath string) (*OsqueryAdapter, error) {
	client, err := osquery.NewClient(socketPath, 5*time.Second)
	if err != nil {
		return nil, err
	}
	return &OsqueryAdapter{Client: client}, nil
}
