package com.exasol.extensionmanager.itest;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.empty;

import java.nio.file.Path;
import java.nio.file.Paths;
import java.sql.Connection;
import java.sql.SQLException;

import org.junit.jupiter.api.*;

import com.exasol.dbbuilder.dialects.exasol.ExasolObjectConfiguration;
import com.exasol.dbbuilder.dialects.exasol.ExasolObjectFactory;
import com.exasol.exasoltestsetup.ExasolTestSetup;
import com.exasol.exasoltestsetup.ExasolTestSetupFactory;
import com.exasol.udfdebugging.UdfTestSetup;

class ExtensionManagerSetupIT {
    private static ExasolTestSetup exasolTestSetup;

    @BeforeAll
    static void setupExasol() {
        exasolTestSetup = new ExasolTestSetupFactory(Path.of("dummy-config")).getTestSetup();
    }

    @AfterAll
    static void tearDownExasol() throws Exception {
        exasolTestSetup.close();
    }

    private Connection connection;
    private UdfTestSetup udfTestSetup;
    private ExtensionManagerSetup extensionManager;

    @BeforeEach
    void setup() throws SQLException {
        this.connection = exasolTestSetup.createConnection();
        this.udfTestSetup = new UdfTestSetup(exasolTestSetup, connection);
        final ExasolObjectFactory exasolObjectFactory = new ExasolObjectFactory(connection,
                ExasolObjectConfiguration.builder().withJvmOptions(udfTestSetup.getJvmOptions()).build());
        extensionManager = ExtensionManagerSetup.create(exasolTestSetup, exasolObjectFactory,
                Paths.get("testing-extension"));
    }

    @Test
    void listExtensions() {
        assertThat(extensionManager.client().getExtensions(), empty());
    }

    @AfterEach
    void teardown() throws SQLException {
        udfTestSetup.close();
        connection.close();
    }
}
