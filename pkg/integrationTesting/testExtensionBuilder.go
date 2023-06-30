package integrationTesting

import (
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"testing"

	log "github.com/sirupsen/logrus"
)

const defaultExtensionApiVersion = "0.2.0"

func CreateTestExtensionBuilder(t *testing.T) *TestExtensionBuilder {
	builder := TestExtensionBuilder{testing: t}
	builder.extensionApiVersion = defaultExtensionApiVersion
	builder.findInstallationsFunc = "return []"
	builder.installFunc = "context.sqlClient.execute('select 1')"
	builder.uninstallFunc = ""
	builder.addInstanceFunc = "return undefined"
	builder.findInstancesFunc = "return []"
	builder.deleteInstanceFunc = "context.sqlClient.execute(`drop instance ${instanceId}`)"
	builder.getInstanceParameterDefinitionsFunc = "return []"
	return &builder
}

type TestExtensionBuilder struct {
	testing                             *testing.T
	extensionApiVersion                 string
	bucketFsUploads                     []BucketFsUploadParams
	findInstallationsFunc               string
	installFunc                         string
	uninstallFunc                       string
	addInstanceFunc                     string
	findInstancesFunc                   string
	deleteInstanceFunc                  string
	getInstanceParameterDefinitionsFunc string
}

type BucketFsUploadParams struct {
	Name                     string `json:"name"`
	DownloadUrl              string `json:"downloadUrl"`
	LicenseUrl               string `json:"licenseUrl"`
	BucketFsFilename         string `json:"bucketFsFilename"`
	LicenseAgreementRequired bool   `json:"licenseAgreementRequired"`
	FileSize                 int    `json:"fileSize"`
}

func (builder *TestExtensionBuilder) WithFindInstallationsFunc(tsFunctionCode string) *TestExtensionBuilder {
	builder.findInstallationsFunc = tsFunctionCode
	return builder
}

func (builder *TestExtensionBuilder) WithInstallFunc(tsFunctionCode string) *TestExtensionBuilder {
	builder.installFunc = tsFunctionCode
	return builder
}

func (builder *TestExtensionBuilder) WithUninstallFunc(tsFunctionCode string) *TestExtensionBuilder {
	builder.uninstallFunc = tsFunctionCode
	return builder
}

func (builder *TestExtensionBuilder) WithAddInstanceFunc(tsFunctionCode string) *TestExtensionBuilder {
	builder.addInstanceFunc = tsFunctionCode
	return builder
}

func (builder *TestExtensionBuilder) WithFindInstancesFunc(tsFunctionCode string) *TestExtensionBuilder {
	builder.findInstancesFunc = tsFunctionCode
	return builder
}

func (builder *TestExtensionBuilder) WithDeleteInstanceFunc(tsFunctionCode string) *TestExtensionBuilder {
	builder.deleteInstanceFunc = tsFunctionCode
	return builder
}

func (builder *TestExtensionBuilder) WithGetInstanceParameterDefinitionFunc(tsFunctionCode string) *TestExtensionBuilder {
	builder.getInstanceParameterDefinitionsFunc = tsFunctionCode
	return builder
}

// MockFindInstallationsFunction creates a JS findInstallations function with extension name and version
func MockFindInstallationsFunction(extensionName string, version string) string {
	template := `return [{name: "$NAME$", version: "$VERSION$"}]`
	filledTemplate := strings.Replace(template, "$NAME$", extensionName, 1)
	filledTemplate = strings.Replace(filledTemplate, "$VERSION$", version, 1)
	return filledTemplate
}

func (builder *TestExtensionBuilder) WithBucketFsUpload(upload BucketFsUploadParams) *TestExtensionBuilder {
	builder.bucketFsUploads = append(builder.bucketFsUploads, upload)
	return builder
}

//go:embed extensionForTesting/tsconfig.json
var tscConfig []byte

func (builder TestExtensionBuilder) Build() *BuiltExtension {
	workDir := builder.createWorkDir()
	err := os.WriteFile(path.Join(workDir, "package.json"), []byte(builder.createPackageJsonContent()), 0600)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path.Join(workDir, "extensionForTesting.ts"), []byte(builder.createExtensionTsContent()), 0600)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path.Join(workDir, "tsconfig.json"), tscConfig, 0600)
	if err != nil {
		panic(err)
	}
	content := builder.runBuild(workDir)
	return &BuiltExtension{content: content, testing: builder.testing}
}

