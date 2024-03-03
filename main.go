package main

import (
	"github.com/farisamirmudin/gowatch/lib"
	"github.com/farisamirmudin/gowatch/view/components"
	"github.com/farisamirmudin/gowatch/view/pages"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	router := gin.Default()
	godotenv.Load()
	router.GET("/", func(ctx *gin.Context) {
		pages.Index().Render(ctx, ctx.Writer)
	})
	router.GET("/search", func(ctx *gin.Context) {
		films := lib.GetFilm(ctx.Query("keyword"))
		components.Result(films).Render(ctx, ctx.Writer)
	})
	router.GET("/episodes", func(ctx *gin.Context) {
		episodes := lib.GetFilmEpisodes(ctx.Query("path"))
		components.EpisodeSelector(episodes).Render(ctx, ctx.Writer)
	})
	router.GET("/servers", func(ctx *gin.Context) {
		filmSrc := lib.GetFilmServers(ctx.Query("path"))
		components.Player(filmSrc).Render(ctx, ctx.Writer)
	})

	router.Run()
}
