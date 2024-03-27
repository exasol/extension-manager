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
	"github.com/exasol/extension-manager/pkg/extensionAPI"
	"github.com/exasol/extension-manager/pkg/extensionController"

	log "github.com/sirupsen/logrus"

	"github.com/stretchr/testify/suite"
)

// Create file manual-test.properties, add configuration and
// run this test with `go test -v ./cmd/...`

type ManualITestSuite struct {
	suite.Suite
	config configProperties
	db     *sql.DB
	ctrl   extensionController.TransactionController // Add this line
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
	suite.db = suite.createDBConnection()
	suite.ctrl = suite.createController()
}

func (suite *ManualITestSuite) TearDownSuite() {
	suite.NoError(suite.db.Close())
}

func (suite *ManualITestSuite) createController() extensionController.TransactionController {
	ctrl, err := extensionController.CreateWithValidatedConfig(suite.createExtensionManagerConfig())
	if err != nil {
		suite.FailNow("Error creating controller: " + err.Error())
	}
	return ctrl
}

func (suite *ManualITestSuite) createExtensionManagerConfig() extensionController.ExtensionManagerConfig {
	return extensionController.ExtensionManagerConfig{
		ExtensionRegistryURL: suite.getConfigValue("extensionRegistryURL"),
		BucketFSBasePath:     suite.getConfigValue("bucketFSBasePath"),
		ExtensionSchema:      suite.getConfigValue("extensionSchema"),
	}
}

func (suite *ManualITestSuite) TestListInstalledExtensions() {
	t0 := time.Now()
	extensions := suite.getAllExtensions()
	suite.GreaterOrEqual(len(extensions), 6, "Expected at least six extension")
	log.Infof("Found %d extensions, listing installed extensions...", len(extensions))
	installed := suite.getInstalledExtensions()
	suite.GreaterOrEqual(len(installed), 1, "Expected at least one installations")
	duration := time.Since(t0)
	log.Infof("Total duration %dms: %v", duration.Milliseconds(), installed)
	suite.LessOrEqual(duration.Milliseconds(), int64(2000), "Process took too long")
}

func (suite *ManualITestSuite) getAllExtensions() []*extensionController.Extension {
	extensions, err := suite.ctrl.GetAllExtensions(context.Background(), suite.db)
	if err != nil {
		suite.FailNow("Error getting extensions: " + err.Error())
	}
	return extensions
}

func (suite *ManualITestSuite) getInstalledExtensions() []*extensionAPI.JsExtInstallation {
	installed, err := suite.ctrl.GetInstalledExtensions(context.Background(), suite.db)
	if err != nil {
		suite.FailNow("Error getting installed extensions: " + err.Error())
	}
	return installed
}

func (suite *ManualITestSuite) TestInstallAndUninstallExtension() {
	t0 := time.Now()
	extensions := suite.getAllExtensions()
	ext := extensions[0]
	log.Infof("Installing extension %q...", ext.Id)
	err := suite.ctrl.InstallExtension(context.Background(), suite.db, ext.Id, ext.InstallableVersions[0].Name)
	if err != nil {
		suite.FailNow("Error installing extension: " + err.Error())
	}

	err = suite.ctrl.UninstallExtension(context.Background(), suite.db, ext.Id, ext.InstallableVersions[0].Name)
	if err != nil {
		suite.FailNow("Error uninstalling extension: " + err.Error())
	}
	duration := time.Since(t0)
	log.Infof("Total duration %dms", duration.Milliseconds())
	suite.LessOrEqual(duration.Milliseconds(), int64(2500), "Process took too long")
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
