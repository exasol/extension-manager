package extensionApi

import (
	"fmt"
	"github.com/robertkrimen/otto"
	"os"
)

func GetExtensionFromFile(fileName string, sqlClient SimpleSqlClient) (Extension, error) {
	vm := otto.New()
	sqlClientJs, err := addSqlClient(vm, sqlClient)
	if err != nil {
		return nil, err
	}
	extensionJs, err := loadExtension(vm, fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to load extension %v. Cause: %v", fileName, err.Error())
	}
	extensionName, err := readRequiredStringProperty(extensionJs, "name")
	if err != nil {
		return nil, err
	}
	extension := ExtensionImpl{*extensionName, extensionJs, sqlClientJs}
	return extension, nil
}

func readRequiredStringProperty(extensionJs *otto.Object, propertyName string) (*string, error) {
	result, err := extensionJs.Get(propertyName)
	if err != nil {
		return nil, fmt.Errorf("failed to read required extension property %v. Cause: %v", propertyName, err.Error())
	}
	stringResult, err := result.ToString()
	if err != nil {
		return nil, fmt.Errorf("invalid value for extension property %v. The value must be a string. Cause: %v", propertyName, err.Error())
	}
	return &stringResult, nil
}

func loadExtension(vm *otto.Otto, fileName string) (*otto.Object, error) {
	const extensionVariable = "installedExtension"
	err := vm.Set(extensionVariable, otto.NullValue())
	if err != nil {
		return nil, fmt.Errorf("failed to set installedExtension = null. Cause: %v", err.Error())
	}
	file, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open extension file %v. Cause: %v", fileName, err.Error())
	}
	_, err = vm.Run(file)
	if err != nil {
		return nil, fmt.Errorf("failed to run extension file %v. Cause %v", fileName, err.Error())
	}
	installedExtension, err := vm.Get(extensionVariable)
	if err != nil {
		return nil, fmt.Errorf("failed to read installedExtension variable. Cause: %v", err.Error())
	}
	installedExtensionObject := installedExtension.Object()
	if installedExtensionObject == nil {
		return nil, fmt.Errorf("invalid installedExtension. The provided JS file did not set the installedExtension variable or it is not an object. Make sure that the extension fill calls installExtension")
	}
	err = assertApiVersion(installedExtensionObject)
	if err != nil {
		return nil, err
	}
	extension, err := installedExtensionObject.Get("extension")
	if err != nil {
		return nil, fmt.Errorf("failed to read the value of installedExtension.extension. Cause: %v", err.Error())
	}
	extensionObject := extension.Object()
	if extensionObject == nil {
		return nil, fmt.Errorf("invalid installedExtension. The provided JS file did not set the installedExtension.extension variable or it is not an object. Make sure that the extension fill calls installExtension")
	}
	return extensionObject, nil
}

func assertApiVersion(installedExtensionObject *otto.Object) error {
	apiVersion, err := getApiVersion(installedExtensionObject)
	if err != nil {
		return err
	}
	if *apiVersion != "0.1.0" {
		return fmt.Errorf("incompatible extension API version %v. Please update the extension to use a supported version of the extension API", *apiVersion)
	}
	return nil
}

func getApiVersion(installedExtensionObject *otto.Object) (*string, error) {
	apiVersion, err := installedExtensionObject.Get("apiVersion")
	if err != nil {
		return nil, fmt.Errorf("invalid installedExtension. Could not read extension.apiVersion. Cause: %v", err.Error())
	}
	if !apiVersion.IsString() {
		return nil, fmt.Errorf("invalid installedExtension.apiVersion. The field should be a string")
	}
	apiVersionString := apiVersion.String()
	return &apiVersionString, nil
}

type SimpleSqlClient interface {
	RunSqlQuery(query string)
}

func addSqlClient(vm *otto.Otto, sqlClient SimpleSqlClient) (*otto.Object, error) {
	sqlClientJs, err := vm.Object(`sqlClient={}`)
	if err != nil {
		return nil, fmt.Errorf("failed to set SQL client. Cause %v\n", err.Error())
	}
	err = sqlClientJs.Set("runQuery", func(call otto.FunctionCall) otto.Value {
		sqlClient.RunSqlQuery(call.Argument(0).String())
		return otto.Value{}
	})
	if err != nil {
		return nil, fmt.Errorf("failed to install sqlClient.runQuery function. Cause: %v", err.Error())
	}
	return sqlClientJs, nil
}

type LoggingSimpleSqlClient struct {
}

func (client LoggingSimpleSqlClient) RunSqlQuery(query string) {
	fmt.Printf("sql: %v\n", query)
}

type Extension interface {
	GetName() string
	Install() error
}

type ExtensionImpl struct {
	name        string
	extensionJs *otto.Object
	sqlClientJs *otto.Object
}

func (extension ExtensionImpl) GetName() string {
	return extension.name
}

func (extension ExtensionImpl) Install() error {
	_, err := extension.extensionJs.Call("install", extension.sqlClientJs)
	if err != nil {
		return fmt.Errorf("failed to run install function. Cause: %v", err)
	}
	return nil
}
