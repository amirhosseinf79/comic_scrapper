package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/amirhosseinf79/comic_scrapper/internal/domain/model"
	scrapper2 "github.com/amirhosseinf79/comic_scrapper/internal/service/scrapper"
	system2 "github.com/amirhosseinf79/comic_scrapper/internal/service/system"
)

func main() {
	webURL := "https://readcomiconline.li"
	logger := model.InitLog()
	scrapper := scrapper2.NewRod(false).SetConfig(logger, webURL)
	system := system2.NewSystem("#/comic/Dreadstar-and-Company")
	defer scrapper.Close()

	err := scrapper.CallPage("/Comic/Dreadstar-and-Company")
	if err != nil {
		log.Fatal(err)
	}

	comic := scrapper.GeneratePageInfo()
	fmt.Println("Title:", comic.Title)
	fmt.Println("PubDate:", comic.PublishDate)
	fmt.Println("Genres:", comic.Categories)
	fmt.Println("Publisher:", comic.Publisher)
	fmt.Println("Writers:", comic.Authors)
	fmt.Println("Status:", comic.Status)
	fmt.Println("Has Description?", len(comic.Description) > 0)
	fmt.Println("Has Thumbnail?", len(comic.ThumbnailFileAddress) > 0)

	episodes := scrapper.GetPageEpisodes()
	firstEpisode := episodes[len(episodes)-1]
	fmt.Println("First Episode:", firstEpisode.Title)

	for _, episode := range episodes {
		scrapper.GenerateEpisodes(episode.Url, &comic)
	}

	path := system.MakeDir("")
	content, err := json.Marshal(comic)
	if err != nil {
		log.Fatal(err)
	}
	err = system.SaveURL(string(content), path)
	if err != nil {
		fmt.Println(err)
	}

	logger.Output = comic

	fmt.Println("Done.")
	fmt.Printf("Overal Status: %v", logger.Status.String())
	fmt.Printf("Info Scrapped: %v", logger.HasInfo)
	fmt.Printf("Total Episode Scrapped: %d/%d", logger.ProcessedEpisodes, logger.TotalEpisodes)
	fmt.Printf("Total File Scrapped: %d/%d", logger.ProcessedFiles, logger.TotalFiles)
	fmt.Println("Failed logs:")
	for _, cmd := range logger.Console {
		if strings.Contains(cmd, "Failed") {
			fmt.Println(cmd)
		}
	}
}
