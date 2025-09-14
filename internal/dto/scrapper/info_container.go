package scrapper

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
