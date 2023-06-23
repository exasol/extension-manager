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

const DefaultExtensionApiVersion = "0.1.16"

func CreateTestExtensionBuilder(t *testing.T) *TestExtensionBuilder {
	builder := TestExtensionBuilder{testing: t}
	builder.extensionApiVersion = DefaultExtensionApiVersion
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

func (b *TestExtensionBuilder) WithFindInstallationsFunc(tsFunctionCode string) *TestExtensionBuilder {
	b.findInstallationsFunc = tsFunctionCode
	return b
}

func (b *TestExtensionBuilder) WithInstallFunc(tsFunctionCode string) *TestExtensionBuilder {
	b.installFunc = tsFunctionCode
	return b
}

func (b *TestExtensionBuilder) WithUninstallFunc(tsFunctionCode string) *TestExtensionBuilder {
	b.uninstallFunc = tsFunctionCode
	return b
}

func (b *TestExtensionBuilder) WithAddInstanceFunc(tsFunctionCode string) *TestExtensionBuilder {
	b.addInstanceFunc = tsFunctionCode
	return b
}

func (b *TestExtensionBuilder) WithFindInstancesFunc(tsFunctionCode string) *TestExtensionBuilder {
	b.findInstancesFunc = tsFunctionCode
	return b
}

func (b *TestExtensionBuilder) WithDeleteInstanceFunc(tsFunctionCode string) *TestExtensionBuilder {
	b.deleteInstanceFunc = tsFunctionCode
	return b
}

func (b *TestExtensionBuilder) WithGetInstanceParameterDefinitionFunc(tsFunctionCode string) *TestExtensionBuilder {
	b.getInstanceParameterDefinitionsFunc = tsFunctionCode
	return b
}

// MockFindInstallationsFunction creates a JS findInstallations function with extension name and version
func MockFindInstallationsFunction(extensionName string, version string) string {
	template := `return [{name: "$NAME$", version: "$VERSION$"}]`
	filledTemplate := strings.Replace(template, "$NAME$", extensionName, 1)
	filledTemplate = strings.Replace(filledTemplate, "$VERSION$", version, 1)
	return filledTemplate
}

func (b *TestExtensionBuilder) WithBucketFsUpload(upload BucketFsUploadParams) *TestExtensionBuilder {
	b.bucketFsUploads = append(b.bucketFsUploads, upload)
	return b
}

//go:embed extensionForTesting/tsconfig.json
var tscConfig []byte

func (b TestExtensionBuilder) Build() *BuiltExtension {
	workDir := b.createWorkDir()
	err := os.WriteFile(path.Join(workDir, "package.json"), []byte(b.createPackageJsonContent()), 0600)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path.Join(workDir, "extensionForTesting.ts"), []byte(b.createExtensionTsContent()), 0600)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(path.Join(workDir, "tsconfig.json"), tscConfig, 0600)
	if err != nil {
		panic(err)
	}
	content := b.runBuild(workDir)
	return &BuiltExtension{content: content, testing: b.testing}
}

func (b TestExtensionBuilder) createWorkDir() string {
	workDir := path.Join(os.TempDir(), "extension-manager-test-extension", "api-"+b.extensionApiVersion)
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

func (b TestExtensionBuilder) createPackageJsonContent() string {
	return strings.Replace(packageJsonTemplate, "$EXTENSION_API_VERSION$", b.extensionApiVersion, 1)
}

//go:embed extensionForTesting/extensionForTestingTemplate.ts
var extensionTemplate string

func (b TestExtensionBuilder) createExtensionTsContent() string {
	bfsUploadsJson, err := json.Marshal(b.bucketFsUploads)
	if err != nil {
		panic(err)
	}
	content := strings.Replace(extensionTemplate, "$UPLOADS$", string(bfsUploadsJson), 1)
	content = strings.Replace(content, "$FIND_INSTALLATIONS$", b.findInstallationsFunc, 1)
	content = strings.Replace(content, "$INSTALL_EXTENSION$", b.installFunc, 1)
	content = strings.Replace(content, "$UNINSTALL_EXTENSION$", b.uninstallFunc, 1)
	content = strings.Replace(content, "$ADD_INSTANCE$", b.addInstanceFunc, 1)
	content = strings.Replace(content, "$FIND_INSTANCES$", b.findInstancesFunc, 1)
	content = strings.Replace(content, "$DELETE_INSTANCE$", b.deleteInstanceFunc, 1)
	content = strings.Replace(content, "$GET_INSTANCE_PARAMETER_DEFINITIONS$", b.getInstanceParameterDefinitionsFunc, 1)
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

func (e BuiltExtension) WriteToTmpFile() (fileName string) {
	extensionFile, err := os.CreateTemp(e.testing.TempDir(), "extension-*.js")
	if err != nil {
		panic(err)
	}
	defer func() {
		err := extensionFile.Close()
		if err != nil {
			panic(err)
		}
	}()
	_, err = extensionFile.Write(e.content)
	if err != nil {
		panic(err)
	}
	return extensionFile.Name()
}

func (e BuiltExtension) WriteToFile(fileName string) {
	err := os.WriteFile(fileName, e.content, 0600)
	if err != nil {
		panic(err)
	}
	cleanupFile(e.testing, fileName)
}

func (e BuiltExtension) Publish(server *MockRegistryServer, id string) {
	path := "/" + id + ".js"
	extensionUrl := server.BaseUrl() + path
	server.SetRegistryContent(fmt.Sprintf(`{"extensions":[{"id": "%s", "url": "%s"}]}`, id, extensionUrl))
	server.SetPathContent(path, e.AsString())
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

func (b TestExtensionBuilder) runBuild(workDir string) []byte {
	buildLock.Lock()
	b.runNpmInstall(workDir)
	b.runNpmBuild(workDir)
	path := path.Join(workDir, "dist.js")
	cleanupFile(b.testing, path)
	builtExtension, err := os.ReadFile(path)
	if err != nil {
		b.testing.Fatalf("failed to read %s: %v", path, err)
	}
	buildLock.Unlock()
	return builtExtension
}

var isNpmInstallCalledForVersion = make(map[string]bool)

func (b TestExtensionBuilder) runNpmInstall(workDir string) {
	if isNpmInstallCalledForVersion[b.extensionApiVersion] {
		// running "npm install" once for each version is enough
		return
	}
	b.testing.Logf("Running npm install in %s", workDir)
	installCommand := exec.Command("npm", "install")
	installCommand.Dir = workDir
	output, err := installCommand.CombinedOutput()
	if err != nil {
		log.Fatalf("Failed to run 'npm install' in dir %v. Cause: %v, Output:\n%s", workDir, err, output)
	} else {
		isNpmInstallCalledForVersion[b.extensionApiVersion] = true
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
