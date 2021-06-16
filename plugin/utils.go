// Copyright (c) 2021, the Drone Plugins project authors.
// Please see the AUTHORS file for details. All rights reserved.
// Use of this source code is governed by an Apache 2.0 license that can be
// found in the LICENSE file.

package plugin

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash/adler32"
	"hash/crc32"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

var (
	fileExistsValues = map[string]bool{
		"overwrite": true,
		"fail":      true,
		"skip":      true,
	}
)

func readStringOrFile(input string) (string, error) {
	// Check if input is a file path
	if _, err := os.Stat(input); err != nil && os.IsNotExist(err) {
		// No file found => use input as result
		return input, nil
	} else if err != nil {
		return "", err
	}
	result, err := ioutil.ReadFile(input)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

func checksum(r io.Reader, method string) (string, error) {
	b, err := ioutil.ReadAll(r)

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
	}

	return "", fmt.Errorf("hashing method %s is not supported", method)
}

func writeChecksums(files, methods []string) ([]string, error) {
	checksums := make(map[string][]string)

	for _, method := range methods {
		for _, file := range files {
			handle, err := os.Open(file)

			if err != nil {
				return nil, fmt.Errorf("failed to read %s artifact: %s", file, err)
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
