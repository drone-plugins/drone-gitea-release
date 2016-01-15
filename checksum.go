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
	"io/ioutil"
	"strconv"
)

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
	default:
		return "", fmt.Errorf("hashing method %s is not supported", method)
	}
}
