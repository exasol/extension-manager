package com.exasol.extensionmanager.itest;

import static com.exasol.extensionmanager.itest.IntegrationTestCommon.TESTING_EXTENSION_SOURCE_DIR;
import static com.exasol.matcher.ResultSetStructureMatcher.table;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.*;

import java.io.IOException;
import java.io.UncheckedIOException;
import java.net.URI;
import java.nio.file.Files;
import java.nio.file.Path;
import java.sql.PreparedStatement;
import java.sql.SQLException;

import org.junit.jupiter.api.*;

import com.exasol.dbbuilder.dialects.exasol.ExasolSchema;
import com.exasol.exasoltestsetup.ExasolTestSetup;
import com.exasol.extensionmanager.itest.builder.ExtensionBuilder;

class ExtensionManagerSetupIT {

    private static ExasolTestSetup exasolTestSetup;

    @BeforeAll
    static void setupExasol() {
        exasolTestSetup = IntegrationTestCommon.createExasolTestSetup();
    }

    @AfterAll
    static void teardownExasol() throws Exception {
        exasolTestSetup.close();
    }

    private ExtensionManagerSetup extensionManager;

    @BeforeEach
    void setup() throws SQLException {
        extensionManager = ExtensionManagerSetup.create(exasolTestSetup, ExtensionBuilder.createDefaultNpmBuilder(
                TESTING_EXTENSION_SOURCE_DIR, TESTING_EXTENSION_SOURCE_DIR.resolve("dist/testing-extension.js")));
    }

    @AfterEach
    void teardown() {
        extensionManager.close();
    }

    @Test
    void cleanupConnections() {
        extensionManager.addConnectionToCleanupQueue("connectionToDelete");
        assertDoesNotThrow(extensionManager::cleanup);
    }

    @Test
    void cleanupVirtualSchema() {
        extensionManager.addVirtualSchemaToCleanupQueue("virtualSchemaToDelete");
        assertDoesNotThrow(extensionManager::cleanup);
    }

    @Test
    void cleanupFile() throws IOException {
        final Path file = Files.createTempFile(getClass().getName(), ".tmp");
        extensionManager.addFileToCleanupQueue(file);
        assertTrue(Files.exists(file), "file still exists after adding it to the cleanup queue");
        extensionManager.cleanup();
        assertFalse(Files.exists(file), "file was deleted during cleanup");
    }

    @Test
    void fetchExtension() throws IOException {
        final String tempExtensionFileName = extensionManager.fetchExtension(URI.create(
                "https://extensions-internal.exasol.com/com.exasol/s3-document-files-virtual-schema/2.6.2/s3-vs-extension.js"));
        final Path file = extensionManager.extensionFolder.resolve(tempExtensionFileName);
        assertAll(() -> assertTrue(Files.exists(file), "file downloaded"),
                () -> assertThat(Files.size(file), equalTo(20875L)));

        extensionManager.cleanup();
        assertFalse(Files.exists(file), "file was deleted during cleanup");
    }

    @Test
    void fetchExtensionFails() throws IOException {
        final URI uri = URI.create("https://invalid-url");
        final UncheckedIOException exception = assertThrows(UncheckedIOException.class,
                () -> extensionManager.fetchExtension(uri));
        assertThat(exception.getMessage(), startsWith("E-EMIT-29: Failed to download"));
    }

    @Test
    void createExtensionSchema() {
        final ExasolSchema schema = extensionManager.createExtensionSchema();
        assertThat(schema, not(nullValue()));
    }

    @Test
    void cleanup() throws SQLException {
        extensionManager.createExtensionSchema();
        extensionManager.cleanup();
        assertNoSchemaExists();
    }

    private void assertNoSchemaExists() throws SQLException {
        try (PreparedStatement statement = exasolTestSetup.createConnection()
                .prepareStatement("select schema_name from EXA_ALL_SCHEMAS")) {
            assertThat(statement.executeQuery(), table("VARCHAR").matches());
        }
    }
}
