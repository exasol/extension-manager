package main

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/exasol/exasol-driver-go"
	"github.com/exasol/extension-manager/pkg/extensionController"
	"github.com/sirupsen/logrus"
)

type MyFormatter struct{}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s: %s\n", entry.Level, entry.Message)), nil
}

func TestMain(m *testing.M) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(new(MyFormatter))
	config := extensionController.ExtensionManagerConfig{
		ExtensionRegistryURL: "https://dpgtxmuoxieps.cloudfront.net/registry.json",
		BucketFSBasePath:     "/buckets/bfssaas",
		ExtensionSchema:      "EXA_EXTENSIONS",
	}
	// Create controller and handle configuration validation error
	ctrl, err := extensionController.CreateWithValidatedConfig(config)
	if err != nil {
		panic("Error creating controller: " + err.Error())
	}

	// Create database connection (required as an argument for all controller methods)
	var db *sql.DB = createDBConnection()
	fmt.Printf("Fetching extensions...\n")
	extensions, err := ctrl.GetAllExtensions(context.Background(), db)
	if err != nil {
		panic("Error getting extensions: " + err.Error())
	}
	for _, extension := range extensions {
		fmt.Printf("- %v\n", extension)
	}
}

func createDBConnection() *sql.DB {
	// Connect to the database
	fmt.Printf("Connecting to Exasol database...\n")
	db, err := sql.Open("exasol", exasol.
		NewConfigWithRefreshToken("exa_pat_c6Dsqr6O6MlKo6rTIhqdb9aQiBm6sqdJtRfidrwC3ukaIa").
		//NewConfig("christoph_pirkl", "exa_pat_c6Dsqr6O6MlKo6rTIhqdb9aQiBm6sqdJtRfidrwC3ukaIa").
		Host("i2whdqt7vzhitldfu4sjuse2s4.clusters.exasol.com").
		Port(8563).
		Autocommit(false).
		String())
	if err != nil {
		panic("Error connecting to database: " + err.Error())
	}
	testConnection(db)
	return db
}

func testConnection(db *sql.DB) {
	row := db.QueryRow("SELECT 'a'")
	var value sql.NullString
	err := row.Scan(&value)
	if err != nil {
		panic("Error scanning row: " + err.Error())
	}
	fmt.Printf("Connected to Exasol database: %v\n", value.String)
}
