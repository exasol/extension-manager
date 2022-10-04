package com.exasol.extensionmanager.itest.installer;

import java.nio.file.Path;

import com.exasol.extensionmanager.itest.ExtensionTestConfig;

/**
 * This class installs the extension manager depending on the given configuration either from GitHub or from a local
 * directory.
 */
public interface ExtensionManagerInstaller {

    /**
     * Create a new installer depending on the configuration. If a
     * {@link ExtensionTestConfig#getLocalExtensionManagerProject()} is given, then the installer uses the configured
     * local folder. Else the installer will install the extension manager from GitHub.
     * 
     * @param config test configuration
     * @return a new installer
     */
    public static ExtensionManagerInstaller forConfig(final ExtensionTestConfig config) {
        if (config.getLocalExtensionManagerProject().isPresent()) {
            return new InstallerFromLocalFolder(config);
        }
        return new InstallerFromGitHub(config);
    }

    /**
     * Install the extension manager.
     *
     * @return path of the executable
     */
    Path install();
}
