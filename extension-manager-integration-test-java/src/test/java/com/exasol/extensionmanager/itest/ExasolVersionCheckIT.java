package com.exasol.extensionmanager.itest;

import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;

import org.junit.jupiter.api.*;

import com.exasol.exasoltestsetup.ExasolTestSetup;

class ExasolVersionCheckIT {

    private static ExasolTestSetup exasolTestSetup;

    @BeforeAll
    static void setup() {
        exasolTestSetup = IntegrationTestCommon.createExasolTestSetup();
    }

    @AfterAll
    static void tearDown() throws Exception {
        exasolTestSetup.close();
    }

    @Test
    void testAssertSupportedExasolVersionDoesNotThrow() {
        assertDoesNotThrow(() -> ExasolVersionCheck.assertExasolVersionSupported(exasolTestSetup));
    }

    @Test
    void testAssumeSupportedExasolVersionDoesNotThrow() {
        assertDoesNotThrow(() -> ExasolVersionCheck.assumeSupportedExasolVersion(exasolTestSetup));
    }
}
