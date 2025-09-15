package scrapper

import (
	"fmt"
	"strings"

	"github.com/amirhosseinf79/comic_scrapper/internal/domain/enum"
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/interfaces"
	"github.com/amirhosseinf79/comic_scrapper/internal/domain/model"
	"github.com/amirhosseinf79/comic_scrapper/internal/dto/comic"
)

func GenerateComicInfo(scrapper0 interfaces.Scrapper, logger *model.Log, path string) (*comic.Info, error) {
	webURL := "https://readcomiconline.li"
	scrapper := scrapper0.SetConfig(logger, webURL)
	err := scrapper.CallPage(path)
	defer scrapper.ClosePage()
	defer scrapper.Close()
	if err != nil {
		return nil, err
	}

	comicM := scrapper.GeneratePageInfo()
	fmt.Println("Title:", comicM.Title)
	fmt.Println("PubDate:", comicM.PublishDate)
	fmt.Println("Genres:", comicM.Categories)
	fmt.Println("Publisher:", comicM.Publisher)
	fmt.Println("Writers:", comicM.Authors)
	fmt.Println("Status:", comicM.Status)
	fmt.Println("Has Description?", len(comicM.Description) > 0)
	fmt.Println("Has Thumbnail?", len(comicM.ThumbnailFileAddress) > 0)

	episodes := scrapper.GetPageEpisodes()
	if len(episodes) == 0 {
		scrapper.ConsoleAdd("GenerateComicInfo", enum.Failed, "no episode found")
		return nil, fmt.Errorf("no episode found")
	}
	firstEpisode := episodes[len(episodes)-1]
	fmt.Println("First Episode:", firstEpisode.Title)

	for _, episode := range episodes {
		scrapper.GenerateEpisodes(episode.Url, &comicM)
	}

	logger.Output = comicM

	fmt.Println("Done.")
	fmt.Printf("Overal Status: %v\n", logger.Status.String())
	fmt.Printf("Info Scrapped: %v\n", logger.HasInfo)
	fmt.Printf("Total Episode Scrapped: %d/%d\n", logger.ProcessedEpisodes, logger.TotalEpisodes)
	fmt.Printf("Total File Scrapped: %d/%d\n", logger.ProcessedFiles, logger.TotalFiles)
	fmt.Println("Failed logs:")
	for _, cmd := range logger.Console {
		if strings.Contains(cmd, "Failed") {
			fmt.Println(cmd)
		}
	}
	return &comicM, nil
}
