package config_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"yourmodule/internal/config"
)

func writeTempConfig(t *testing.T, content string) string {
	t.Helper()

	dir := t.TempDir()
	file := filepath.Join(dir, "config.yaml")
	err := os.WriteFile(file, []byte(content), 0644)
	require.NoError(t, err)

	return file
}

func TestLoad_ValidConfig(t *testing.T) {
	yaml := `
listen_addr: ":8080"
backends:
  - http://localhost:9001
  - http://localhost:9002
health_check_interval: 5s
health_check_timeout: 2s
`
	path := writeTempConfig(t, yaml)
	cfg, err := config.Load(path)

	require.NoError(t, err)
	assert.Equal(t, ":8080", cfg.ListenAddr)
	assert.Len(t, cfg.Backends, 2)
	assert.Equal(t, 5*time.Second, cfg.HealthCheckInterval)
	assert.Equal(t, 2*time.Second, cfg.HealthCheckTimeout)
}

func TestLoad_InvalidYAML(t *testing.T) {
	yaml := `listen_addr: ":8080" bad_yaml`
	path := writeTempConfig(t, yaml)

	_, err := config.Load(path)
	require.Error(t, err)
}

func TestLoad_MissingRequiredFields(t *testing.T) {
	yaml := `
listen_addr: ""
backends: []
health_check_interval: 0s
health_check_timeout: 0s
`
	path := writeTempConfig(t, yaml)

	_, err := config.Load(path)
	require.Error(t, err)
}
