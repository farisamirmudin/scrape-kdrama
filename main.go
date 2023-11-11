package main

import (
	"html/template"

	"github.com/farisamirmudin/gowatch/lib"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()
	godotenv.Load()
	router.GET("/", func(ctx *gin.Context) {
		indexTemplate := template.Must(template.ParseFiles("views/index.html"))
		indexTemplate.Execute(ctx.Writer, nil)
	})
	router.GET("/search", func(ctx *gin.Context) {
		films := lib.GetFilm(ctx.Query("keyword"))
		resultsTemplate := template.Must(template.ParseFiles("views/results.html"))
		resultsTemplate.Execute(ctx.Writer, films)
	})
	router.GET("/episodes", func(ctx *gin.Context) {
		episodes := lib.GetFilmEpisodes(ctx.Query("path"))
		episodesTemplate, _ := template.New("episodes.html").Funcs(template.FuncMap{"ExtractTitleAndEpisode": lib.ExtractTitleAndEpisode}).ParseFiles("views/episodes.html")
		episodesTemplate.Execute(ctx.Writer, episodes)
	})
	router.GET("/servers", func(ctx *gin.Context) {
		player := lib.GetFilmServers(ctx.Query("path"))
		playerTemplate := template.Must(template.ParseFiles("views/player.html"))
		playerTemplate.Execute(ctx.Writer, player)
	})

	router.Run()
}
