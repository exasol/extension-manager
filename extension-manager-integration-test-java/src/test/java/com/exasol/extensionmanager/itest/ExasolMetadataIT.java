package com.exasol.extensionmanager.itest;

import static com.exasol.extensionmanager.itest.IntegrationTestCommon.TESTING_EXTENSION_SOURCE_DIR;
import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;

import org.junit.jupiter.api.*;

import com.exasol.exasoltestsetup.ExasolTestSetup;
import com.exasol.extensionmanager.itest.builder.ExtensionBuilder;

class ExasolMetadataIT {

    private static ExasolTestSetup exasolTestSetup;
    private static ExtensionManagerSetup extensionManager;
    private static ExasolMetadata metadata;

    @BeforeAll
    static void setupExasol() {
        exasolTestSetup = IntegrationTestCommon.createExasolTestSetup();
        extensionManager = ExtensionManagerSetup.create(exasolTestSetup,
                ExtensionBuilder.createDefaultNpmBuilder(TESTING_EXTENSION_SOURCE_DIR,
                        TESTING_EXTENSION_SOURCE_DIR.resolve("dist").resolve(IntegrationTestCommon.EXTENSION_ID)));
        metadata = extensionManager.exasolMetadata();
    }

    @AfterAll
    static void tearDownExasol() throws Exception {
        extensionManager.close();
        exasolTestSetup.close();
    }

    @Test
    void assertNoConnections() {
        assertDoesNotThrow(metadata::assertNoConnections);
    }

    @Test
    void assertNoVirtualSchema() {
        assertDoesNotThrow(metadata::assertNoVirtualSchema);
    }

    @Test
    void assertNoScripts() {
        assertDoesNotThrow(metadata::assertNoScripts);
    }
}
