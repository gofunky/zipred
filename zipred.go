package zipred

import (
	"errors"
	"github.com/francoispqt/onelog"
	"github.com/gofunky/zipstream"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var log *onelog.Logger

func init() {
	log = onelog.New(
		os.Stdout,
		onelog.WARN|onelog.ERROR|onelog.FATAL,
	)
}

// UseLogger sets a custom Logger for the package.
func UseLogger(logger *onelog.Logger) {
	log = logger
}

const (
	errorConst = "error"
)

// Zipred allows the implementation of a context that dynamically defines the boundaries of the zip filtering.
type Zipred interface {
	// URL to download the archive from
	URL() string
	// Predicate indicates if the given file should be read or not.
	// It is to return a zero string to discard the file, otherwise the key name.
	// If error is nonempty, the download is aborted and the error is passed on.
	Predicate(fileInfo os.FileInfo) (key string, err error)
	// Done indicates if enough data has been read and the download can be aborted ahead of the EOF.
	// isEOF is true if the end of the zip file has been reached.
	// If error is nonempty, the download is aborted and the error is passed on.
	Done(isEOF bool) (finish bool, err error)
}

// FilterFileInfo filters the given URL's zip file and returns all file names on which the given predicate applies.
func FilterFileInfo(z Zipred) (keys []string, err error) {
	keys = make([]string, 0)
	err = parseFiles(z, func(reader *zipstream.Reader, predKey string) (storeErr error) {
		keys = append(keys, predKey)
		return nil
	})
	if err != nil {
		log.ErrorWith("the archive filtering failed").
			Err(errorConst, err).Write()
		return nil, err
	}
	log.Info("the archive's files were successfully downloaded and filtered")
	return keys, nil
}

// FilterZipContent filters the given URL's zip file and returns all file contents on which the given predicate applies.
func FilterZipContent(z Zipred) (files map[string][]byte, err error) {
	files = make(map[string][]byte, 0)
	err = parseFiles(z, func(reader *zipstream.Reader, predKey string) (storeErr error) {
		err = readNext(predKey, reader, files)
		if err != nil {
			log.ErrorWith("the given zipped file could not be read").
				String("file", predKey).Err(errorConst, err).Write()
			return err
		}
		return nil
	})
	if err != nil {
		log.ErrorWith("the archive filtering failed").
			Err(errorConst, err).Write()
		return nil, err
	}
	log.Info("the archive's files were successfully downloaded and filtered")
	return files, nil
}

// parseFiles downloads the given URL, parses the zip and passes the zip reader and predicate key to the given func.
func parseFiles(z Zipred, store func(reader *zipstream.Reader, predKey string) (storeErr error)) (useErr error) {
	return downloadFiles(z.URL(), func(data io.ReadCloser) (useErr error) {
		zipReader := zipstream.NewReader(data)
		for {
			header, useErr := zipReader.Next()
			if useErr == io.EOF {
				_, useErr := z.Done(true)
				if useErr != nil {
					log.ErrorWith("the archive parser has reached an erroneous state").
						Bool("isEOF", true).Err(errorConst, useErr).Write()
					return useErr
				}
				return nil
			} else if useErr != nil {
				log.ErrorWith("the downloading archive could not be read").
					Err(errorConst, useErr).Write()
				return useErr
			}
			predKey, useErr := z.Predicate(header.FileInfo())
			if useErr != nil {
				log.ErrorWith("the given predicate failed").
					Err(errorConst, useErr).Write()
				return useErr
			}
			if predKey != "" {
				useErr := store(zipReader, predKey)
				if useErr != nil {
					log.ErrorWith("the archive parser could not store the result").
						Err(errorConst, useErr).Write()
					return useErr
				}
			}
			done, useErr := z.Done(false)
			if useErr != nil {
				log.ErrorWith("the archive parser has reached an erroneous state").
					Err(errorConst, useErr).Write()
				return useErr
			}
			if done {
				return nil
			}
		}
	})
}

// downloadFiles downloads the given URL and passes the resulting reader to the given usage func.
func downloadFiles(URL string, usage func(data io.ReadCloser) (useErr error)) (err error) {
	// URL must not be empty
	if URL == "" {
		log.ErrorWith("the given URL is empty").
			String("url", URL).Write()
		return errors.New("URL is empty")
	}

	// Download file ony the fly
	resp, err := http.Get(URL)
	if err != nil {
		log.ErrorWith("the given archive archive could not be downloaded").
			String("url", URL).Err(errorConst, err).Write()
		return err
	}
	defer resp.Body.Close()

	// Apply predicates
	err = usage(resp.Body)

	return err
}

// readNext reads the next template and stores it in the given target map.
func readNext(key string, reader io.Reader, target map[string][]byte) (err error) {
	if key == "" {
		return errors.New("empty key was given")
	}
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	target[key] = content
	return nil
}
