// Package downloader provides utilities to download web pages and save them to files
package downloader

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Service interface {
	GetPage(u string) (GettedFile, error)
	ParsePage(f *GettedFile) error
}

type GettedFile struct {
	file *os.File
	URL  *url.URL
}

type service struct {
	storage Storage
}

func NewService(storage Storage) Service {
	return &service{storage: storage}
}

const BaseDir = "list/"

var Dirs = []string{
	// TODO:"/recursive",
	"/assets/img",
	"/assets/css",
	"/assets/js",
}

var (
	errGetURL          = errors.New("ошибка запроса: ")
	errFileCreate      = errors.New("ошибка создания файла: ")
	errFileCopy        = errors.New("ошибка копирования файла: ")
	errFileOpen        = errors.New("ошибка копирования файла: ")
	errAlreadyDownload = errors.New("файл с таким именем уже существует: ")
	errURLParse        = errors.New("ошибка парсинга URL: ")
	errMkDir           = errors.New("ошибка создания директории: ")
)

// GetPage implements Service
func (s *service) GetPage(inURL string) (GettedFile, error) {

	//---Get Page---
	resp, err := http.Get(inURL)
	if err != nil {
		return GettedFile{}, fmt.Errorf("%w %v", errGetURL, err)
	}
	defer resp.Body.Close()

	//---URL parsing with lib [net/url]---
	u, err := url.Parse(inURL)
	if err != nil {
		return GettedFile{}, fmt.Errorf("%w %v", errURLParse, err)
	}
	fmt.Printf("\033[31mPATH URL\033[0m> %v\n", u.Host)

	// ---Make directory---
	curRepoPath, err := s.storage.MkDir(u)
	if err != nil {
		return GettedFile{}, fmt.Errorf("%w %v", errMkDir, err)
	}

	//---File Creation---
	file, err := s.storage.SaveFile(u, curRepoPath)
	if err != nil {
		return GettedFile{}, fmt.Errorf("%w %v", errFileCreate, err)
	}
	defer file.Close()

	//---Body copy---
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return GettedFile{}, fmt.Errorf("%w %v", errFileCopy, err)
	}

	GettedFile := GettedFile{file: file, URL: u}

	fmt.Printf("Страница сохранена под именем: %v\n", file.Name())
	return GettedFile, nil
}

func (s *service) ParsePage(gf *GettedFile) error {
	// ---Open file---
	file, err := os.Open(gf.file.Name())
	if err != nil {
		return fmt.Errorf("%w: %v", errFileOpen, err)
	}
	defer file.Close()

	// ---goquery document---
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return fmt.Errorf("%w: %v", err, err)
	}

	var parseErrs []error 
	//---goqery parse--
	doc.Find("img,link[rel='stylesheet'],script[src]").Each(func(i int, s *goquery.Selection) {
		var src string
		if val, exists := s.Attr("src"); exists {
			src = val
		}
		if val, exists := s.Attr("href"); exists {
			src = val
		}
		if src == "" {
			return
		}

		// --- Build abs URL ---
		absURL := src
		if !strings.HasPrefix(src, "http") {
			ref, err := url.Parse(src)
			if err != nil {
				parseErrs = append(parseErrs, fmt.Errorf("ошибка парсинга ссылки %q: %v", src, err))
				return
			}
			absURL = gf.URL.ResolveReference(ref).String()
		}

		resp, err := http.Get(absURL)
		if err != nil {
			parseErrs = append(parseErrs, fmt.Errorf("ошибка загрузки %q: %v", absURL, err))
			return
		}
		defer resp.Body.Close()

		var filePath string

		switch {
		case s.Is("img"):
			filePath, err = GetStatic(absURL, gf, "img")
		case s.Is("link"):
			filePath, err = GetStatic(absURL, gf, "link")
		case s.Is("script"):
			filePath, err = GetStatic(absURL, gf, "script")
		}

		if err != nil {
			parseErrs = append(parseErrs, fmt.Errorf("%w: %v", errGetURL, err))
			return
		}
		if s.Is("link") {
			s.SetAttr("href", filePath)
		} else {
			s.SetAttr("src", filePath)
		}
	})

	// ---err return---
	if len(parseErrs) > 0 {
		var sb strings.Builder
		for _, e := range parseErrs {
			sb.WriteString(e.Error() + "\n")
		}
		return fmt.Errorf("ошибки при разборе страницы:\n%s", sb.String())
	}

	return nil
}

func GetStatic(absURL string, gf *GettedFile, fileType string) (string, error) {
	// ---Parse URL---
	u, err := url.Parse(absURL)
	if err != nil {
		return "", fmt.Errorf("%w %v", errURLParse, err)
	}

	// --- Determine directory based on fileType ---
	var subDir string
	switch fileType {
	case "img":
		subDir = "/assets/img"
	case "css", "link":
		subDir = "/assets/css"
	case "js", "script":
		subDir = "/assets/js"
	default:
		subDir = "/assets/other"
	}

	// --- Create full directory path ---
	basePath := fmt.Sprintf("%s%s%s", BaseDir, gf.URL.Host, subDir)
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return "", fmt.Errorf("%w %v", errMkDir, err)
	}

	// --- Extract file name from URL path ---
	splitted := strings.Split(u.Path, "/")
	filename := splitted[len(splitted)-1]
	if filename == "" {
		filename = "index"
	}

	// --- Full local file path ---
	localPath := fmt.Sprintf("%s/%s", basePath, filename)

	// --- Check if file already exists ---
	if _, err := os.Stat(localPath); err == nil {
		fmt.Printf("Файл уже существует: %s\n", localPath)
		return localPath, nil
	}

	// --- Download file ---
	resp, err := http.Get(absURL)
	if err != nil {
		return "", fmt.Errorf("%w %v", errGetURL, err)
	}
	defer resp.Body.Close()

	// --- Create file ---
	out, err := os.Create(localPath)
	if err != nil {
		return "", fmt.Errorf("%w %v", errFileCreate, err)
	}
	defer out.Close()

	// --- Copy body ---
	if _, err = io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("%w %v", errFileCopy, err)
	}

	// --- Return relative path for HTML ---
	relPath := strings.TrimPrefix(localPath, BaseDir)
	return relPath, nil
}
