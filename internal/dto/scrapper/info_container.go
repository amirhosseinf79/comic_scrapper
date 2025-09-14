package scrapper

import "github.com/amirhosseinf79/comic_scrapper/internal/domain/enum"

type InfoMap struct {
	Genres          string
	Publisher       string
	Writer          string
	Artist          string
	PublicationDate string
	Status          string
	Description     string
	Views           string
}

type Container struct {
	ComicTitle    string
	ComicCover    string
	InfoContainer string
	ImageDiv      string
	InfoTitles    InfoMap
}

type Episode struct {
	Title string
	Url   string
}

type Status struct {
	Failed  enum.LogStatus
	Pending enum.LogStatus
	Success enum.LogStatus
}
