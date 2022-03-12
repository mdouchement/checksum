package main

import (
	"bytes"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"hash/crc32"
	"io"
	"os"
	"strings"

	"github.com/muesli/coral"
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
)

var (
	version  = "dev"
	revision = "none"
	date     = "unknown"
)

type controller struct {
	supported []string
	algs      []string
	append    string
	verify    string

	mhashes map[string]hash.Hash
}

func main() {
	c := &controller{
		supported: []string{"crc32", "md5", "sha1", "sha256", "sha512", "blake2b", "blake2b512"},
		mhashes:   make(map[string]hash.Hash, 0),
	}

	cmd := &coral.Command{
		Use:          "checksum file",
		Short:        "File checksum",
		Long:         "File checksum",
		SilenceUsage: true,
		Version:      fmt.Sprintf("%s - build %.7s @ %s", version, revision, date),
		Args:         coral.ExactArgs(1),
		RunE: func(_ *coral.Command, args []string) error {
			filename := strings.TrimSpace(args[0])

			//

			if c.verify != "" {
				return c.validate(filename)
			}

			//

			err := c.compute(filename)
			if err != nil {
				return errors.Wrap(err, "compute")
			}

			if c.append != "" {
				return c.writeToFile(filename)
			}

			// STDOUT
			fmt.Println("Checksums:")
			for _, alg := range c.supported {
				if h, ok := c.mhashes[alg]; ok {
					fmt.Printf("%12s: %x\n", alg, h.Sum(nil))
				}
			}

			return nil
		},
	}
	cmd.Flags().StringSliceVarP(&c.algs, "algs", "", c.supported, `List of used hash algorithm (e.g. --algs="md5,sha1" --algs="sha256")`)
	cmd.Flags().StringVarP(&c.append, "append-to", "", "", "File to append checksums to")
	cmd.Flags().StringVarP(&c.verify, "verify", "", "", "Verify checksum of the file")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func (c *controller) validate(filename string) error {
	err := c.compute(filename)
	if err != nil {
		return errors.Wrap(err, "compute")
	}

	checksum := strings.TrimSpace(c.verify)
	for _, alg := range c.supported {
		if h, ok := c.mhashes[alg]; ok {
			hex.EncodeToString(h.Sum(nil))
			if checksum == hex.EncodeToString(h.Sum(nil)) {
				fmt.Println("Validated with", alg)
				return nil
			}
		}
	}

	return errors.New("invalid checksum")
}

func (c *controller) compute(filename string) error {
	var err error
	hashes := []io.Writer{}

	for _, alg := range c.algs {
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
		c.mhashes[alg] = h
	}

	//

	f, err := os.Open(filename)
	if err != nil {
		return errors.Wrap(err, "file")
	}
	defer f.Close()

	//

	w := io.MultiWriter(hashes...)
	_, err = io.Copy(w, f)
	return errors.Wrap(err, "checksum")
}

func (c *controller) writeToFile(filename string) error {
	b, err := os.ReadFile(c.append)
	if err != nil {
		b = []byte{}
	}
	buf := bytes.NewBuffer(b)

	for _, alg := range c.supported {
		if h, ok := c.mhashes[alg]; ok {
			buf.WriteString(fmt.Sprintf("%x  %s\n", h.Sum(nil), filename))
		}
	}

	return os.WriteFile(c.append, buf.Bytes(), 0644)
}
