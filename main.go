package main

import (
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pocockn/shouts-api/config"
	"github.com/pocockn/shouts-api/persistance"
	"github.com/pocockn/shouts-api/shouts/handler"
	"github.com/pocockn/shouts-api/shouts/repository"
	"log"
	"net/http"
)

func main() {
	config := config.NewConfig()

	db, err := persistance.NewConnection(config)
	if err != nil {
		log.Fatal(err)
	}

	connection, err := db.Connect()
	if err != nil {
		log.Fatal(err)
	}

	e := echo.New()

	//TODO: this should come from config.
	e.Use(
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: []string{"http://localhost:8081"},
			AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		}),
		middleware.Logger(),
	)

	shoutRepo := repository.NewShoutsRepository(connection)
	handler.NewShoutHandler(config, e, shoutRepo)

	e.Logger.Fatal(e.Start(":1323"))
}
