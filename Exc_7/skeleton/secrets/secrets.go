package secrets

import (
	"fmt"
	"os"
	"strings"
)

// LoadSecretOrEnv tries to load a secret from a file defined by <key>_FILE environment variable.
// If that fails or is not set, it tries to load the value directly from <key> environment variable.
func LoadSecretOrEnv(key string) (string, error) {
	// 1. Try to read from file defined by <key>_FILE
	fileEnvKey := key + "_FILE"
	filePath := os.Getenv(fileEnvKey)
	if filePath != "" {
		content, err := os.ReadFile(filePath)
		if err != nil {
			return "", fmt.Errorf("failed to read secret file from %s: %w", filePath, err)
		}
		return strings.TrimSpace(string(content)), nil
	}

	// 2. Try to read directly from <key>
	val := os.Getenv(key)
	if val != "" {
		return val, nil
	}

	return "", fmt.Errorf("neither %s nor %s environment variables are set", fileEnvKey, key)
}
