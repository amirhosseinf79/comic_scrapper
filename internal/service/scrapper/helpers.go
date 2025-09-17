package scrapper

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/amirhosseinf79/comic_scrapper/internal/domain/enum"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/scrapper"
	"github.com/go-rod/rod/lib/proto"
)

func (r *rodS) ConsoleAdd(state string, status enum.LogStatus, cmd ...string) {
	err := r.logger.AutoUpdate(r.log, state, status, cmd...)
	if err != nil {
		fmt.Println(err)
	}
}

func (r *rodS) GetReaderTitle() string {
	r.ConsoleAdd("GetReaderTitle", r.status.Pending)
	info, err := r.page.Info()
	if err != nil {
		r.ConsoleAdd("GetReaderTitle", r.status.Failed, err.Error())
		return ""
	}
	pageTitle := info.Title
	re := regexp.MustCompile("(.+) - Read")
	raw := re.FindStringSubmatch(pageTitle)
	if len(raw) > 0 {
		return raw[1]
	}
	if pageTitle == "" {
		r.ConsoleAdd("GetReaderTitle", r.status.Failed)
		r.log.HasInfo = false
	}
	return pageTitle
}

func (r *rodS) GetPageTitle() string {
	r.ConsoleAdd("GetPageTitle", r.status.Pending)
	getTitleElement, err := r.page.Timeout(5 * time.Second).Element(r.infoContainer.ComicTitle)
	if err != nil {
		return ""
	}
	title, err := getTitleElement.Text()
	if err != nil {
		r.ConsoleAdd("GetPageTitle", r.status.Failed)
		r.log.HasInfo = false
		return ""
	}
	return title
}

func (r *rodS) GetPageCover() string {
	r.ConsoleAdd("GetPageCover", r.status.Pending)
	getImageElement, err := r.page.Timeout(5 * time.Second).Element(r.infoContainer.ComicCover)
	if err != nil {
		return ""
	}
	imgSrc, err := getImageElement.Attribute("src")
	if err != nil || imgSrc == nil {
		r.ConsoleAdd("GetPageCover", r.status.Failed)
		r.log.HasInfo = false
		return ""
	}
	return r.webURL + *imgSrc
}

func (r *rodS) GetPageGenres() []string {
	r.ConsoleAdd("GetPageGenres", r.status.Pending)
	genresElement, _ := r.page.Timeout(5 * time.Second).Elements("p:nth-of-type(1) a.dotUnder")
	genres := make([]string, 0)
	for _, g := range genresElement {
		text, err := g.Text()
		if err != nil {
			continue
		}
		genres = append(genres, text)
	}
	if len(genres) == 0 {
		r.ConsoleAdd("GetPageGenres", r.status.Failed)
		r.log.HasInfo = false
	}
	return genres
}

func (r *rodS) GetPageInfo(title string) string {
	r.ConsoleAdd("GetPageInfo", r.status.Pending, title)
	rawTitle := fmt.Sprintf("%v:", title)
	infoElements, _ := r.page.Timeout(5 * time.Second).Elements(r.infoContainer.InfoContainer)
	for _, p := range infoElements {
		text, _ := p.Text()
		txt := strings.TrimSpace(text)
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
			sum, err := p.Next()
			if err != nil {
				r.log.HasInfo = false
				break
			}
			next, _ := sum.Text()
			return strings.TrimSpace(next)
		}
	}
	r.ConsoleAdd("GetPageInfo", r.status.Failed, title)
	r.log.HasInfo = false
	return ""
}

func (r *rodS) GetPageInfoList(title string) []string {
	r.ConsoleAdd("GetPageInfoList", r.status.Pending, title)
	rawTitle := fmt.Sprintf("%v:", title)
	paragraphs, _ := r.page.Timeout(5 * time.Second).Elements(r.infoContainer.InfoContainer)
	for _, p := range paragraphs {
		title, _ := p.Text()
		txt := strings.TrimSpace(title)
		if !strings.HasPrefix(txt, r.infoContainer.InfoTitles.Description) {
			if strings.HasPrefix(txt, rawTitle) {
				raw := strings.TrimSpace(txt[len([]rune(rawTitle)):])
				items := make([]string, 0)
				list := strings.Split(raw, ",")
				for _, item := range list {
					items = append(items, strings.TrimSpace(item))
				}
				return items
			}
		}
	}
	r.ConsoleAdd("GetPageInfoList", r.status.Failed, title)
	r.log.HasInfo = false
	return nil
}

func (r *rodS) GetPageEpisodes() []scrapper.Episode {
	r.ConsoleAdd("GetPageEpisodes", r.status.Pending)
	trs, _ := r.page.Timeout(5 * time.Second).Elements("table.listing tbody tr")
	episodes := make([]scrapper.Episode, 0)
	for _, tr := range trs {
		linkEl, _ := tr.Element("a")
		if linkEl == nil {
			continue
		}
		title, err := linkEl.Text()
		if title == "" || err != nil {
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
	if len(episodes) == 0 {
		r.ConsoleAdd("GetPageEpisodes", r.status.Failed)
	}
	r.log.TotalEpisodes = len(episodes)
	return episodes
}

func (r *rodS) GetReaderImageURLs() []string {
	r.ConsoleAdd("GetReaderImageURLs", r.status.Pending)
	images, _ := r.page.Timeout(10 * time.Second).Elements(r.infoContainer.ImageDiv)
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
	if len(imgURLs) == 0 {
		r.ConsoleAdd("GetReaderImageURLs", r.status.Failed)
	}
	return imgURLs
}

func (r *rodS) NextReaderImage() string {
	r.ConsoleAdd("NextReaderImage", r.status.Pending)
	waitTime := 50 * time.Millisecond

	for range 2 {
		nextBtn, err := r.page.Timeout(10 * time.Second).Element("#btnNext")
		if err == nil && nextBtn != nil {
			err = nextBtn.Timeout(10*time.Second).Click(proto.InputMouseButtonLeft, 1)
			if err != nil {
				r.ConsoleAdd("NextReaderImage", r.status.Failed, err.Error())
			}
			time.Sleep(waitTime)
		} else {
			r.ConsoleAdd("NextReaderImage", r.status.Failed)
		}
	}
	info, err := r.page.Info()
	if err != nil {
		r.ConsoleAdd("NextReaderImage", r.status.Failed, err.Error())
		return ""
	}
	return info.URL[len(r.webURL):]
}
