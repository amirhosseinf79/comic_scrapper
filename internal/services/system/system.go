package system

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/amirhosseinf79/comic_scrapper/internal/domain/interfaces"
)

type sys struct {
	basePath string
}

func NewSystem(path string) interfaces.System {
	return &sys{
		basePath: path,
	}
}

func (s *sys) MakeDir(folder string) string {
	path := s.basePath + folder
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return path
}

func (s *sys) DownloadFile(url, filepath string, number int) error {
	imagePath := filepath + "/" + fmt.Sprintf("page_%03d.jpg", number)
	fmt.Println("Downloading image", number, "to", imagePath)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(imagePath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func (s *sys) SaveURL(url, filepath string) error {
	imagePath := filepath + "/urls.json"
	file, err := os.OpenFile(imagePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(url)
	if err != nil {
		return err
	}
	return nil
}
