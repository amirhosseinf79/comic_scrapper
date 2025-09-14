package main

import (
	"encoding/json"
	"fmt"
	"log"

	scrapper2 "github.com/amirhosseinf79/comic_scrapper/internal/services/scrapper"
	system2 "github.com/amirhosseinf79/comic_scrapper/internal/services/system"
)

func main() {
	webURL := "https://readcomiconline.li"
	scrapper := scrapper2.NewRod(false).Setconfig(webURL)
	system := system2.NewSystem("#/comic/")
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

	//scrapper.GenerateEpisodes(firstEpisode.Url, &comic)

	path := system.MakeDir("")
	content, err := json.Marshal(comic)
	if err != nil {
		log.Fatal(err)
	}
	err = system.SaveURL(string(content), path)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Done.")
	fmt.Println("Total Episode Downloaded:", len(comic.Episodes))
}
