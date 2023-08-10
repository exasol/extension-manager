package com.exasol.extensionmanager.itest;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.notNullValue;
import static org.junit.jupiter.api.Assertions.assertFalse;
import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.mockito.Mockito.lenient;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.sql.*;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.api.io.TempDir;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import com.exasol.exasoltestsetup.ExasolTestSetup;
import com.exasol.exasoltestsetup.SqlConnectionInfo;
import com.exasol.extensionmanager.itest.builder.ExtensionBuilder;

@ExtendWith(MockitoExtension.class)
class ExtensionManagerSetupTest {

    @Mock
    private ExasolTestSetup exasolTestSetupMock;
    @Mock
    private ExtensionBuilder extensionBuilderMock;
    @Mock
    private Connection dbConnectionMock;
    @Mock
    private Statement statementMock;
    private ExtensionManagerSetup extensionManager;

    @TempDir
    Path tempDir;

    @BeforeEach
    void setup() throws SQLException, IOException {
        final Path file = Files.createFile(tempDir.resolve("built-extension.js"));
        lenient().when(extensionBuilderMock.getExtensionFile()).thenReturn(file);
        lenient().when(exasolTestSetupMock.getConnectionInfo())
                .thenReturn(new SqlConnectionInfo("host", 8563, "user", "pass"));
        lenient().when(exasolTestSetupMock.createConnection()).thenReturn(dbConnectionMock);
        lenient().when(dbConnectionMock.createStatement()).thenReturn(statementMock);
        extensionManager = ExtensionManagerSetup.create(exasolTestSetupMock, extensionBuilderMock);
    }

    @Test
    void cleanupFile() throws IOException {
        final Path file = Files.createFile(tempDir.resolve("temp"));
        extensionManager.addFileToCleanupQueue(file);
        assertTrue(Files.exists(file), "file still exists after adding it to the cleanup queue");
        extensionManager.cleanup();
        assertFalse(Files.exists(file), "file was deleted during cleanup");
    }

    @Test
    void previousVersionManager() {
        assertThat(extensionManager.previousVersionManager(), notNullValue());
    }
}
