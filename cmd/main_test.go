package main

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/exasol/exasol-driver-go"
	"github.com/exasol/extension-manager/pkg/extensionController"
	"github.com/sirupsen/logrus"
)

type MyFormatter struct{}

func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(fmt.Sprintf("%s: %s\n", entry.Level, entry.Message)), nil
}

// Run this test with `go test -v ./cmd/...`

func TestMain(m *testing.M) {
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(new(MyFormatter))
	var db *sql.DB = createDBConnection()
	t0 := time.Now()
	config := extensionController.ExtensionManagerConfig{
		ExtensionRegistryURL: "https://d3d6d68cbkri8h.cloudfront.net/registry.json", // test
		//ExtensionRegistryURL: "https://dpgtxmuoxieps.cloudfront.net/registry.json", // prod
		//BucketFSBasePath: "/buckets/bfssaas",
		BucketFSBasePath: "/buckets/",
		ExtensionSchema:  "CHP_EM_TESTING",
	}
	ctrl, err := extensionController.CreateWithValidatedConfig(config)
	if err != nil {
		panic("Error creating controller: " + err.Error())
	}

	fmt.Printf("Fetching extensions...\n")
	extensions, err := ctrl.GetAllExtensions(context.Background(), db)
	if err != nil {
		panic("Error getting extensions: " + err.Error())
	}
	for _, extension := range extensions {
		fmt.Printf("- %v\n", extension)
	}
	fmt.Printf("Listing installed extensions...\n")
	installed, err := ctrl.GetInstalledExtensions(context.Background(), db)
	if err != nil {
		panic("Error getting installed extensions: " + err.Error())
	}
	fmt.Printf("Found %d installed extensions in %dms: %v\n", len(installed), time.Since(t0).Milliseconds(), installed)
}

func createDBConnection() *sql.DB {
	// Connect to the database
	fmt.Printf("Connecting to Exasol database...\n")
	db, err := sql.Open("exasol", exasol.
		NewConfigWithRefreshToken("exa_pat_c6Dsqr6O6MlKo6rTIhqdb9aQiBm6sqdJtRfidrwC3ukaIa").
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
