package main

import (
	"fmt"

	"github.com/gcarrenho/guidemysteps/src/internal/adapters/handlers"
	osmrepository "github.com/gcarrenho/guidemysteps/src/internal/adapters/repositories/open_street_map"
	postgresql "github.com/gcarrenho/guidemysteps/src/internal/adapters/repositories/postgresql"
	"github.com/gcarrenho/guidemysteps/src/internal/core/services"

	i18nRepo "github.com/gcarrenho/guidemysteps/src/internal/adapters/repositories/i18"

	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var (
	bundle = i18n.NewBundle(language.English)
)

// main executes application
func main() {
	InitRoutes()
}

func InitRoutes() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	group := router.Group("/guide-my-steps/")

	//Init db
	postgreConf := postgresql.NewPostgreSQLConf()
	postgreSQLDB, err := postgreConf.InitPostgreSQLDB()
	if err != nil {
		fmt.Println(err)
	}

	// repositories
	osm := osmrepository.NewRoutingRepo("https://routing.openstreetmap.de")
	translate := i18nRepo.NewI18nRepo(bundle)
	userRepo := postgresql.NewUserRepository(postgreSQLDB)

	// services
	routingSvc := services.NewRoutingSvc(osm, translate)
	userSvc := services.NewUserSvc(userRepo)
	handlers.NewRoutingHandler(group, routingSvc, userSvc)

	fmt.Println("Listening and serving HTTP on :8080")
	err = router.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
