package gitignore

import (
	"errors"
	"github.com/gofunky/zipred"
	"os"
	"strings"
)

// gitIgnoreContent is the context that implements zipred.Zipred.
type gitIgnoreContent struct {
	// patterns to filter
	patterns []string
	// number of patterns found
	count int
}

// URL to download the archive from
func (c *gitIgnoreContent) URL() string {
	return archiveURL
}

// Predicate indicates if the given file should be read or not.
// It is to return a zero string to discard the file, otherwise the key name.
// If error is nonempty, the download is aborted and the error is passed on.
func (c *gitIgnoreContent) Predicate(fileInfo os.FileInfo) (key string, err error) {
	fileName := fileInfo.Name()
	if strings.HasSuffix(fileName, gitignoreSuffix) {
		alias := strings.ToLower(strings.TrimSuffix(fileName, gitignoreSuffix))
		for _, pat := range c.patterns {
			if strings.ToLower(pat) == alias {
				c.count++
				return alias, nil
			}
		}
	}
	return
}

// Done indicates if enough data has been read and the download can be aborted ahead of the EOF.
// isEOF is true if the end of the zip file has been reached.
// If error is nonempty, the download is aborted and the error is passed on.
func (c *gitIgnoreContent) Done(isEOF bool) (finish bool, err error) {
	if c.count == len(c.patterns) {
		return true, nil
	} else if isEOF {
		return true, errors.New("not all given gitignore patterns could be found")
	}
	return
}

// Get the given gitignore patterns.
func Get(patterns []string) (files map[string][]byte, err error) {
	context := &gitIgnoreContent{
		patterns: patterns,
	}
	return zipred.FilterZipContent(context)
}

// GetAll available gitignore patterns.
func GetAll() (files map[string][]byte, err error) {
	allPatterns, err := List()
	if err != nil {
		return nil, err
	}
	context := &gitIgnoreContent{
		patterns: allPatterns,
	}
	return zipred.FilterZipContent(context)
}
