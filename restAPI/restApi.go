package restAPI

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"sync"

	cont "github.com/exasol/extension-manager/extensionController"
	"github.com/exasol/extension-manager/restAPI/models"
	log "github.com/sirupsen/logrus"

	"github.com/exasol/exasol-driver-go"
	"github.com/gin-gonic/gin"
)

// RestAPI is the interface that provides the REST API server of the extension-manager.
type RestAPI interface {
	// Serve starts the server. This method blocks until the server is stopped or fails.
	Serve()
	// Stop stops the server
	Stop()
}

// Description for swagger must be at the end of the file!

// @title          Exasol extension manager REST API
// @version        0.1.0
// @contact.name   Exasol Integration team
// @contact.email  opensource@exasol.com
// @license.name   MIT
// @license.url    https://github.com/exasol/extension-manager/blob/main/LICENSE
// @host localhost:8080
// @BasePath  /
// @accept   json
// @produce  json
// @query.collection.format csv
// @schemes http https
// @tag.name extension
// @tag.description List, install and uninstall extensions
// @tag.name instance
// @tag.description List, create and remove instances of an extension

// Create creates a new RestAPI.
func Create(controller cont.TransactionController, serverAddress string) RestAPI {
	return &restAPIImpl{controller: controller, serverAddress: serverAddress}
}

type restAPIImpl struct {
	controller    cont.TransactionController
	serverAddress string
	server        *http.Server
	stopped       *bool
	stoppedMutex  *sync.Mutex
}

func (api *restAPIImpl) Serve() {
	if api.server != nil {
		panic("server already running")
	}
	api.setStopped(false)

	handler, _, err := setupStandaloneAPI(api.controller)
	if err != nil {
		log.Fatalf("failed to setup api: %v", err)
	}
	api.server = &http.Server{
		Addr:    api.serverAddress,
		Handler: handler,
	}
	log.Printf("Starting server on %s...\n", api.serverAddress)
	err = api.server.ListenAndServe() // blocking
	if err != nil && !api.isStopped() {
		log.Fatalf("failed to start server: %v", err)
	}
}

func (api *restAPIImpl) setStopped(stopped bool) {
	if api.stopped == nil {
		stopped := false
		api.stopped = &stopped
		api.stoppedMutex = &sync.Mutex{}
	}
	api.stoppedMutex.Lock()
	defer api.stoppedMutex.Unlock()
	*api.stopped = stopped
}

func (api *restAPIImpl) isStopped() bool {
	api.stoppedMutex.Lock()
	defer api.stoppedMutex.Unlock()
	return *api.stopped
}

func (api *restAPIImpl) Stop() {
	if api.server == nil {
		panic("cant stop server since it's not running")
	}
	api.setStopped(true)
	err := api.server.Shutdown(context.Background())
	if err != nil {
		panic(fmt.Sprintf("failed to shutdown rest API server. Cause: %v", err))
	}
	api.server = nil
}

// @Summary      Get all extensions
// @Description  Get a list of all available extensions, i.e. extensions that can be installed.
// @Id           getAvailableExtensions
// @tags         extension
// @Produce      json
// @Success      200 {object} ExtensionsResponse
// @Param        dbHost query string true "Hostname of the Exasol DB to manage"
// @Param        dbPort query int true "Port number of the Exasol DB to manage"
// @Param        dbUser query string false "Username of the Exasol DB to manage"
// @Param        dbPassword query string false "Password of the Exasol DB to manage"
// @Param        dbAccessToken query string false "Access token of the Exasol DB to manage"
// @Param        dbRefreshToken query string false "Refresh token of the Exasol DB to manage"
// @Failure      500 {object} string
// @Router       /extensions [get]
func (api *restAPIImpl) handleGetExtensions(c *gin.Context) {
	response, err := api.getExtensions(c)
	api.sendResponse(c, response, err)
}

func (api *restAPIImpl) getExtensions(c *gin.Context) (*models.ExtensionsResponse, error) {
	db, err := api.openDBConnection(c)
	if err != nil {
		return nil, err
	}
	defer closeDbConnection(db)
	extensions, err := api.controller.GetAllExtensions(c, db)
	if err != nil {
		return nil, err
	}
	convertedExtensions := make([]models.ExtensionsResponseExtension, 0, len(extensions))
	for _, extension := range extensions {
		ext := models.ExtensionsResponseExtension{Id: extension.Id, Name: extension.Name, Description: extension.Description, InstallableVersions: extension.InstallableVersions}
		convertedExtensions = append(convertedExtensions, ext)
	}
	response := models.ExtensionsResponse{
		Extensions: convertedExtensions,
	}
	return &response, nil
}

