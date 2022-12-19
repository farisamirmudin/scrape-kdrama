package main

import (
	"os"

	"github.com/farisamirmudin/lib"
	"github.com/fatih/color"
	"github.com/joho/godotenv"
)

func main() {

	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	keyword := lib.Search()
	searchLink, _ := lib.GetLink(os.Getenv("SEARCH_URL")+keyword, nil)
	chosenLink, title := lib.GetLink(os.Getenv("BASE_URL")+searchLink, &keyword)
	embedded := lib.GetEmbeddedLink(os.Getenv("BASE_URL") + chosenLink)
	videoLink := lib.GetVideoLink(embedded)
	color.Red(title)
	lib.Play(title, embedded, videoLink)
}
