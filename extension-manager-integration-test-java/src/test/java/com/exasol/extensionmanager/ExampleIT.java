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
import com.exasol.extensionmanager.itest.ExtensionManagerSetup;
import com.exasol.extensionmanager.itest.builder.ExtensionBuilder;

/**
 * This integration test illustrates the usage of {@link ExtensionManagerSetup}. See the
 * <a href="https://github.com/exasol/extension-manager/blob/main/doc/extension_developer_guide.md">Extension Developer
 * Guide</a> for details.
 */
// [doc -> dsn~eitfj-start-extension-manager~1]
// [doc -> dsn~eitfj-access-extension-manager-rest-interface~1]
class ExampleIT {
    private static final Path EXTENSION_SOURCE_DIR = Paths.get("testing-extension").toAbsolutePath();
    private static final String EXTENSION_ID = "testing-extension.js";
    private static ExasolTestSetup exasolTestSetup;
    private static ExtensionManagerSetup setup;

    @BeforeAll
    static void setup() {
        exasolTestSetup = new ExasolTestSetupFactory(Path.of("cloud-setup")).getTestSetup();
        setup = ExtensionManagerSetup.create(exasolTestSetup, ExtensionBuilder.createDefaultNpmBuilder(
                EXTENSION_SOURCE_DIR, EXTENSION_SOURCE_DIR.resolve("dist").resolve(EXTENSION_ID)));
    }

    @AfterAll
    static void teardown() throws Exception {
        setup.close();
        exasolTestSetup.close();
    }

    @AfterEach
    void cleanup() {
        setup.cleanup();
    }

    @Test
    void listExtensions() {
        final List<ExtensionsResponseExtension> extensions = setup.client().getExtensions();
        assertThat(extensions, hasSize(1));
    }
}
