package interfaces

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/model"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/comic"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/scrapper"
)

type Scrapper interface {
	SetConfig(log *model.Log, webURL string) Scrapper
	Close()
	ClosePage()
	CallPage(url string) error
	GetPageTitle() string
	GetReaderTitle() string
	GetPageInfo(title string) string
	GetPageInfoList(title string) []string
	GetPageEpisodes() []scrapper.Episode
	GetReaderImageURLs() []string
	NextReaderImage() string

	GeneratePageInfo() comic.Info
	GenerateEpisodes(initURL string, comicInfo *comic.Info)
	GenerateComicInfo(logger *model.Log, path string) (*comic.Info, error)
}
