package compress

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GZip implements a compress which based gzip
type GZip struct {
	additionBinaries []string
}

// NewGZip creates an instance of GZip
// additionBinaries could be empty or nil
func NewGZip(additionBinaries []string) *GZip {
	return &GZip{additionBinaries: additionBinaries}
}

// make sure GZip implements the interface Compress
var _ Compress = &GZip{}

// ExtractFiles extracts files from a target compress file
func (c *GZip) ExtractFiles(sourceFile, targetName string) (err error) {
	if targetName == "" {
		err = errors.New("target filename is empty")
		return
	}

	var f *os.File
	var gzf *gzip.Reader
	if f, err = os.Open(sourceFile); err != nil {
		return
	}
	defer func() {
		_ = f.Close()
	}()

	if gzf, err = gzip.NewReader(f); err != nil {
		return
	}

	tarReader := tar.NewReader(gzf)
	var header *tar.Header
	var found bool
	for {
		if header, err = tarReader.Next(); err == io.EOF {
			err = nil
			break
		} else if err != nil {
			break
		}
		name := header.Name

		switch header.Typeflag {
		case tar.TypeReg:
			if err = extraFile(name, targetName, sourceFile, header, tarReader); err == nil {
				found = true
			} else {
				break
			}

			for i := range c.additionBinaries {
				addition := c.additionBinaries[i]
				if err = extraFile(addition, addition, sourceFile, header, tarReader); err != nil {
					return
				}
			}
		}
	}

	if err == nil && !found {
		err = fmt.Errorf("cannot found item '%s' from '%s'", targetName, sourceFile)
	}
	return
}

func extraFile(name, targetName, tarFile string, header *tar.Header, tarReader *tar.Reader) (err error) {
	if name != targetName && !strings.HasSuffix(name, "/"+targetName) {
		return
	}
	var targetFile *os.File
	if targetFile, err = os.OpenFile(fmt.Sprintf("%s/%s", filepath.Dir(tarFile), targetName),
		os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode)); err != nil {
		return
	}
	if _, err = io.Copy(targetFile, tarReader); err != nil {
		return
	}
	_ = targetFile.Close()
	return
}
