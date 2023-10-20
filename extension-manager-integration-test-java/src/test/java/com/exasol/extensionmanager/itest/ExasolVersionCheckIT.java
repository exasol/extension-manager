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
    void testAssertExasolVersion8DoesNotThrow() {
        assertDoesNotThrow(() -> ExasolVersionCheck.assertExasolVersion8(exasolTestSetup));
    }

    @Test
    void testAssumeExasolVersion8DoesNotThrow() {
        assertDoesNotThrow(() -> ExasolVersionCheck.assumeExasolVersion8(exasolTestSetup));
    }
}
