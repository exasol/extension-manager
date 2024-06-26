package extensionController

import (
	"fmt"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/exasol/extension-manager/pkg/apiErrors"
	"github.com/exasol/extension-manager/pkg/extensionAPI"
	"github.com/exasol/extension-manager/pkg/extensionAPI/context"
	"github.com/exasol/extension-manager/pkg/extensionAPI/exaMetadata"
	"github.com/exasol/extension-manager/pkg/extensionController/bfs"
	"github.com/exasol/extension-manager/pkg/extensionController/registry"
	"github.com/exasol/extension-manager/pkg/extensionController/transaction"

	"github.com/exasol/extension-manager/pkg/parameterValidator"
)

// controller is the core part of the extension-manager that provides the extension handling functionality.
type controller interface {
	// GetAllExtensions reports all extension definitions.
	GetAllExtensions(bfsFiles []bfs.BfsFile) ([]*Extension, error)

	// GetAllInstallations searches for installations of any extensions.
	GetAllInstallations(txCtx *transaction.TransactionContext) ([]*extensionAPI.JsExtInstallation, error)

	// GetParameterDefinitions returns the parameter definitions required for installing a given extension version.
	GetParameterDefinitions(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string) ([]parameterValidator.ParameterDefinition, error)

	// InstallExtension installs an extension.
	// extensionId is the ID of the extension to install
	// extensionVersion is the version of the extension to install
	InstallExtension(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string) error

	// UninstallExtension removes an extension.
	// extensionId is the ID of the extension to uninstall
	// extensionVersion is the version of the extension to uninstall
	UninstallExtension(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string) error

	// UpgradeExtension upgrades an installed extension to the latest version.
	// extensionId is the ID of the extension to uninstall
	UpgradeExtension(txCtx *transaction.TransactionContext, extensionId string) (*extensionAPI.JsUpgradeResult, error)

	// CreateInstance creates a new instance of an extension, e.g. a virtual schema and returns it's name.
	CreateInstance(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string, parameterValues []ParameterValue) (*extensionAPI.JsExtInstance, error)

	// FindInstances returns a list of all instances for the given version.
	FindInstances(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string) ([]*extensionAPI.JsExtInstance, error)

	// DeleteInstance deletes instance with the given ID.
	DeleteInstance(txCtx *transaction.TransactionContext, extensionId, extensionVersion, instanceId string) error
}

type controllerImpl struct {
	registry       registry.Registry
	config         ExtensionManagerConfig
	metaDataReader exaMetadata.ExaMetadataReader
}

func createImpl(config ExtensionManagerConfig) controller {
	return &controllerImpl{
		registry:       registry.NewRegistry(config.ExtensionRegistryURL),
		metaDataReader: exaMetadata.CreateExaMetaDataReader(),
		config:         config,
	}
}

/* [impl -> dsn~list-extensions~1]. */
func (c *controllerImpl) GetAllExtensions(bfsFiles []bfs.BfsFile) ([]*Extension, error) {
	jsExtensions, err := c.getAllExtensions()
	if err != nil {
		return nil, err
	}
	var extensions []*Extension
	for _, jsExtension := range jsExtensions {
		if c.requiredFilesAvailable(jsExtension, bfsFiles) {
			extensions = append(extensions, convertExtension(jsExtension))
		}
	}
	log.Infof("Found %d of %d extensions with required files (%d files available in total)", len(extensions), len(jsExtensions), len(bfsFiles))
	return extensions, nil
}

func convertExtension(jsExtension *extensionAPI.JsExtension) *Extension {
	return &Extension{
		Id:                  jsExtension.Id,
		Name:                jsExtension.Name,
		Category:            jsExtension.Category,
		Description:         jsExtension.Description,
		InstallableVersions: jsExtension.InstallableVersions}
}

func (c *controllerImpl) requiredFilesAvailable(extension *extensionAPI.JsExtension, bfsFiles []bfs.BfsFile) bool {
	for _, requiredFile := range extension.BucketFsUploads {
		if !existsFileInBfs(bfsFiles, requiredFile) {
			log.Debugf("Ignoring extension %q since the required file %q does not exist or has a wrong file size.\n", extension.Name, requiredFile.BucketFsFilename)
			return false
		}
	}
	return true
}

func existsFileInBfs(bfsFiles []bfs.BfsFile, requiredFile extensionAPI.BucketFsUpload) bool {
	for _, existingFile := range bfsFiles {
		if fileMatches(requiredFile, existingFile) {
			return true
		}
	}
	log.Tracef("Required file %q of size %db not found", requiredFile.Name, requiredFile.FileSize)
	return false
}

