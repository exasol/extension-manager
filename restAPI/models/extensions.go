package models

// ExtensionsResponse contains all available extensions.
type ExtensionsResponse struct {
	Extensions []ExtensionsResponseExtension `json:"extensions"` // All available extensions.
}

// ExtensionsResponseExtension contains information about an available extension that can be installed.
type ExtensionsResponseExtension struct {
	Id                  string   `json:"id"`                  // ID of the extension. Don't store this as it may change in the future.
	Name                string   `json:"name"`                // The name of the extension to be displayed to the user.
	Description         string   `json:"description"`         // The description of the extension to be displayed to the user.
	InstallableVersions []string `json:"installableVersions"` // A list of versions of this extension available for installation.
}

// InstallationsResponse contains all installed extensions.
type InstallationsResponse struct {
	Installations []InstallationsResponseInstallation `json:"installations"`
}

// InstallationsResponseInstallation contains information about installed extensions.
type InstallationsResponseInstallation struct {
	Name               string        `json:"name"`
	Version            string        `json:"version"`
	InstanceParameters []interface{} `json:"instanceParameters"`
}
