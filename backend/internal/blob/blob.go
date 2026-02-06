package blob

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/klauspost/compress/zstd"
)

type Store struct {
	Root string
}

func New(root string) *Store {
	return &Store{Root: root}
}

func (s *Store) WriteCompressed(relPath string, r io.Reader) (string, int64, error) {
	fullPath := filepath.Join(s.Root, relPath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		return "", 0, err
	}
	file, err := os.Create(fullPath)
	if err != nil {
		return "", 0, err
	}
	defer file.Close()

	encoder, err := zstd.NewWriter(file)
	if err != nil {
		return "", 0, err
	}
	defer encoder.Close()

	n, err := io.Copy(encoder, r)
	return fullPath, n, err
}

func (s *Store) ReadDecompressed(relPath string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.Root, relPath)
	file, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	decoder, err := zstd.NewReader(file)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	return &combinedCloser{reader: decoder, file: file}, nil
}

func (s *Store) HashName(input string) string {
	sum := sha256.Sum256([]byte(input))
	return hex.EncodeToString(sum[:])
}

type combinedCloser struct {
	reader *zstd.Decoder
	file   *os.File
}

func (c *combinedCloser) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}

func (c *combinedCloser) Close() error {
	c.reader.Close()
	return c.file.Close()
}

func (s *Store) RelativePath(category, name string) string {
	return filepath.ToSlash(filepath.Join(category, fmt.Sprintf("%s.zst", s.HashName(name))))
}
