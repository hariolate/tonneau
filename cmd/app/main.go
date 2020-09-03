package main

import (
	"context"
	"flag"
	"github.com/gin-gonic/gin"
	"gtihub.com/hariolate/tonneau/config"
	"gtihub.com/hariolate/tonneau/service"
	"log"
	"net/http"
)

var (
	configFile = flag.String("config", "config.json", "specify config file")
	debug      = flag.Bool("debug", true, "is in debug env or not")
)

func main() {
	flag.Parse()

	if !*debug {
		gin.SetMode(gin.ReleaseMode)
	}

	conf := config.MustFromFile(*configFile).MustParse()
	conf.Context = context.Background()

	srv := service.FromConfig(conf)

	router := gin.Default()
	router.Use(conf.MakeCheckMaintenanceStatusMiddleware())

	api := router.Group("/api")
	srv.SetupUserHandlersFor(api)

	router.Static("/static", "./static")

	router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/static/index.html")
	})

	if err := http.Serve(conf.Listener, router); err != nil {
		log.Fatalln(err)
	}
}
