package downloader

import (
	"fmt"
	"net/url"
	"os"
	"path"
	"strings"
)

type absURL string
type staticTag string

type Storage interface {
	SaveFile(URL *url.URL, curRepoPath string) (*os.File, error)
	MkDir(dir *url.URL) (string, error)
}

type FileStorage struct{}

func (fs *FileStorage) SaveFile(URL *url.URL, curRepoPath string) (*os.File, error) {
	trimmed := strings.Trim(URL.Path, "/")
	fPrefix := strings.ReplaceAll(trimmed, "/", "-") + ".html"

	if _, err := os.Stat(fPrefix); !os.IsNotExist(err) {
		// file already exists
		return nil, fmt.Errorf("%w %v", errAlreadyDownload, fPrefix)
	}

	filePath := path.Join(curRepoPath, fPrefix)
	fmt.Printf("\033[33mFilePath\033[0m> %v\n", filePath)

	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w %v", errFileCreate, err)
	}
	return file, nil

}

func (fs *FileStorage) MkDir(u *url.URL) (string, error) {
	curRepoPath := BaseDir + u.Host

	if err := os.MkdirAll(curRepoPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("%w: %v", errMkDir, err)
	}

	for _, f := range Dirs {
		assetPath := path.Join(curRepoPath, f)
		if err := os.MkdirAll(assetPath, os.ModePerm); err != nil {
			fmt.Printf("%v: %v\n, %v", errMkDir, curRepoPath, err)
			continue
		}
	}

	return curRepoPath, nil
}
