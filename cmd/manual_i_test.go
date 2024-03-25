package main

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/exasol/exasol-driver-go"
	"github.com/exasol/extension-manager/pkg/extensionController"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/suite"
)

// Create file manual-test.properties, add configuration and
// run this test with `go test -v ./cmd/...`

type ManualITestSuite struct {
	suite.Suite
	config configProperties
}

func TestManualITestSuite(t *testing.T) {
	suite.Run(t, new(ManualITestSuite))
}

func (suite *ManualITestSuite) SetupSuite() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(new(simpleFormatter))
	config, err := readPropertiesFile("../manual-test.properties")
	if err != nil {
		suite.T().Skip("Skipping manual integration tests: " + err.Error())
	}
	suite.config = config

}

func (suite *ManualITestSuite) TestListInstalledExtensions() {
	var db *sql.DB = suite.createDBConnection()
	t0 := time.Now()
	ctrl, err := extensionController.CreateWithValidatedConfig(suite.createExtensionManagerConfig())
	if err != nil {
		suite.FailNow("Error creating controller: " + err.Error())
	}
	extensions, err := ctrl.GetAllExtensions(context.Background(), db)
	if err != nil {
		suite.FailNow("Error getting extensions: " + err.Error())
	}
	log.Infof("Found %d extensions, listing installed extensions...", len(extensions))
	installed, err := ctrl.GetInstalledExtensions(context.Background(), db)
	if err != nil {
		suite.FailNow("Error getting installed extensions: " + err.Error())
	}
	log.Infof("Total duration %dms: %v", time.Since(t0).Milliseconds(), installed)
}

func (suite *ManualITestSuite) createExtensionManagerConfig() extensionController.ExtensionManagerConfig {
	return extensionController.ExtensionManagerConfig{
		ExtensionRegistryURL: suite.getConfigValue("extensionRegistryURL"),
		BucketFSBasePath:     suite.getConfigValue("bucketFSBasePath"),
		ExtensionSchema:      suite.getConfigValue("extensionSchema"),
	}
}

func (suite *ManualITestSuite) createDBConnection() *sql.DB {
	dbHost := suite.getConfigValue("databaseHost")
	log.Debugf("Connecting to database at %q...", dbHost)
	db, err := sql.Open("exasol", exasol.
		NewConfigWithRefreshToken(suite.getConfigValue("databaseToken")).
		Host(dbHost).
		Port(8563).
		Autocommit(false).
		String())
	if err != nil {
		suite.FailNow("Error connecting to database: " + err.Error())
	}
	suite.testConnection(db)
	return db
}

func (suite *ManualITestSuite) testConnection(db *sql.DB) {
	row := db.QueryRow("SELECT 'a'")
	var value sql.NullString
	err := row.Scan(&value)
	if err != nil {
		suite.FailNow("Error scanning row: " + err.Error())
	}
}

func (suite *ManualITestSuite) getConfigValue(key string) string {
	value := suite.config[key]
	if value == "" {
		suite.FailNow(fmt.Sprintf("Key %q not found in config file", key))
	}
	return value
}

type configProperties map[string]string

// readPropertiesFile reads a Java properties file and returns a map of key-value pairs.
// Based on https://stackoverflow.com/a/46860900
func readPropertiesFile(filename string) (configProperties, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}
	config := configProperties{}
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error opening file %q: %w", path, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "#") {
			continue
		}
		if equal := strings.Index(line, "="); equal >= 0 {
			if key := strings.TrimSpace(line[:equal]); len(key) > 0 {
				value := ""
				if len(line) > equal {
					value = strings.TrimSpace(line[equal+1:])
				}
				config[key] = value
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return config, nil
}
