package com.exasol.extensionmanager.itest;

import java.nio.file.Path;
import java.nio.file.Paths;

class IntegrationTestCommon {
    static final Path TESTING_EXTENSION_SOURCE_DIR = Paths.get("testing-extension").toAbsolutePath();
    static final Path CONFIG_FILE = Paths.get("extension-test.properties").toAbsolutePath();

    private IntegrationTestCommon() {
        // Not instantiable
    }
}
