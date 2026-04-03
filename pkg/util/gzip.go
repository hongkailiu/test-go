// Package util has utility functions
package util

import (
	"archive/tar"
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

	if !strings.HasSuffix(cleanPath, ".gz") {
		retErr = os.WriteFile(cleanPath, data, 0600)

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
	if retErr != nil {
		return fmt.Errorf("error compressing data to file %s: %w", cleanPath, retErr)
	}

	return nil
}

func AddFileToTar(tw *tar.Writer, trim string, filePath string) (errRet error) {
	file, err := os.Open(filePath)
	if err != nil {
		errRet = err
		return
	}
	defer func() {
		if err := file.Close(); err != nil && errRet == nil {
			errRet = err
		}
	}()

	stat, err := file.Stat()
	if err != nil {
		errRet = err
		return
	}

	// Create the tar header based on the file system stats
	header, err := tar.FileInfoHeader(stat, "")
	if err != nil {
		errRet = err
		return
	}
	// Ensure the header has the correct base filename
	//header.Name = filepath.Base(filePath)
	header.Name = strings.TrimPrefix(filePath, trim)

	// Write the header to the tar stream
	if err := tw.WriteHeader(header); err != nil {
		errRet = err
		return
	}

	// Stream the file contents directly into the tar writer
	_, err = io.Copy(tw, file)
	errRet = err
	return
}