func fileMatches(requiredFile extensionAPI.BucketFsUpload, existingFile bfs.BfsFile) bool {
	if requiredFile.BucketFsFilename != existingFile.Name {
		return false
	}
	if requiredFile.FileSize < 0 {
		log.Tracef("Found required file %q of size %db ignoring file size", existingFile.Name, existingFile.Size)
		return true
	}
	if requiredFile.FileSize == existingFile.Size {
		log.Tracef("Found required file %q of size %db", existingFile.Name, existingFile.Size)
		return true
	}
	log.Debugf("File %q exists but has wrong size %d, expected %d bytes", existingFile.Name, existingFile.Size, requiredFile.FileSize)
	return false
}

func (c *controllerImpl) getAllExtensions() ([]*extensionAPI.JsExtension, error) {
	t0 := time.Now()
	extensionIds, err := c.registry.FindExtensions()
	if err != nil {
		return nil, err
	}
	extensions := make([]*extensionAPI.JsExtension, 0, len(extensionIds))
	for _, id := range extensionIds {
		extension, err := c.loadExtensionById(id)
		if err != nil {
			return nil, fmt.Errorf("failed to load extension %q: %w", id, err)
		}
		extensions = append(extensions, extension)
	}
	log.Debugf("Loaded %d extensions JS files in %dms", len(extensions), time.Since(t0).Milliseconds())
	return extensions, nil
}

func (c *controllerImpl) loadExtensionById(id string) (*extensionAPI.JsExtension, error) {
	content, err := c.registry.ReadExtension(id)
	if err != nil {
		return nil, err
	}
	extension, err := extensionAPI.LoadExtension(id, content)
	if err != nil {
		return nil, err
	}
	return extension, nil
}

func (c *controllerImpl) GetAllInstallations(txCtx *transaction.TransactionContext) ([]*extensionAPI.JsExtInstallation, error) {
	metadata, err := c.metaDataReader.ReadMetadataTables(txCtx.GetTransaction(), c.config.ExtensionSchema)
	if err != nil {
		return nil, fmt.Errorf("failed to read metadata tables. Cause: %w", err)
	}
	extensions, err := c.getAllExtensions()
	if err != nil {
		return nil, err
	}
	extensionContext := c.createExtensionContext(txCtx)
	var allInstallations []*extensionAPI.JsExtInstallation
	for _, extension := range extensions {
		installations, err := extension.FindInstallations(extensionContext, metadata)
		if err != nil {
			return nil, apiErrors.NewAPIErrorWithCause(fmt.Sprintf("failed to find installations for extension %q", extension.Name), err)
		}
		addExtensionId(extension.Id, installations)
		c.logInstallations(extension, installations)
		allInstallations = append(allInstallations, installations...)
	}
	return allInstallations, nil
}

func addExtensionId(extensionID string, installations []*extensionAPI.JsExtInstallation) {
	for _, i := range installations {
		i.ID = extensionID
	}
}

func (*controllerImpl) logInstallations(extension *extensionAPI.JsExtension, installations []*extensionAPI.JsExtInstallation) {
	if len(installations) == 0 {
		log.Debugf("Found no installations for extension %q", extension.Id)
		return
	}
	log.Debugf("Found %d installations for extension %q", len(installations), extension.Id)
	if log.IsLevelEnabled(log.DebugLevel) {
		for _, installation := range installations {
			log.Debugf("- %q: version %s", installation.Name, installation.Version)
		}
	}
}

func (c *controllerImpl) GetParameterDefinitions(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string) ([]parameterValidator.ParameterDefinition, error) {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return nil, extensionLoadingFailed(extensionId, err)
	}
	rawDefinitions, err := extension.GetParameterDefinitions(c.createExtensionContext(txCtx), extensionVersion)
	if err != nil {
		return nil, err
	}
	definitions, err := parameterValidator.ConvertDefinitions(rawDefinitions)
	if err != nil {
		return nil, err
	}
	return definitions, nil
}

func (c *controllerImpl) InstallExtension(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string) error {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return extensionLoadingFailed(extensionId, err)
	}
	err = c.ensureSchemaExists(txCtx)
	if err != nil {
		return err
	}
	return extension.Install(c.createExtensionContext(txCtx), extensionVersion)
}

func (c *controllerImpl) UninstallExtension(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string) error {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return extensionLoadingFailed(extensionId, err)
	}
	extensionCtx := c.createExtensionContext(txCtx)
	err = c.verifyNoInstances(extension, extensionCtx, extensionVersion)
	if err != nil {
		return fmt.Errorf("cannot uninstall extension because instances remain: %w", err)
	}
	return extension.Uninstall(extensionCtx, extensionVersion)
}

