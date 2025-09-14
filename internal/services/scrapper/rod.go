package scrapper

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/amirhosseinf79/comic_scrapper/internal/domain/interfaces"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/comic"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/scrapper"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type rodS struct {
	webURL               string
	browser              *rod.Browser
	page                 *rod.Page
	infoContainer        scrapper.Container
	episodeListContainer string
}

func NewRod(headless bool) interfaces.Rod {
	u := launcher.New().Headless(headless).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	localBrowser := browser.MustIncognito()

	return &rodS{
		browser: localBrowser,
	}
}

func (r *rodS) Setconfig(webURL string) interfaces.Rod {
	return &rodS{
		webURL:  webURL,
		browser: r.browser,
		infoContainer: scrapper.Container{
			ComicTitle:    "a.bigChar",
			ComicCover:    "div.rightBox div.barContent img",
			InfoContainer: ".barContent p",
			ImageDiv:      "div#divImage img",
			InfoTitles: scrapper.InfoMap{
				Genres:          "Genres",
				Publisher:       "Publisher",
				Writer:          "Writer",
				Artist:          "Artist",
				PublicationDate: "Publication date",
				Status:          "Status",
				Description:     "Summary",
				Views:           "Views",
			},
		},
		episodeListContainer: "listing",
	}
}

func (r *rodS) Close() {
	_ = r.browser.Close()
}

func (r *rodS) ClosePage() {
	r.page.MustClose()
}

func (r *rodS) CallPage(url string) error {
	ctx := context.Background()
	var width = 600
	var height = 300
	for {
		var err error
		r.page, err = r.browser.Context(ctx).Page(proto.TargetCreateTarget{
			URL:    r.webURL + url,
			Width:  &width,
			Height: &height,
		})
		if err == nil {
			err = r.page.WaitLoad()
		}
		if err != nil {
			if errors.Is(err, context.Canceled) {
				fmt.Println("Cancelled")
				return err
			} else if errors.Is(err, context.DeadlineExceeded) {
				fmt.Println("DeadlineExceeded")
				return err
			}
			fmt.Println(err)
			continue
		}
		break
	}
	return nil
}

func (r *rodS) GenerateAuthors(list []string) []comic.Author {
	authors := make([]comic.Author, 0)
	for _, item := range list {
		splitFullName := strings.SplitN(item, " ", 2)
		firstName := splitFullName[0]
		lastName := ""
		if len(splitFullName) > 1 {
			lastName = splitFullName[1]
		}
		authors = append(authors, comic.Author{
			FirstName: firstName,
			LastName:  lastName,
			Type:      "string",
		})
	}
	return authors
}

func (r *rodS) GenerateDateTime(str string) time.Time {
	layout := "January 2006"
	t, err := time.Parse(layout, str)
	if err != nil {
		return time.Now()
	}
	return t
}

func (r *rodS) GenerateStatus(str string) int {
	switch str {
	case "Completed":
		return 2
	case "Ongoing":
		return 1
	default:
		return 3
	}
}

func (r *rodS) GeneratePageInfo() comic.Info {
	pageTitle := r.GetPageTitle()
	pageCover := r.GetPageCover()
	pageStatus := r.GetPageInfo(r.infoContainer.InfoTitles.Status)
	writerList := r.GetPageInfoList(r.infoContainer.InfoTitles.Writer)
	pubDateStr := r.GetPageInfo(r.infoContainer.InfoTitles.PublicationDate)

	return comic.Info{
		Title:                pageTitle,
		ImageFileAddress:     pageCover,
		BannerFileAddress:    pageCover,
		ThumbnailFileAddress: pageCover,
		Status:               r.GenerateStatus(pageStatus),
		Publisher:            r.GetPageInfo(r.infoContainer.InfoTitles.Publisher),
		Description:          r.GetPageInfo(r.infoContainer.InfoTitles.Description),
		PublishDate:          r.GenerateDateTime(pubDateStr),
		Categories:           r.GetPageInfoList(r.infoContainer.InfoTitles.Genres),
		Authors:              r.GenerateAuthors(writerList),
	}
}

func (r *rodS) GenerateEpisodeFiles(index int, urls []string) *comic.EpisodeFile {
	imageTitle := fmt.Sprintf("image_%d", index)
	if len(urls) == 0 {
		return nil
	}
	if len(urls) > 1 {
		return &comic.EpisodeFile{
			Title:        imageTitle,
			Description:  imageTitle,
			ImageAddress: urls[1],
		}
	}
	return &comic.EpisodeFile{
		Title:        imageTitle,
		Description:  imageTitle,
		ImageAddress: urls[0],
	}
}

func (r *rodS) GenerateEachEpisodeFiles(comicInfo *comic.Info, episodes []comic.Episode) ([]comic.Episode, string, bool) {
	currentEpisodeTitle := r.GetReaderTitle()
	fmt.Println("Generating episodes for:", currentEpisodeTitle)
	episode := comic.Episode{
		Title:        currentEpisodeTitle,
		Description:  currentEpisodeTitle,
		ImageAddress: comicInfo.ThumbnailFileAddress,
		Files:        make([]comic.EpisodeFile, 0),
	}
	isFinished := false
	initURL := ""
	newURL := ""
	counter := 0
	for {
		counter++
		episodeImageFiles := r.GetReaderImageURLs()
		fmt.Println("Generating", len(episodeImageFiles), "files:", counter)
		generatedFiles := r.GenerateEpisodeFiles(counter, episodeImageFiles)
		if generatedFiles != nil {
			episode.Files = append(episode.Files, *generatedFiles)
		}

		newURL = r.NextReaderImage()
		if newURL == initURL {
			fmt.Println("Finished:", len(episodes), "Episodes Generated")
			isFinished = true
			break
		}
		initURL = newURL
		newEpisodeTitle := r.GetReaderTitle()
		if newEpisodeTitle != currentEpisodeTitle {
			fmt.Println("Going to generate next episode:", newEpisodeTitle)
			break
		}
	}
	episodes = append(episodes, episode)
	return episodes, newURL, isFinished
}

func (r *rodS) GenerateEpisodes(initURL string, comicInfo *comic.Info) {
	episodes := make([]comic.Episode, 0)
	isFinished := false
	newURL := ""

	for {
		r.ClosePage()
		err := r.CallPage(initURL)
		if err != nil {
			fmt.Println(err)
			break
		}
		episodes, newURL, isFinished = r.GenerateEachEpisodeFiles(comicInfo, episodes)
		if isFinished {
			break
		}
		initURL = newURL
	}
	comicInfo.Episodes = episodes
}
