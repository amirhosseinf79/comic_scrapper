package scrapper

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/amirhosseinf79/comic_scrapper/internal/domain/enum"
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/interfaces"
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/model"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/comic"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/scrapper"
	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/launcher"
	"github.com/go-rod/rod/lib/proto"
)

type rodS struct {
	ctx                  context.Context
	ctxCancel            context.CancelFunc
	log                  *model.Log
	webURL               string
	browser              *rod.Browser
	page                 *rod.Page
	status               scrapper.Status
	infoContainer        scrapper.Container
	logger               interfaces.Logger
	episodeListContainer string
}

func New(headless bool, logger interfaces.Logger) interfaces.Scrapper {
	u := launcher.New().Headless(headless).MustLaunch()
	browser := rod.New().ControlURL(u).MustConnect()
	localBrowser := browser.MustIncognito()

	return &rodS{
		logger:  logger,
		browser: localBrowser,
		status: scrapper.Status{
			Failed:  enum.Failed,
			Pending: enum.Pending,
			Success: enum.Succeed,
		},
	}
}

func (r *rodS) SetConfig(log *model.Log, webURL string) interfaces.Scrapper {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	r.ctxCancel = cancel
	r.ctx = ctx

	r.log = log
	r.webURL = webURL
	r.episodeListContainer = "listing"
	r.infoContainer = scrapper.Container{
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
	}
	return r
}

func (r *rodS) Close() {
	_ = r.browser.Close()
}

func (r *rodS) ClosePage() {
	r.page.MustClose()
}

func (r *rodS) CallPage(url string) error {
	r.ConsoleAdd("CallPage", r.status.Pending, url)
	var width = 600
	var height = 300
	var err error
	r.page, err = r.browser.Context(r.ctx).Page(proto.TargetCreateTarget{
		URL:    r.webURL + url,
		Width:  &width,
		Height: &height,
	})
	if err != nil {
		r.ConsoleAdd("CallPage", r.status.Failed, err.Error()+": "+url)
		if errors.Is(err, context.Canceled) {
			return err
		} else if errors.Is(err, context.DeadlineExceeded) {
			return err
		}
		return err
	}
	err = r.page.WaitLoad()
	if err != nil {
		r.ConsoleAdd("CallPage", r.status.Failed, err.Error()+": "+url)
	}
	return err
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
	r.ConsoleAdd("GenerateDateTime", r.status.Pending, str)
	layout := "January 2006"
	t, err := time.Parse(layout, str)
	if err != nil {
		layout = "Jan 2, 2006"
		t, err = time.Parse(layout, str)
		if err != nil {
			r.ConsoleAdd("GenerateDateTime", r.status.Failed, err.Error())
			return time.Now()
		}
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
	r.ConsoleAdd("GenerateEachEpisodeFiles", r.status.Pending, currentEpisodeTitle)
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
		log := fmt.Sprint("Generating", len(episodeImageFiles), "files:", counter)
		r.ConsoleAdd("GenerateEachEpisodeFiles", r.status.Pending, log)
		generatedFiles := r.GenerateEpisodeFiles(counter, episodeImageFiles)
		if generatedFiles != nil {
			r.ConsoleAdd("GenerateEachEpisodeFiles", r.status.Success, log)
			episode.Files = append(episode.Files, *generatedFiles)
		} else {
			r.ConsoleAdd("GenerateEachEpisodeFiles", r.status.Failed, log)
		}

		newURL = r.NextReaderImage()
		if newURL == initURL {
			log = fmt.Sprint("Finished:", "Episodes Generated")
			r.ConsoleAdd("GenerateEachEpisodeFiles", r.status.Success, log)
			isFinished = true
			break
		}
		initURL = newURL
		newEpisodeTitle := r.GetReaderTitle()
		if newEpisodeTitle != currentEpisodeTitle {
			log = fmt.Sprint("Going to generate next episode:", newEpisodeTitle)
			r.ConsoleAdd("GenerateEachEpisodeFiles", r.status.Success, log)
			isFinished = true
			break
		}
	}
	episodes = append(episodes, episode)

	r.log.TotalFiles += counter
	r.log.ProcessedFiles += len(episode.Files)

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
			break
		}
		episodes, newURL, isFinished = r.GenerateEachEpisodeFiles(comicInfo, episodes)
		if isFinished {
			break
		}
		initURL = newURL
	}
	comicInfo.Episodes = append(comicInfo.Episodes, episodes...)
	r.log.ProcessedEpisodes = len(comicInfo.Episodes)
}