// @Summary      Get all installed extensions.
// @Description  Get a list of all installed extensions. Installation means, that an extension is installed in the database (e.g. JAR files added to BucketFS and Adapter Script created).
// @Id           getInstalledExtensions
// @tags         extension
// @Produce      json
// @Success      200 {object} InstallationsResponse
// @Param        dbHost query string true "Hostname of the Exasol DB to manage"
// @Param        dbPort query int true "Port number of the Exasol DB to manage"
// @Param        dbUser query string false "Username of the Exasol DB to manage"
// @Param        dbPassword query string false "Password of the Exasol DB to manage"
// @Param        dbAccessToken query string false "Access token of the Exasol DB to manage"
// @Param        dbRefreshToken query string false "Refresh token of the Exasol DB to manage"
// @Failure      500 {object} string
// @Router       /installations [get]
func (api *restAPIImpl) handleGetInstallations(c *gin.Context) {
	response, err := api.getInstallations(c)
	api.sendResponse(c, response, err)
}

func (api *restAPIImpl) getInstallations(c *gin.Context) (*models.InstallationsResponse, error) {
	db, err := api.openDBConnection(c)
	if err != nil {
		return nil, err
	}
	defer closeDbConnection(db)
	installations, err := api.controller.GetAllInstallations(c, db)
	if err != nil {
		return nil, err
	}
	convertedInstallations := make([]models.InstallationsResponseInstallation, 0, len(installations))
	for _, installation := range installations {
		convertedInstallations = append(convertedInstallations, models.InstallationsResponseInstallation{installation.Name, installation.Version, installation.InstanceParameters})
	}
	response := models.InstallationsResponse{
		Installations: convertedInstallations,
	}
	return &response, nil
}

// @Summary      Install an extension.
// @Description  This installs an extension in a given version, e.g. by creating Adapter Scripts.
// @Id           installExtension
// @tags         extension
// @Produce      json
// @Success      200 {object} string
// @Param        dbHost query string true "Hostname of the Exasol DB to manage"
// @Param        dbPort query int true "Port number of the Exasol DB to manage"
// @Param        dbUser query string false "Username of the Exasol DB to manage"
// @Param        dbPassword query string false "Password of the Exasol DB to manage"
// @Param        dbAccessToken query string false "Access token of the Exasol DB to manage"
// @Param        dbRefreshToken query string false "Refresh token of the Exasol DB to manage"
// @Param        extensionId query string true "ID of the extension to install"
// @Param        extensionVersion query string true "Version of the extension to install"
// @Param        dummy body string false "dummy body" default()
// @Failure      500 {object} string
// @Router       /installations [put]
func (api *restAPIImpl) handlePutInstallation(c *gin.Context) {
	result, err := api.installExtension(c)
	api.sendResponse(c, result, err)
}

func (api *restAPIImpl) installExtension(c *gin.Context) (string, error) {
	db, err := api.openDBConnection(c)
	if err != nil {
		return "", err
	}
	defer closeDbConnection(db)
	query := c.Request.URL.Query()
	extensionId := query.Get("extensionId")
	if extensionId == "" {
		return "", fmt.Errorf("missing parameter extensionId")
	}
	extensionVersion := query.Get("extensionVersion")
	if extensionVersion == "" {
		return "", fmt.Errorf("missing parameter extensionVersion")
	}

	err = api.controller.InstallExtension(c, db, extensionId, extensionVersion)

	if err != nil {
		return "", fmt.Errorf("error installing extension: %v", err)
	}
	return "", nil
}

// handlePutInstance creates a new instance.
// @Summary      Create an instance of an extension.
// @Description  This creates a new instance of an extension, e.g. a virtual schema.
// @Id           createInstance
// @tags         instance
// @Produce      json
// @Success      200 {object} CreateInstanceResponse
// @Param        dbHost query string true "Hostname of the Exasol DB to manage"
// @Param        dbPort query int true "Port number of the Exasol DB to manage"
// @Param        dbUser query string false "Username of the Exasol DB to manage"
// @Param        dbPassword query string false "Password of the Exasol DB to manage"
// @Param        dbAccessToken query string false "Access token of the Exasol DB to manage"
// @Param        dbRefreshToken query string false "Refresh token of the Exasol DB to manage"
// @Param        createInstanceRequest body CreateInstanceRequest true "Request data for creating an instance"
// @Failure      500 {object} string
// @Router       /instances [put]
func (api *restAPIImpl) handlePutInstance(c *gin.Context) {
	result, err := api.createInstance(c)
	api.sendResponse(c, result, err)
}

