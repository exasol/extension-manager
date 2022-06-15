package restAPI

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"sync"

	cont "github.com/exasol/extension-manager/extensionController"

	"github.com/exasol/exasol-driver-go"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"

	// docs are generated by Swag CLI, you have to import it.
	_ "github.com/exasol/extension-manager/generatedApiDocs"
)

// RestAPI is the interface that provides the REST API server of the extension-manager.
type RestAPI interface {
	// Serve starts the server. This method blocks until the server is stopped or fails.
	Serve()
	// Stop stops the server
	Stop()
}

// @title           Exasol extension manager REST API
// @version         0.1.0
// @description     This is a REST API for managing extensions in an Exasol database.

// @contact.name   Exasol Integration team
// @contact.email  opensource@exasol.com

// @license.name  MIT
// @license.url   https://github.com/exasol/extension-manager/blob/main/LICENSE

// @BasePath  /

func Create(controller cont.ExtensionController) RestAPI {
	return &restAPIImpl{Controller: controller}
}

type restAPIImpl struct {
	Controller   cont.ExtensionController
	server       *http.Server
	stopped      *bool
	stoppedMutex *sync.Mutex
}

func (restApi *restAPIImpl) Serve() {
	if restApi.server != nil {
		panic("server already running")
	}
	restApi.setStopped(false)
	router := gin.Default()
	router.GET("/extensions", restApi.handleGetExtensions)
	router.GET("/installations", restApi.handleGetInstallations)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	const port = "8080"
	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	restApi.server = srv
	log.Printf("Starting server on port %s...\n", port)
	err := restApi.server.ListenAndServe() // blocking
	if err != nil && !restApi.isStopped() {
		panic(fmt.Sprintf("failed to start rest API server. Cause: %v", err.Error()))
	}
}

func (restApi *restAPIImpl) setStopped(stopped bool) {
	if restApi.stopped == nil {
		stopped := false
		restApi.stopped = &stopped
		restApi.stoppedMutex = &sync.Mutex{}
	}
	restApi.stoppedMutex.Lock()
	defer restApi.stoppedMutex.Unlock()
	*restApi.stopped = stopped
}

func (restApi *restAPIImpl) isStopped() bool {
	restApi.stoppedMutex.Lock()
	defer restApi.stoppedMutex.Unlock()
	return *restApi.stopped
}

func (restApi *restAPIImpl) Stop() {
	if restApi.server == nil {
		panic("cant stop server since it's not running")
	}
	restApi.setStopped(true)
	err := restApi.server.Shutdown(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to shutdown rest API server. Cause: %v", err.Error()))
	}
	restApi.server = nil
}

// @Summary      Get all extensions
// @Description  Get a list of all available extensions.
// @Produce      json
// @Success      200  {object}  ExtensionsResponse
// @Param 	     dbHost query string true "Hostname of the Exasol DB to manage"
// @Param 	     dbPort query int true "port number of the Exasol DB to manage"
// @Param 	     dbUser query string true "username of the Exasol DB to manage"
// @Param 	     dbPass query string true "password of the Exasol DB to manage"
// @Failure      500  {object}  string
// @Router       /extensions [get]
func (restApi *restAPIImpl) handleGetExtensions(c *gin.Context) {
	response, err := restApi.getExtensions(c)
	restApi.sendResponse(c, response, err)
}

func (restApi *restAPIImpl) getExtensions(c *gin.Context) (*ExtensionsResponse, error) {
	dbConnectionWithNoAutocommit, err := restApi.openDBConnection(c)
	if err != nil {
		return nil, err
	}
	defer closeDbConnection(dbConnectionWithNoAutocommit)
	extensions, err := restApi.Controller.GetAllExtensions(dbConnectionWithNoAutocommit)
	if err != nil {
		return nil, err
	}
	convertedExtensions := make([]ExtensionsResponseExtension, 0, len(extensions))
	for _, extension := range extensions {
		convertedExtensions = append(convertedExtensions, ExtensionsResponseExtension{Name: extension.Name, Description: extension.Description, InstallableVersions: extension.InstallableVersions})
	}
	response := ExtensionsResponse{
		Extensions: convertedExtensions,
	}
	return &response, nil
}

