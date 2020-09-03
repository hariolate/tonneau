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

//package main
//
//import (
//	"flag"
//)
//
//var (
//	configFile = flag.String("config", "config.json", "specify config file")
//	debug      = flag.Bool("debug", true, "is in debug env or not")
//)

//type User struct {
//	gorm.Model
//	CreditCard CreditCard
//}
//
//type CreditCard struct {
//	gorm.Model
//	Number string
//	UserID uint
//}
//
//type MatchResult struct {
//	gorm.Model
//
//	Players []User
//	Scores  []int
//
//	RoundCount int
//}
//
//type User struct {
//	gorm.Model
//	UserProfile UserProfile
//	Email       string `gorm:"unique" validate:"required,email"`
//	Password    string `validate:"required,min=3"`
//}
//
//type UserProfile struct {
//	gorm.Model
//	Alias    string `json:"alias"`
//	Picture  []byte `json:"picture"`
//	Trophies int    `json:"trophies"`
//	UserID uint
//}
//
//func main() {
//	flag.Parse()
//	conf := config.MustFromFile(*configFile).MustParse()
//	conf.DBConn.AutoMigrate(&User{}, &UserProfile{})
//	var u User
//	conf.DBConn.Create(&u)
//}
