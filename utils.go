package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash/adler32"
	"hash/crc32"
	"io"
	"os"
	"strconv"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
)

var (
	fileExistsValues = map[string]bool{
		"overwrite": true,
		"fail":      true,
		"skip":      true,
	}
)

func checksum(r io.Reader, method string) (string, error) {
	b, err := io.ReadAll(r)

	if err != nil {
		return "", err
	}

	switch method {
	case "md5":
		return fmt.Sprintf("%x", md5.Sum(b)), nil
	case "sha1":
		return fmt.Sprintf("%x", sha1.Sum(b)), nil
	case "sha256":
		return fmt.Sprintf("%x", sha256.Sum256(b)), nil
	case "sha512":
		return fmt.Sprintf("%x", sha512.Sum512(b)), nil
	case "adler32":
		return strconv.FormatUint(uint64(adler32.Checksum(b)), 10), nil
	case "crc32":
		return strconv.FormatUint(uint64(crc32.ChecksumIEEE(b)), 10), nil
	case "blake2b":
		return fmt.Sprintf("%x", blake2b.Sum256(b)), nil
	case "blake2s":
		return fmt.Sprintf("%x", blake2s.Sum256(b)), nil
	}

	return "", fmt.Errorf("Hashing method %s is not supported", method)
}

func writeChecksums(files, methods []string) ([]string, error) {
	checksums := make(map[string][]string)

	for _, method := range methods {
		for _, file := range files {
			handle, err := os.Open(file)

			if err != nil {
				return nil, fmt.Errorf("Failed to read %s artifact: %s", file, err)
			}

			hash, err := checksum(handle, method)

			if err != nil {
				return nil, err
			}

			checksums[method] = append(checksums[method], hash, file)
		}
	}

	for method, results := range checksums {
		filename := method + "sum.txt"
		f, err := os.Create(filename)

		if err != nil {
			return nil, err
		}

		for i := 0; i < len(results); i += 2 {
			hash := results[i]
			file := results[i+1]

			if _, err := f.WriteString(fmt.Sprintf("%s  %s\n", hash, file)); err != nil {
				return nil, err
			}
		}

		files = append(files, filename)
	}

	return files, nil
}
