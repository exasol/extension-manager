package com.exasol.extensionmanager;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.hasSize;

import java.nio.file.Path;
import java.nio.file.Paths;
import java.util.List;

import org.junit.jupiter.api.*;

import com.exasol.exasoltestsetup.ExasolTestSetup;
import com.exasol.exasoltestsetup.ExasolTestSetupFactory;
import com.exasol.extensionmanager.client.model.ExtensionsResponseExtension;
import com.exasol.extensionmanager.itest.ExasolVersionCheck;
import com.exasol.extensionmanager.itest.ExtensionManagerSetup;
import com.exasol.extensionmanager.itest.builder.ExtensionBuilder;

/**
 * This integration test illustrates the usage of {@link ExtensionManagerSetup} for testing extensions for Extension
 * Manager (EM). See the
 * <a href="https://github.com/exasol/extension-manager/blob/main/doc/extension_developer_guide.md">Extension Developer
 * Guide</a> for details.
 */
// [doc -> dsn~eitfj-start-extension-manager~1]
// [doc -> dsn~eitfj-access-extension-manager-rest-interface~1]
class ExampleIT {
    /** Relative path to the directory containing the extension definition sources */
    private static final Path EXTENSION_SOURCE_DIR = Paths.get("testing-extension").toAbsolutePath();
    /** File name of the built JavaScript file (= extension ID) */
    private static final String EXTENSION_ID = "testing-extension.js";

    private static ExasolTestSetup exasolTestSetup;
    private static ExtensionManagerSetup setup;

    @BeforeAll
    static void setup() {
        // Overwrite default Exasol version
        System.setProperty("com.exasol.dockerdb.image", "8.32.0");

        exasolTestSetup = new ExasolTestSetupFactory(Path.of("cloud-setup")).getTestSetup();

        // Skip test in case this is an older Exasol version
        ExasolVersionCheck.assumeExasolVersion8(exasolTestSetup);

        // Create EM setup
        setup = ExtensionManagerSetup.create(exasolTestSetup, ExtensionBuilder.createDefaultNpmBuilder(
                EXTENSION_SOURCE_DIR, EXTENSION_SOURCE_DIR.resolve("dist").resolve(EXTENSION_ID)));
    }

    @AfterAll
    static void teardown() throws Exception {
        // Close EM setup and Exasol DB after running all tests
        setup.close();
        exasolTestSetup.close();
    }

    @AfterEach
    void cleanup() {
        // Cleanup resources after each test
        setup.cleanup();
    }

    @Test
    void listExtensions() {
        // Use the EM client to list available extensions and verify the result
        final List<ExtensionsResponseExtension> extensions = setup.client().getExtensions();
        assertThat(extensions, hasSize(1));
    }
}
