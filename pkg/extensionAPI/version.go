package extensionAPI

import (
	"fmt"

	"golang.org/x/mod/semver"
)

const supportedApiVersion = "0.2.0"

/* [impl -> dsn~extension-compatibility~1] */
func validateExtensionIsCompatibleWithApiVersion(extensionId, currentExtensionApiVersion string) error {
	prefixedVersion := "v" + currentExtensionApiVersion
	if !semver.IsValid(prefixedVersion) {
		return fmt.Errorf("extension %q uses invalid API version number %q", extensionId, currentExtensionApiVersion)
	}
	major := semver.Major(prefixedVersion)
	if major != getSupportedMajorVersion() {
		return fmt.Errorf("extension %q uses incompatible API version %q. Please update the extension to use supported version %q", extensionId, currentExtensionApiVersion, supportedApiVersion)
	}
	return nil
}

func getSupportedMajorVersion() string {
	prefixedVersion := "v" + supportedApiVersion
	if !semver.IsValid(prefixedVersion) {
		panic(fmt.Errorf("version %q has an invalid format", supportedApiVersion))
	}
	return semver.Major(prefixedVersion)
}
