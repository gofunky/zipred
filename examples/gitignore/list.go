package gitignore

import (
	"github.com/gofunky/zipred"
	"os"
	"strings"
)

// gitIgnoreList is the context that implements zipred.Zipred.
type gitIgnoreList struct{}

// URL to download the archive from
func (c *gitIgnoreList) URL() string {
	return archiveURL
}

// Predicate indicates if the given file should be read or not.
// It is to return a zero string to discard the file, otherwise the key name.
// If error is nonempty, the download is aborted and the error is passed on.
func (c *gitIgnoreList) Predicate(fileInfo os.FileInfo) (key string, err error) {
	fileName := fileInfo.Name()
	if strings.HasSuffix(fileName, gitignoreSuffix) {
		return strings.ToLower(strings.TrimSuffix(fileName, gitignoreSuffix)), nil
	}
	return
}

// Done indicates if enough data has been read and the download can be aborted ahead of the EOF.
// isEOF is true if the end of the zip file has been reached.
// If error is nonempty, the download is aborted and the error is passed on.
func (c *gitIgnoreList) Done(isEOF bool) (finish bool, err error) {
	return
}

// List all available gitignore patterns.
func List() (patterns []string, err error) {
	return zipred.FilterFileInfo(&gitIgnoreList{})
}
