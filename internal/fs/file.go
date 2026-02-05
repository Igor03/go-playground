package fs

import (
	"bufio"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func ReadFile(filePath string) ([]string, error) {
	f, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}

	defer f.Close()

	s := bufio.NewScanner(f)

	o := []string{}

	for s.Scan() {
		l := s.Text()
		o = append(o, l)

	}
	return o, nil
}

func CreateHashForFile(filePath string) (string, error) {
	f, err := os.Open(filePath)

	if err != nil {
		return "", err
	}

	defer f.Close()

	hasher := sha256.New()

	if _, err := io.Copy(hasher, f); err != nil {
		return "", err
	}

	h := hex.EncodeToString(hasher.Sum(nil))

	return h, nil

}

func CreateDummyFile(size int64, fileType string, fileName string, outputPath string) error {
	f, err := os.Create(fmt.Sprintf("%s/%s.%s", outputPath, fileName, fileType))
	tSize := size * 1024 * 1024 // This is the target size in bytes

	if err != nil {
		return err
	}

	defer f.Close()

	b := make([]byte, 5*1024*1024) // 5 MiB chunks
	var written int64

	for written < tSize {
		tWrite := int64(len(b))

		if remaining := tSize - written; remaining < tWrite {
			tWrite = remaining
		}

		// Make the values in the buffer random
		rand.Read(b[:tWrite])

		f.Write(b[:tWrite])

		written += tWrite
	}

	return f.Sync()
}
