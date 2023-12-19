package main

import (
	"net/http"
	"github.com/labstack/echo/v4"
)

func main() {

	e := echo.New()
	//CORS設定でlocalhost:3000とlocalhost:3001を許可

	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Access-Control-Allow-Headers, Authorization, X-Requested-With")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			return next(c)
		}
	})

	e.Static("/", "static/build")

	e.GET("/api/get_twitch_streaming_user", func(c echo.Context) error {
		twitchLives, err := GetTwtichLives()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string][]map[string]interface{}{
			"users": twitchLives,
		})
	})

	e.GET("/api/get_youtube_movies", func(c echo.Context) error {
		ytbMovieList, err := GetYoutubeMovieInfo()
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, map[string][]map[string]interface{}{
			"movies": ytbMovieList,
		})

	})

	e.Start(":3000")

}