func (*controllerImpl) verifyNoInstances(extension *extensionAPI.JsExtension, extensionCtx *context.ExtensionContext, extensionVersion string) error {
	if !extension.SupportsListInstances(extensionCtx, extensionVersion) {
		return nil
	}
	instances, err := extension.ListInstances(extensionCtx, extensionVersion)
	if err != nil {
		return fmt.Errorf("failed to check existing instances: %w", err)
	}
	if len(instances) > 0 {
		instanceNames := concatInstanceNames(instances)
		return apiErrors.NewBadRequestErrorF("cannot uninstall extension because %d instance(s) still exist: %s", len(instances), instanceNames)
	}
	return nil
}

func concatInstanceNames(instances []*extensionAPI.JsExtInstance) string {
	instanceNames := ""
	for _, inst := range instances {
		if instanceNames != "" {
			instanceNames += ", "
		}
		instanceNames += inst.Name
	}
	return instanceNames
}

/* [impl -> dsn~upgrade-extension~1]. */
func (c *controllerImpl) UpgradeExtension(txCtx *transaction.TransactionContext, extensionId string) (*extensionAPI.JsUpgradeResult, error) {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return nil, extensionLoadingFailed(extensionId, err)
	}
	return extension.Upgrade(c.createExtensionContext(txCtx))
}

func (c *controllerImpl) CreateInstance(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string, parameterValues []ParameterValue) (*extensionAPI.JsExtInstance, error) {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return nil, extensionLoadingFailed(extensionId, err)
	}
	err = c.ensureSchemaExists(txCtx)
	if err != nil {
		return nil, err
	}

	params, err := c.convertAndValidate(txCtx, extensionId, extensionVersion, parameterValues)
	if err != nil {
		return nil, err
	}

	extensionContext := c.createExtensionContext(txCtx)
	instance, err := extension.AddInstance(extensionContext, extensionVersion, &params)
	if err != nil {
		return nil, err
	}
	if instance == nil {
		return nil, fmt.Errorf("extension %q did not return an instance", extensionId)
	}
	return instance, nil
}

func (c *controllerImpl) convertAndValidate(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string, parameterValues []ParameterValue) (extensionAPI.ParameterValues, error) {
	paramDefinitions, err := c.GetParameterDefinitions(txCtx, extensionId, extensionVersion)
	if err != nil {
		return extensionAPI.ParameterValues{}, fmt.Errorf("failed to get parameter definitions: %w", err)
	}
	params := convertParameters(parameterValues)
	err = validateParameters(paramDefinitions, params)
	if err != nil {
		return extensionAPI.ParameterValues{}, err
	}
	return params, nil
}

func convertParameters(parameterValues []ParameterValue) extensionAPI.ParameterValues {
	values := []extensionAPI.ParameterValue{}
	for _, p := range parameterValues {
		values = append(values, extensionAPI.ParameterValue{Name: p.Name, Value: p.Value})
	}
	return extensionAPI.ParameterValues{Values: values}
}

func (c *controllerImpl) DeleteInstance(txCtx *transaction.TransactionContext, extensionId, extensionVersion, instanceId string) error {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return extensionLoadingFailed(extensionId, err)
	}
	return extension.DeleteInstance(c.createExtensionContext(txCtx), extensionVersion, instanceId)
}

func (c *controllerImpl) FindInstances(txCtx *transaction.TransactionContext, extensionId string, extensionVersion string) ([]*extensionAPI.JsExtInstance, error) {
	extension, err := c.loadExtensionById(extensionId)
	if err != nil {
		return nil, extensionLoadingFailed(extensionId, err)
	}
	return extension.ListInstances(c.createExtensionContext(txCtx), extensionVersion)
}

func (c *controllerImpl) createExtensionContext(txCtx *transaction.TransactionContext) *context.ExtensionContext {
	return context.CreateContext(txCtx, c.config.ExtensionSchema)
}

func (c *controllerImpl) ensureSchemaExists(txCtx *transaction.TransactionContext) error {
	_, err := txCtx.GetTransaction().Exec(fmt.Sprintf(`CREATE SCHEMA IF NOT EXISTS "%s"`, c.config.ExtensionSchema))
	if err != nil {
		return fmt.Errorf("failed to create schema: %w", err)
	}
	return nil
}

/* [impl -> dsn~parameter-types~1]. */
func validateParameters(parameterDefinitions []parameterValidator.ParameterDefinition, params extensionAPI.ParameterValues) error {
	validator, err := parameterValidator.New()
	if err != nil {
		return fmt.Errorf("failed to create parameter validator: %w", err)
	}
	result, err := validator.ValidateParameters(parameterDefinitions, params)
	if err != nil {
		return fmt.Errorf("failed to validate parameters: %w", err)
	}
	message := ""
	for _, r := range result {
		message += r.Message + ", "
	}
	message = strings.TrimSuffix(message, ", ")
	if message != "" {
		return apiErrors.NewBadRequestErrorF("invalid parameters: %s", message)
	}
	return nil
}

func extensionLoadingFailed(extensionId string, err error) error {
	return fmt.Errorf("failed to load extension %q: %w", extensionId, err)
}
