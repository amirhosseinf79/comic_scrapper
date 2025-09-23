package comic

import "time"

type Author struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Type      string `json:"type"`
}

type Episode struct {
	Title        string        `json:"title"`
	Description  string        `json:"description"`
	ImageAddress string        `json:"imageAddress"`
	Files        []EpisodeFile `json:"files"`
}

type EpisodeFile struct {
	Title        string `json:"title"`
	Description  string `json:"description"`
	ImageAddress string `json:"imageAddress"`
}

type Info struct {
	Title                string    `json:"title"`
	Description          string    `json:"description"`
	Status               int       `json:"status"`
	PublishDate          time.Time `json:"publishDate"`
	EndDate              time.Time `json:"endDate"`
	ImageFileAddress     string    `json:"imageFileAddress"`
	BannerFileAddress    string    `json:"bannerFileAddress"`
	ThumbnailFileAddress string    `json:"thumbnailFileAddress"`
	Publisher            string    `json:"publisher"`
	Authors              []Author  `json:"authors"`
	Categories           []string  `json:"categories"`
	Episodes             []Episode `json:"episodes"`
}
