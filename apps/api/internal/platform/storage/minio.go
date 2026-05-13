package storage

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/moodora/moodora/apps/api/internal/platform/config"
)

func Connect(cfg config.StorageConfig) (*minio.Client, error) {
	endpoint, useSSL, err := normalizeEndpoint(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Region: cfg.Region,
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create object storage client: %w", err)
	}

	return client, nil
}

func normalizeEndpoint(rawEndpoint string) (string, bool, error) {
	endpoint := strings.TrimSpace(rawEndpoint)
	if endpoint == "" {
		return "", false, fmt.Errorf("S3_ENDPOINT is required")
	}

	parsed, err := url.Parse(endpoint)
	if err != nil {
		return "", false, fmt.Errorf("parse S3_ENDPOINT: %w", err)
	}

	if parsed.Scheme == "" {
		return endpoint, false, nil
	}

	switch parsed.Scheme {
	case "http":
		return parsed.Host, false, nil
	case "https":
		return parsed.Host, true, nil
	default:
		return "", false, fmt.Errorf("S3_ENDPOINT must use http or https")
	}
}
