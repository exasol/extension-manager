package com.exasol.extensionmanager.itest;

import static org.junit.jupiter.api.Assertions.assertEquals;
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
     * Ensures that the given test setup is connected to a Exasol DB in version 8.
     * 
     * @param testSetup test setup to check
     * @throws AssertionError if the major version number is not {@code 8}
     */
    public static void assertExasolVersion8(final ExasolTestSetup testSetup) {
        final String version = getExasolMajorVersion(testSetup);
        assertEquals("8", version, "Exasol version");
    }

    /**
     * Assumes that the given test setup is connected to a Exasol DB in version 8.
     * 
     * @param testSetup test setup to check
     * @throws org.opentest4j.TestAbortedException if the major version number is not {@code 8}
     */
    public static void assumeExasolVersion8(final ExasolTestSetup testSetup) {
        final String version = getExasolMajorVersion(testSetup);
        assumeTrue("8".equals(version), "Expected Exasol version 8 but got '" + version + "'");
    }

    static String getExasolMajorVersion(final ExasolTestSetup testSetup) {
        try (Statement stmt = testSetup.createConnection().createStatement()) {
            final ResultSet result = stmt
                    .executeQuery("SELECT PARAM_VALUE FROM SYS.EXA_METADATA WHERE PARAM_NAME='databaseMajorVersion'");
            assertTrue(result.next(), "no result");
            return result.getString(1);
        } catch (final SQLException exception) {
            throw new IllegalStateException("Failed to query Exasol version: " + exception.getMessage(), exception);
        }
    }
}