type ExtensionsResponse struct {
	Extensions []ExtensionsResponseExtension `json:"extensions"`
}

type ExtensionsResponseExtension struct {
	Name                string   `json:"name"`
	Description         string   `json:"description"`
	InstallableVersions []string `json:"installableVersions"`
}

// @Summary      Get all installations
// @Description  Get a list of all installations. Installation means, that an extension is installed in the database (e.g. JAR files added to BucketFS, Adapter Script created).
// @Produce      json
// @Success      200  {object}  InstallationsResponse
// @Param 	     dbHost query string true "Hostname of the Exasol DB to manage"
// @Param 	     dbPort query int true "port number of the Exasol DB to manage"
// @Param 	     dbUser query string true "username of the Exasol DB to manage"
// @Param 	     dbPass query string true "password of the Exasol DB to manage"
// @Failure      500  {object}  string
// @Router       /installations [get]
func (restApi *restAPIImpl) handleGetInstallations(c *gin.Context) {
	response, err := restApi.getInstallations(c)
	restApi.sendResponse(c, response, err)
}

func (restApi *restAPIImpl) getInstallations(c *gin.Context) (*InstallationsResponse, error) {
	dbConnection, err := restApi.openDBConnection(c)
	if err != nil {
		return nil, err
	}
	defer closeDbConnection(dbConnection)
	installations, err := restApi.Controller.GetAllInstallations(dbConnection)
	if err != nil {
		return nil, err
	}
	convertedInstallations := make([]InstallationsResponseInstallation, 0, len(installations))
	for _, installation := range installations {
		convertedInstallations = append(convertedInstallations, InstallationsResponseInstallation{installation.Name, installation.Version, installation.InstanceParameters})
	}
	response := InstallationsResponse{
		Installations: convertedInstallations,
	}
	return &response, nil
}

func (restApi *restAPIImpl) sendResponse(c *gin.Context, response interface{}, err error) {
	if err != nil {
		c.String(500, "Internal error.")
		fmt.Println(err.Error())
		return
	}
	c.JSON(200, response)
}

func closeDbConnection(database *sql.DB) {
	err := database.Close()
	if err != nil {
		// Strange but not critical. So we just log it and go on.
		fmt.Printf("failed to close db connection. Cause %v", err.Error())
	}
}

func (restApi *restAPIImpl) openDBConnection(c *gin.Context) (*sql.DB, error) {
	config, err := getDbConfig(c)
	if err != nil {
		return nil, fmt.Errorf("failed to get db config: %w", err)
	}
	config.Autocommit(false).ValidateServerCertificate(false)
	database, err := sql.Open("exasol", config.String())
	if err != nil {
		return nil, fmt.Errorf("failed to open a database connection. Cause: %w", err)
	}
	_, err = database.Exec("select 1")
	if err != nil {
		return nil, fmt.Errorf("database connection test failed. Cause: %v", err.Error())
	}
	return database, nil
}

func getDbConfig(c *gin.Context) (*exasol.DSNConfigBuilder, error) {
	query := c.Request.URL.Query()
	host := query.Get("dbHost")
	if host == "" {
		return nil, fmt.Errorf("missing parameter dbHost")
	}
	portString := query.Get("dbPort")
	if portString == "" {
		return nil, fmt.Errorf("missing parameter dbPort")
	}
	port, err := strconv.Atoi(portString)
	if err != nil {
		return nil, fmt.Errorf("invalid value %q for parameter dbPort", portString)
	}
	user := query.Get("dbUser")
	if user == "" {
		return nil, fmt.Errorf("missing parameter dbUser")
	}
	password := query.Get("dbPass")
	if password == "" {
		return nil, fmt.Errorf("missing parameter dbPass")
	}
	config := exasol.NewConfig(user, password).Port(port).Host(host)
	return config, nil
}

type InstallationsResponse struct {
	Installations []InstallationsResponseInstallation `json:"installations"`
}

type InstallationsResponseInstallation struct {
	Name               string        `json:"name"`
	Version            string        `json:"version"`
	InstanceParameters []interface{} `json:"instanceParameters"`
}
