package interfaces

import (
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/comic"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/scrapper"
)

type Rod interface {
	Setconfig(url string) Rod
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
}
