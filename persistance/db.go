package persistance

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pocockn/shouts-api/config"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

type (
	// GormDB holds a database connection.
	GormDB struct {
		maxConnections int
		url            string
	}
)

// NewConnection creates a new connection for the database.
func NewConnection(config config.Config) (*GormDB, error) {
	return &GormDB{
		maxConnections: config.Database.MaxConnections,
		url:            generateURL(config),
	}, nil
}

// Connect connects to the database and passes back the connection so we can
// use it throughout the application
func (g GormDB) Connect() (*gorm.DB, error) {
	gormDb, err := gorm.Open("mysql", g.url)
	if err != nil {
		return nil, fmt.Errorf("Unable to open DB connection using GORM: %s", err)
	}

	maxConnsPerContainer := g.maxConnections / 4
	gormDb.DB().SetMaxOpenConns(maxConnsPerContainer / 2)

	return gormDb, nil
}

func generateURL(config config.Config) string {
	templateString := "%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4"

	return fmt.Sprintf(
		templateString,
		config.Database.Username,
		config.Database.Password,
		config.Database.Host,
		config.Database.Port,
		config.Database.DatabaseName,
	)
}
