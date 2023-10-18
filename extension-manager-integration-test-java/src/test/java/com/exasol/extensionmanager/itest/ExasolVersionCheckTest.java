package com.exasol.extensionmanager.itest;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.when;

import java.sql.*;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import org.opentest4j.TestAbortedException;

import com.exasol.exasoltestsetup.ExasolTestSetup;

@ExtendWith(MockitoExtension.class)
class ExasolVersionCheckTest {

    @Mock
    ExasolTestSetup exasolTestSetupMock;
    @Mock
    Connection connectionMock;
    @Mock
    Statement statementMock;
    @Mock
    ResultSet resultSetMock;

    @Test
    void testAssertExasolVersion8DoesNotThrow() throws SQLException {
        simulateExasolVersion("8");
        assertDoesNotThrow(() -> ExasolVersionCheck.assertExasolVersion8(exasolTestSetupMock));
    }

    @Test
    void testAssertExasolVersion8Throws() throws SQLException {
        simulateExasolVersion("7");
        final AssertionError error = assertThrows(AssertionError.class,
                () -> ExasolVersionCheck.assertExasolVersion8(exasolTestSetupMock));
        assertThat(error.getMessage(), equalTo("Exasol version ==> expected: <8> but was: <7>"));
    }

    @Test
    void testAssumeExasolVersion8DoesNotThrow() throws SQLException {
        simulateExasolVersion("8");
        assertDoesNotThrow(() -> ExasolVersionCheck.assumeExasolVersion8(exasolTestSetupMock));
    }

    @Test
    void testAssumeExasolVersion8Throws() throws SQLException {
        simulateExasolVersion("7");
        final TestAbortedException error = assertThrows(TestAbortedException.class,
                () -> ExasolVersionCheck.assumeExasolVersion8(exasolTestSetupMock));
        assertThat(error.getMessage(), equalTo("Assumption failed: Expected Exasol version 8 but got '7'"));
    }

    private void simulateExasolVersion(final String version) throws SQLException {
        when(exasolTestSetupMock.createConnection()).thenReturn(connectionMock);
        when(connectionMock.createStatement()).thenReturn(statementMock);
        when(statementMock.executeQuery("SELECT PARAM_VALUE FROM SYS.EXA_METADATA WHERE PARAM_NAME='databaseMajorVersion'")).thenReturn(resultSetMock);
        when(resultSetMock.next()).thenReturn(true);
        when(resultSetMock.getString(1)).thenReturn(version);
    }
}
