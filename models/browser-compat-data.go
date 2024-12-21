package models

import (
	"log"
	"os"
	"time"

	"gorm.io/datatypes"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Meta struct {
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

type Status struct {
	Deprecated    bool `json:"deprecated"`
	Experimental  bool `json:"experimental"`
	StandardTrack bool `json:"standard_track"`
}

type BrowserCompatData struct {
	Type           string                      `json:"type" gorm:"index"`
	MdnUrl         string                      `json:"mdn_url"`
	SpecUrl        string                      `json:"spec_url"`
	Api            string                      `json:"api" gorm:"index:api_browser_version,unique"`
	Browser        string                      `json:"Browser" gorm:"index:api_browser_version,unique"`
	BrowserVersion string                      `json:"BrowserVersion" gorm:"index:api_browser_version,unique"`
	Tags           datatypes.JSONSlice[string] `json:"tags"`
	Status         datatypes.JSONType[Status]  `json:"status"`
}

var db *gorm.DB

func init() {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  false,         // Disable color
		},
	)
	var err error
	db, err = gorm.Open(sqlite.Open("./browser-compat.db"), &gorm.Config{Logger: newLogger})
	if err != nil {
		panic("failed to connect database")
	}

	MigrateDatabase()
}

// MigrateDatabase migrates the BrowserCompatData model to the database.
func MigrateDatabase() {
	// Migrate the schema
	db.AutoMigrate(&BrowserCompatData{})
}

// ImportData imports browser compatibility data into the database.
func ImportData(data []BrowserCompatData) {
	for _, d := range data {
		db.Create(&d)
	}
}

func Create(data BrowserCompatData) error {
	err := db.Create(&data).Error
	// if err != nil {
	// 	log.Printf("data: %#v", data)
	// }
	return err
}

// QueryData queries browser compatibility data from the database.
func QueryData(query interface{}, args ...interface{}) []BrowserCompatData {
	var results []BrowserCompatData
	db.Where(query, args...).Find(&results)
	return results
}
