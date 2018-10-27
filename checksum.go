package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"strings"

	"golang.org/x/crypto/blake2b"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	supported = []string{"crc32", "md5", "sha1", "sha256", "sha512", "blake2b", "blake2b512"}
	algs      []string
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	c := &cobra.Command{
		Use:     "checksum file",
		Short:   "File checksum",
		Long:    "File checksum",
		Version: fmt.Sprintf("%s - build %.7s @ %s", version, commit, date),
		Args:    cobra.ExactArgs(1),
		RunE:    action,
	}
	c.Flags().StringSliceVarP(&algs, "algs", "", supported, `List of used hash algorithm (e.g. --algs="md5,sha1" --algs="sha256")`)

	if err := c.Execute(); err != nil {
		fmt.Println(err)
	}
}

func action(c *cobra.Command, args []string) (err error) {
	hashes := []io.Writer{}
	mhashes := map[string]hash.Hash{}
	for _, alg := range algs {
		var h hash.Hash
		switch alg {
		case "crc32":
			h = crc32.New(crc32.IEEETable)
		case "md5":
			h = md5.New()
		case "sha1":
			h = sha1.New()
		case "sha256":
			h = sha256.New()
		case "sha512":
			h = sha512.New()
		case "blake2b":
			h, err = blake2b.New256(nil)
			if err != nil {
				return errors.Wrap(err, "blake2b")
			}
		case "blake2b512":
			h, err = blake2b.New512(nil)
			if err != nil {
				return errors.Wrap(err, "blake2b")
			}
		default:
			return errors.Errorf("Unsuported algorithm: %s", alg)
		}
		hashes = append(hashes, h)
		mhashes[alg] = h
	}

	f, err := os.Open(strings.TrimSpace(args[0]))
	if err != nil {
		return errors.Wrap(err, "file")
	}
	defer f.Close()

	w := io.MultiWriter(hashes...)
	if _, err := io.Copy(w, f); err != nil {
		return errors.Wrap(err, "checksum")
	}

	fmt.Println("Checksums:")
	for _, alg := range supported {
		if h, ok := mhashes[alg]; ok {
			fmt.Printf("%12s: %x\n", alg, h.Sum(nil))
		}
	}

	return nil
}
