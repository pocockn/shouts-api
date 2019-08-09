package handler

import (
	"encoding/json"
	"github.com/labstack/gommon/log"
	awsWrappersSNS "github.com/pocockn/awswrappers/sns"
	"github.com/pocockn/models/sns"
	"github.com/pocockn/models/sns/lambda"
	"github.com/pocockn/shouts-api/config"
	"github.com/pocockn/shouts-api/services"
	uuid "github.com/satori/go.uuid"
	"net/http"
	"strconv"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/labstack/echo"
	"github.com/pocockn/shouts-api/shouts"
)

// ShoutHandler implements the repository interface.
type ShoutHandler struct {
	Config        config.Config
	Repo          shouts.Repository
	UploadService services.Upload
}

// NewShoutHandler creates a new shouts handler with the routes.
func NewShoutHandler(config config.Config, e *echo.Echo, repo shouts.Repository) {
	handler := &ShoutHandler{
		Config:        config,
		Repo:          repo,
		UploadService: services.NewUpload(config, config.S3.Client),
	}

	e.GET("/shouts/:id", handler.FetchShout)
	e.GET("/shouts", handler.FetchAll)
	e.POST("/shout", handler.Store)
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
	guid := uuid.NewV4().String()
	sourceFile, err := c.FormFile("source")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	targetFile, err := c.FormFile("target")
	if err != nil {
		return c.JSON(http.StatusInternalServerError, nil)
	}

	err = s.UploadService.MultiUpload(sourceFile, targetFile)
	if err != nil {
		return err
	}

	imageSimilarityPayload, err := lambda.NewRawJsonImageSimilarity(
		targetFile.Filename,
		sourceFile.Filename,
		guid,
	)
	if err != nil {
		return err
	}

	snsMessage := sns.Message{
		ID:      guid,
		Payload: &imageSimilarityPayload,
	}

	snsMessagePayload, err := json.Marshal(&snsMessage)
	if err != nil {
		return err
	}

	snsClient := *awsWrappersSNS.NewClient(nil, false, nil)
	messageID, err := snsClient.PublishMessage(
		string(snsMessagePayload),
		s.Config.SNS.Arn,
	)
	if err != nil {
		return err
	}

	log.Infof("SNS notification %s", messageID)

	return c.JSON(http.StatusCreated, nil)
}
