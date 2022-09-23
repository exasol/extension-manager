package com.exasol.extensionmanager.itest.installer;

import java.nio.file.Path;

import com.exasol.extensionmanager.itest.ExtensionTestConfig;

public interface ExtensionManagerInstaller {

    public static ExtensionManagerInstaller forConfig(final ExtensionTestConfig config) {
        if (config.getLocalExtensionManagerProject().isPresent()) {
            return new InstallerFromLocalFolder(config);
        }
        return new InstallerFromGitHub(config);
    }

    /**
     * Install the extension manager in the given version.
     *
     * @return the path of the executable.
     */
    Path install();
}
