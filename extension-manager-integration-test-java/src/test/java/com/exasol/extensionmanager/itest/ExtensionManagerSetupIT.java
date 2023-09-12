package com.exasol.extensionmanager.itest;

import static com.exasol.extensionmanager.itest.IntegrationTestCommon.BUILT_EXTENSION_JS;
import static com.exasol.extensionmanager.itest.IntegrationTestCommon.TESTING_EXTENSION_SOURCE_DIR;
import static com.exasol.matcher.ResultSetStructureMatcher.table;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.not;
import static org.hamcrest.Matchers.nullValue;
import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;

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
        extensionManager = ExtensionManagerSetup.create(exasolTestSetup,
                ExtensionBuilder.createDefaultNpmBuilder(TESTING_EXTENSION_SOURCE_DIR, BUILT_EXTENSION_JS));
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
