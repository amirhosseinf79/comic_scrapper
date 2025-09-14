package scrapper

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/amirhosseinf79/comic_scrapper/internal/dto/scrapper"
)

func (r *rodS) GetReaderTitle() string {
	pageTitle := r.page.MustInfo().Title
	re := regexp.MustCompile("(.+) - Read")
	raw := re.FindStringSubmatch(pageTitle)
	if len(raw) > 0 {
		return raw[1]
	}
	return pageTitle
}

func (r *rodS) GetPageTitle() string {
	getTitleElement, err := r.page.Element(r.infoContainer.ComicTitle)
	if err != nil {
		return ""
	}
	title, err := getTitleElement.Text()
	if err != nil {
		return ""
	}
	return title
}

func (r *rodS) GetPageCover() string {
	getImageElement, err := r.page.Element(r.infoContainer.ComicCover)
	if err != nil {
		return ""
	}
	imgSrc, err := getImageElement.Attribute("src")
	if err != nil || imgSrc == nil {
		return ""
	}
	return r.webURL + *imgSrc
}

func (r *rodS) GetPageGenres() []string {
	genresElement, _ := r.page.Elements("p:nth-of-type(1) a.dotUnder")
	genres := make([]string, 0)
	for _, g := range genresElement {
		genres = append(genres, g.MustText())
	}
	return genres
}

func (r *rodS) GetPageInfo(title string) string {
	rawTitle := fmt.Sprintf("%v:", title)
	infoElements, _ := r.page.Elements(r.infoContainer.InfoContainer)
	for _, p := range infoElements {
		txt := strings.TrimSpace(p.MustText())
		if !strings.HasPrefix(txt, r.infoContainer.InfoTitles.Description) {
			if strings.HasPrefix(txt, rawTitle) {
				raw := txt[len([]rune(rawTitle)):]
				raws := strings.Split(raw, r.infoContainer.InfoTitles.Views)
				if len(raws) > 0 {
					raw = raws[0]
				}
				return strings.TrimSpace(raw)
			}
		} else {
			p.MustNext()
			next := p.MustNext().MustText()
			return strings.TrimSpace(next)
		}
	}
	return ""
}

func (r *rodS) GetPageInfoList(title string) []string {
	rawTitle := fmt.Sprintf("%v:", title)
	paragraphs, _ := r.page.Elements(r.infoContainer.InfoContainer)
	for _, p := range paragraphs {
		txt := strings.TrimSpace(p.MustText())
		if !strings.HasPrefix(txt, r.infoContainer.InfoTitles.Description) {
			if strings.HasPrefix(txt, rawTitle) {
				raw := strings.TrimSpace(txt[len([]rune(rawTitle)):])
				list := strings.Split(raw, ",")
				return list
			}
		}
	}
	return nil
}

func (r *rodS) GetPageEpisodes() []scrapper.Episode {
	trs, _ := r.page.Elements("table.listing tbody tr")
	episodes := make([]scrapper.Episode, 0)
	for _, tr := range trs {
		linkEl, _ := tr.Element("a")
		if linkEl == nil {
			continue
		}
		title := linkEl.MustText()
		if title == "" {
			continue
		}
		href, err := linkEl.Attribute("href")
		if err != nil || href == nil {
			continue
		}
		episodes = append(episodes, scrapper.Episode{
			Title: title,
			Url:   *href + "&quality=hq",
		})
	}
	return episodes
}

func (r *rodS) GetReaderImageURLs() []string {
	images, _ := r.page.Elements(r.infoContainer.ImageDiv)
	imgURLs := make([]string, 0)
	for _, image := range images {
		imageURL, err := image.Attribute("src")
		if err != nil || imageURL == nil {
			continue
		}
		if len(*imageURL) > 10 {
			imgURLs = append(imgURLs, *imageURL)
		}
	}
	return imgURLs
}

func (r *rodS) NextReaderImage() string {
	nextBtn, err := r.page.Element("#btnNext")
	if err == nil && nextBtn != nil {
		nextBtn.MustClick()
		time.Sleep(1 * time.Second)
	} else {
		time.Sleep(1 * time.Second)
		return r.NextReaderImage()
	}
	return r.page.MustInfo().URL[len(r.webURL):]
}
