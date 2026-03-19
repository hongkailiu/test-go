// Package util has utility functions
package util

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ReadFileMaybeGZIP wraps util.ReadBytesMaybeGZIP, returning the decompressed contents.
// if the file is gzipped, or otherwise the raw contents.
func ReadFileMaybeGZIP(path string) ([]byte, error) {
	cleanPath := filepath.Clean(path)

	b, err := os.ReadFile(cleanPath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", cleanPath, err)
	}

	return readBytesMaybeGZIP(b)
}

func readBytesMaybeGZIP(data []byte) ([]byte, error) {
	// check if data contains gzip header: http://www.zlib.org/rfc-gzip.html
	if !bytes.HasPrefix(data, []byte("\x1F\x8B")) {
		// go ahead and return the contents if not gzipped
		return data, nil
	}
	// otherwise decode
	gzipReader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("error creating gzip reader: %w", err)
	}

	ret, err := io.ReadAll(gzipReader)
	if err != nil {
		return nil, fmt.Errorf("error decompressing data: %w", err)
	}

	return ret, nil
}

// WriteBytesMaybeGZIP write data to a file.
// If the file has suffix .gz then the content is compressed.
func WriteBytesMaybeGZIP(file string, data []byte) (retErr error) {
	cleanPath := filepath.Clean(file)

	const perm os.FileMode = 0600 // Unix permission bits
	if !strings.HasSuffix(cleanPath, ".gz") {
		retErr = os.WriteFile(cleanPath, data, perm)

		return fmt.Errorf("error writing file %s: %w", cleanPath, retErr)
	}

	outFile, retErr := os.Create(cleanPath)
	if retErr != nil {
		return fmt.Errorf("error creating file %s: %w", cleanPath, retErr)
	}

	defer func() {
		err := outFile.Close()
		if retErr != nil && err != nil {
			retErr = err
		}
	}()

	gzWriter := gzip.NewWriter(outFile)

	defer func() {
		// defer is a stack (FILO).
		// Close the gzip writer BEFORE the file.
		// This flushes the compression buffers and writes the gzip footer.
		err := gzWriter.Close()
		if retErr != nil && err != nil {
			retErr = err
		}
	}()

	_, retErr = gzWriter.Write(data)

	return fmt.Errorf("error writing file %s: %w", cleanPath, retErr)
}