func (api *restAPIImpl) createInstance(c *gin.Context) (CreateInstanceResponse, error) {
	db, err := api.openDBConnection(c)
	var response CreateInstanceResponse
	if err != nil {
		return response, err
	}
	defer closeDbConnection(db)
	var request CreateInstanceRequest
	if err := c.BindJSON(&request); err != nil {
		return response, fmt.Errorf("invalid request: %w", err)
	}

	var parameters []cont.ParameterValue
	for _, p := range request.ParameterValues {
		parameters = append(parameters, cont.ParameterValue{Name: p.Name, Value: p.Value})
	}
	response.InstanceName, err = api.controller.CreateInstance(c, db, request.ExtensionId, request.ExtensionVersion, parameters)
	if err != nil {
		return response, fmt.Errorf("error installing extension: %v", err)
	}
	return response, nil
}

// @Description Request data for creating a new instance of an extension.
type CreateInstanceRequest struct {
	ExtensionId      string           `json:"extensionId"`      // The ID of the extension
	ExtensionVersion string           `json:"extensionVersion"` // The version of the extension
	ParameterValues  []ParameterValue `json:"parameterValues"`  // The parameters for the new instance
}

// @Description Parameter values for creating a new instance.
type ParameterValue struct {
	Name  string `json:"name"`  // The name of the parameter
	Value string `json:"value"` // The value of the parameter
}

// @Description Response data for creating a new instance of an extension.
type CreateInstanceResponse struct {
	InstanceName string `json:"instanceName"` // The name of the newly created instance
}

func (api *restAPIImpl) sendResponse(c *gin.Context, response interface{}, err error) {
	if err != nil {
		c.String(500, fmt.Sprintf("Request failed: %s", err.Error()))
		log.Printf("Request failed: %v\n", err)
		return
	}
	if s, ok := response.(string); ok {
		c.String(200, s)
	} else {
		c.JSON(200, response)
	}
}

func closeDbConnection(db *sql.DB) {
	err := db.Close()
	if err != nil {
		// Strange but not critical. So we just log it and go on.
		fmt.Printf("failed to close db connection. Cause %v", err)
	}
}

func (api *restAPIImpl) openDBConnection(c *gin.Context) (*sql.DB, error) {
	config, err := createDbConfig(c)
	if err != nil {
		return nil, err
	}
	config.ValidateServerCertificate(false)
	config.Autocommit(false)
	database, err := sql.Open("exasol", config.String())
	if err != nil {
		return nil, fmt.Errorf("failed to open a database connection. Cause: %w", err)
	}
	return database, nil
}

func createDbConfig(c *gin.Context) (*exasol.DSNConfigBuilder, error) {
	query := c.Request.URL.Query()
	config, err := createDbConfigWithAuthentication(c)
	if err != nil {
		return nil, err
	}

	if host := query.Get("dbHost"); host == "" {
		return nil, fmt.Errorf("missing parameter dbHost")
	} else {
		config.Host(host)
	}

	if portString := query.Get("dbPort"); portString == "" {
		return nil, fmt.Errorf("missing parameter dbPort")
	} else {
		if port, err := strconv.Atoi(portString); err != nil {
			return nil, fmt.Errorf("invalid value %q for parameter dbPort", portString)
		} else {
			config.Port(port)
		}
	}
	return config, nil
}

func createDbConfigWithAuthentication(c *gin.Context) (*exasol.DSNConfigBuilder, error) {
	query := c.Request.URL.Query()
	accessToken := query.Get("dbAccessToken")
	if accessToken != "" {
		return exasol.NewConfigWithAccessToken(accessToken), nil
	}

	refreshToken := query.Get("dbRefreshToken")
	if refreshToken != "" {
		return exasol.NewConfigWithRefreshToken(refreshToken), nil
	}

	user := query.Get("dbUser")
	if user == "" {
		return nil, fmt.Errorf("missing parameter dbUser")
	}

	password := query.Get("dbPassword")
	if password == "" {
		return nil, fmt.Errorf("missing parameter dbPassword")
	}

	return exasol.NewConfig(user, password), nil
}

// General API documentation must be at the end of the file

// @description This is a REST API for managing extensions like virtual schemas in an Exasol database.
// @description
// @description It allows you to install a new extension and create multiple instances for it.
// @description
// @description Authentication is done by passing database connection parameters host, port and credentials via URL parameters. Credentials can be either:
// @description - Username and password (parameters `dbUser` and `dbPassword`)
// @description - Access token (parameter `dbAccessToken`)
// @description - Refresh token (parameter `dbRefreshToken`)
