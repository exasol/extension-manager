package com.exasol.extensionmanager.itest;

import java.nio.file.Path;
import java.nio.file.Paths;

import com.exasol.exasoltestsetup.ExasolTestSetup;
import com.exasol.exasoltestsetup.ExasolTestSetupFactory;

class IntegrationTestCommon {
    static final Path TESTING_EXTENSION_SOURCE_DIR = Paths.get("testing-extension").toAbsolutePath();
    static final Path CONFIG_FILE = Paths.get("extension-test.properties").toAbsolutePath();
    static final String EXTENSION_ID = "testing-extension.js";
    static final Path BUILT_EXTENSION_JS = TESTING_EXTENSION_SOURCE_DIR.resolve("dist").resolve(EXTENSION_ID);

    private IntegrationTestCommon() {
        // Not instantiable
    }

    static ExasolTestSetup createExasolTestSetup() {
        if (System.getProperty("com.exasol.dockerdb.image") == null) {
            System.setProperty("com.exasol.dockerdb.image", "8.24.0");
        }
        return new ExasolTestSetupFactory(Path.of("dummy-config")).getTestSetup();
    }
}
