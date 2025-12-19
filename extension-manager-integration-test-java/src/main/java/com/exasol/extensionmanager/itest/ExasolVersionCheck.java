package com.exasol.extensionmanager.itest;

import static org.junit.jupiter.api.Assertions.assertTrue;
import static org.junit.jupiter.api.Assumptions.assumeTrue;

import java.sql.*;

import com.exasol.exasoltestsetup.ExasolTestSetup;

/**
 * This class contains static methods for checking the Exasol DB version.
 */
public class ExasolVersionCheck {
    private ExasolVersionCheck() {
        // Not instantiable
    }

    /**
     * Ensures that the given test setup is connected to a supported Exasol DB version.
     * <p>
     * This executes a SQL query using the given test setup to determine the major version of the database.
     * 
     * @param testSetup test setup to check
     * @throws AssertionError if the major version number is less than {@code 8}
     */
    public static void assertExasolVersionSupported(final ExasolTestSetup testSetup) {
        final int version = getExasolMajorVersion(testSetup);
        assertTrue(version >= 8, "Expected Exasol version >= 8 but got '" + version + "'");
    }

    /**
     * Assumes that the given test setup is connected to a supported Exasol DB version.
     * <p>
     * This executes a SQL query using the given test setup to determine the major version of the database.
     * 
     * @param testSetup test setup to check
     * @throws org.opentest4j.TestAbortedException if the major version number less than {@code 8}
     */
    public static void assumeSupportedExasolVersion(final ExasolTestSetup testSetup) {
        final int version = getExasolMajorVersion(testSetup);
        assumeTrue(version >= 8, "Expected Exasol version >= 8 but got '" + version + "'");
    }

    static int getExasolMajorVersion(final ExasolTestSetup testSetup) {
        try (Statement stmt = testSetup.createConnection().createStatement()) {
            final ResultSet result = stmt
                    .executeQuery("SELECT PARAM_VALUE FROM SYS.EXA_METADATA WHERE PARAM_NAME='databaseMajorVersion'");
            assertTrue(result.next(), "no result");
            return result.getInt(1);
        } catch (final SQLException exception) {
            throw new IllegalStateException("Failed to query Exasol version: " + exception.getMessage(), exception);
        }
    }
}
