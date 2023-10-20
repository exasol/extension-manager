package com.exasol.extensionmanager.itest;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.notNullValue;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.Mockito.lenient;
import static org.mockito.Mockito.when;

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
    ExasolTestSetup exasolTestSetupMock;
    @Mock
    ExtensionBuilder extensionBuilderMock;
    @Mock
    Connection dbConnectionMock;
    @Mock
    Statement statementMock;
    @Mock
    ResultSet resultSetMock;
    @TempDir
    Path tempDir;

    ExtensionManagerSetup extensionManager;

    @BeforeEach
    void setup() throws SQLException, IOException {
        final Path file = Files.createFile(tempDir.resolve("built-extension.js"));
        lenient().when(extensionBuilderMock.getExtensionFile()).thenReturn(file);
        lenient().when(exasolTestSetupMock.getConnectionInfo())
                .thenReturn(new SqlConnectionInfo("host", 8563, "user", "pass"));
        lenient().when(exasolTestSetupMock.createConnection()).thenReturn(dbConnectionMock);
        lenient().when(dbConnectionMock.createStatement()).thenReturn(statementMock);
        simulateExasolVersion("8");
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

    @Test
    void createFailsForOldVersion() throws SQLException {
        simulateExasolVersion("7");
        final AssertionError error = assertThrows(AssertionError.class,
                () -> ExtensionManagerSetup.create(exasolTestSetupMock, null));
        assertThat(error.getMessage(), equalTo("Exasol version ==> expected: <8> but was: <7>"));
    }

    private void simulateExasolVersion(final String version) throws SQLException {
        when(statementMock.executeQuery("SELECT PARAM_VALUE FROM SYS.EXA_METADATA WHERE PARAM_NAME='databaseMajorVersion'")).thenReturn(resultSetMock);
        when(resultSetMock.next()).thenReturn(true);
        when(resultSetMock.getString(1)).thenReturn(version);
    }
}
