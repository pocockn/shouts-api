package handler

import (
	"net/http"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/pocockn/shouts-api/shouts"
)

// ShoutHandler implements the repository interface.
type ShoutHandler struct {
	Repo shouts.Repository
}

// NewShoutHandler creates a new shouts handler with the routes.
func NewShoutHandler(e *echo.Echo, repo shouts.Repository) {
	handler := &ShoutHandler{
		Repo: repo,
	}

	e.GET("/shouts/:id", handler.FetchShout)
	e.GET("/shouts", handler.FetchAll)
}

// FetchShout gets a shout from it's ID.
func (s *ShoutHandler) FetchShout(c echo.Context) error {
	idP, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusNotFound, nil)
	}

	id := uint(idP)

	shout, err := s.Repo.Fetch(id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, shout)
}

// FetchAll fetches all the shouts from the DB.
func (s *ShoutHandler) FetchAll(c echo.Context) error {
	allShouts, err := s.Repo.FetchAll()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	return c.JSON(http.StatusOK, allShouts)
}

// Store takes a shout and stores it in the DB.
func (s *ShoutHandler) Store(c echo.Context) error {
	// move all this into controllers when refactoring
	// get the two images
	sourceFile, err := c.FormFile("source")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}
	src, err := sourceFile.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, nil)
	}
	defer src.Close()

	// upload both images to s3
	return nil
	// store their s3 key name (file name) on the shout object

	// perform the lambda function analysis
}
