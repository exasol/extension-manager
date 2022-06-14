package integrationTesting

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
)

func CreateTestExtensionBuilder() *TestExtensionBuilder {
	builder := TestExtensionBuilder{}
	builder.findInstallationsFunc = "return []"
	return &builder
}

type TestExtensionBuilder struct {
	bucketFsUploads       []BucketFsUploadParams
	findInstallationsFunc string
}

type BucketFsUploadParams struct {
	Name                     string `json:"name"`
	DownloadUrl              string `json:"downloadUrl"`
	LicenseUrl               string `json:"licenseUrl"`
	BucketFsFilename         string `json:"bucketFsFilename"`
	LicenseAgreementRequired bool   `json:"licenseAgreementRequired"`
	FileSize                 int    `json:"fileSize"`
}

func (builder *TestExtensionBuilder) WithFindInstallationsFunc(jsFunctionCode string) *TestExtensionBuilder {
	builder.findInstallationsFunc = jsFunctionCode
	return builder
}

// MockFindInstallationsFunction creates a JS findInstallations function that returns one installation with given JSON array of parameter definitions.
func MockFindInstallationsFunction(extensionName string, version string, parametersJSON string) string {
	template := `return [{
                name: "$NAME$",
                version: "$VERSION$",
                instanceParameters: $PARAMS$
            }]`
	filledTemplate := strings.Replace(template, "$NAME$", extensionName, 1)
	filledTemplate = strings.Replace(filledTemplate, "$VERSION$", version, 1)
	return strings.Replace(filledTemplate, "$PARAMS$", parametersJSON, 1)
}

func (builder *TestExtensionBuilder) WithBucketFsUpload(upload BucketFsUploadParams) *TestExtensionBuilder {
	builder.bucketFsUploads = append(builder.bucketFsUploads, upload)
	return builder
}

//go:embed extensionForTesting/extensionForTestingTemplate.ts
var template string

//go:embed extensionForTesting/package.json
var packageJson []byte

//go:embed extensionForTesting/tsconfig.json
var tscConfig []byte

func (builder TestExtensionBuilder) Build() *BuiltExtension {
	bfsUploadsJson, err := json.Marshal(builder.bucketFsUploads)
	if err != nil {
		panic(err)
	}
	extensionTs := strings.Replace(template, "$UPLOADS$", string(bfsUploadsJson), 1)
	extensionTs = strings.Replace(extensionTs, "$FIND_INSTALLATIONS$", builder.findInstallationsFunc, 1)
	workDir := path.Join(os.TempDir(), "extension-manager-test-extension-build-dir")
	if _, err := os.Stat(workDir); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(workDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	err = ioutil.WriteFile(path.Join(workDir, "package.json"), packageJson, 0600)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path.Join(workDir, "extensionForTesting.ts"), []byte(extensionTs), 0600)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(path.Join(workDir, "tsconfig.json"), tscConfig, 0600)
	if err != nil {
		panic(err)
	}
	return &BuiltExtension{runBuild(workDir)}
}

type BuiltExtension struct {
	content []byte
}

func (extension BuiltExtension) AsString() string {
	return string(extension.content)
}

func (extension BuiltExtension) Bytes() []byte {
	return extension.content
}

func (extension BuiltExtension) WriteToTmpFile() (fileName string) {
	extensionFile, err := ioutil.TempFile(os.TempDir(), "extension-*.js")
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
	err := ioutil.WriteFile(fileName, extension.content, 0600)
	if err != nil {
		panic(err)
	}
}

var buildLock sync.Mutex
var isNpmInstallCalled = false

func runBuild(workDir string) []byte {
	buildLock.Lock()
	runNpmInstall(workDir)
	var stderr bytes.Buffer
	buildCommand := exec.Command("npm", "run", "build")
	buildCommand.Stdout = &stderr
	buildCommand.Stderr = &stderr
	buildCommand.Dir = workDir
	err := buildCommand.Run()
	if err != nil {
		fmt.Println(stderr.String())
		panic(fmt.Sprintf("failed to build extensionForTesting. Cause: %v", err.Error()))
	}
	builtExtension, err := ioutil.ReadFile(path.Join(workDir, "dist.js"))
	if err != nil {
		panic(err)
	}
	buildLock.Unlock()
	return builtExtension
}

func runNpmInstall(workDir string) {
	if !isNpmInstallCalled { // running it once is enough
		var stderr bytes.Buffer
		installCommand := exec.Command("npm", "install")
		installCommand.Dir = workDir
		output, err := installCommand.CombinedOutput()

		if err != nil {
			fmt.Println(stderr.String())
			log.Fatalf("Failed to install node modules (run 'npm install') for extensionForTesting. Cause: %v, Output:\n%s", err.Error(), output)
		}
		isNpmInstallCalled = true
	}
}