func (builder TestExtensionBuilder) createWorkDir() string {
	workDir := path.Join(os.TempDir(), "extension-manager-test-extension", "api-"+builder.extensionApiVersion)
	if _, err := os.Stat(workDir); errors.Is(err, os.ErrNotExist) {
		err := os.MkdirAll(workDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	return workDir
}

//go:embed extensionForTesting/package.json
var packageJsonTemplate string

func (builder TestExtensionBuilder) createPackageJsonContent() string {
	return strings.Replace(packageJsonTemplate, "$EXTENSION_API_VERSION$", builder.extensionApiVersion, 1)
}

//go:embed extensionForTesting/extensionForTestingTemplate.ts
var extensionTemplate string

func (builder TestExtensionBuilder) createExtensionTsContent() string {
	bfsUploadsJson, err := json.Marshal(builder.bucketFsUploads)
	if err != nil {
		panic(err)
	}
	content := strings.Replace(extensionTemplate, "$UPLOADS$", string(bfsUploadsJson), 1)
	content = strings.Replace(content, "$FIND_INSTALLATIONS$", builder.findInstallationsFunc, 1)
	content = strings.Replace(content, "$INSTALL_EXTENSION$", builder.installFunc, 1)
	content = strings.Replace(content, "$UNINSTALL_EXTENSION$", builder.uninstallFunc, 1)
	content = strings.Replace(content, "$ADD_INSTANCE$", builder.addInstanceFunc, 1)
	content = strings.Replace(content, "$FIND_INSTANCES$", builder.findInstancesFunc, 1)
	content = strings.Replace(content, "$DELETE_INSTANCE$", builder.deleteInstanceFunc, 1)
	content = strings.Replace(content, "$GET_INSTANCE_PARAMETER_DEFINITIONS$", builder.getInstanceParameterDefinitionsFunc, 1)
	return content
}

type BuiltExtension struct {
	testing *testing.T
	content []byte
}

func (extension BuiltExtension) AsString() string {
	return string(extension.content)
}

func (extension BuiltExtension) Bytes() []byte {
	return extension.content
}

func (extension BuiltExtension) WriteToTmpFile() (fileName string) {
	extensionFile, err := os.CreateTemp(extension.testing.TempDir(), "extension-*.js")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := extensionFile.Close()
		if err != nil {
			panic(err)
		}
	}()
	_, err = extensionFile.Write(extension.content)
	if err != nil {
		panic(err)
	}
	return extensionFile.Name()
}

func (extension BuiltExtension) WriteToFile(fileName string) {
	err := os.WriteFile(fileName, extension.content, 0600)
	if err != nil {
		panic(err)
	}
	cleanupFile(extension.testing, fileName)
}

func (extension BuiltExtension) Publish(server *MockRegistryServer, id string) {
	path := "/" + id + ".js"
	extensionUrl := server.BaseUrl() + path
	server.SetRegistryContent(fmt.Sprintf(`{"extensions":[{"id": "%s", "url": "%s"}]}`, id, extensionUrl))
	server.SetPathContent(path, extension.AsString())
}

func cleanupFile(t *testing.T, fileName string) {
	t.Cleanup(func() {
		if _, err := os.Stat(fileName); errors.Is(err, os.ErrNotExist) {
			return
		}
		err := os.Remove(fileName)
		if err != nil {
			t.Errorf("failed to delete file: %v", err)
		}
	})
}

var buildLock sync.Mutex

func (builder TestExtensionBuilder) runBuild(workDir string) []byte {
	buildLock.Lock()
	builder.runNpmInstall(workDir)
	builder.runNpmBuild(workDir)
	path := path.Join(workDir, "dist.js")
	cleanupFile(builder.testing, path)
	builtExtension, err := os.ReadFile(path)
	if err != nil {
		builder.testing.Fatalf("failed to read %s: %v", path, err)
	}
	buildLock.Unlock()
	return builtExtension
}

var isNpmInstallCalledForVersion = make(map[string]bool)

func (builder TestExtensionBuilder) runNpmInstall(workDir string) {
	if isNpmInstallCalledForVersion[builder.extensionApiVersion] {
		// running "npm install" once for each version is enough
		return
	}
	builder.testing.Logf("Running npm install in %s", workDir)
	installCommand := exec.Command("npm", "install")
	installCommand.Dir = workDir
	output, err := installCommand.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to run 'npm install' in dir %v. Cause: %v, Output:\n%s", workDir, err, output)
	} else {
		isNpmInstallCalledForVersion[builder.extensionApiVersion] = true
	}
}

func (TestExtensionBuilder) runNpmBuild(workDir string) {
	buildCommand := exec.Command("npm", "run", "build")
	buildCommand.Dir = workDir
	output, err := buildCommand.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to run 'npm run build' in dir %v. Cause: %v, Output:\n%s", workDir, err, output)
	}
}
