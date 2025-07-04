package checksum

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func GenerateFileChecksum(filePath string) string {
	file, err := os.Open(filePath)
	if err != nil {
		panic(fmt.Errorf("failed to open file: %w", err))
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		panic(fmt.Errorf("failed to calculate checksum: %w", err))
	}

	checksum := hex.EncodeToString(hash.Sum(nil))
	return checksum
}
