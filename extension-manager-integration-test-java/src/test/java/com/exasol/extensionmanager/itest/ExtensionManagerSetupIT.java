package com.exasol.extensionmanager.itest;

import static java.util.Collections.emptyList;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertThrows;

import java.nio.file.Path;
import java.nio.file.Paths;
import java.sql.Connection;
import java.sql.SQLException;
import java.time.Duration;
import java.util.List;
import java.util.Map;

import org.junit.jupiter.api.*;

import com.exasol.dbbuilder.dialects.exasol.ExasolObjectConfiguration;
import com.exasol.dbbuilder.dialects.exasol.ExasolObjectFactory;
import com.exasol.exasoltestsetup.ExasolTestSetup;
import com.exasol.exasoltestsetup.ExasolTestSetupFactory;
import com.exasol.extensionmanager.client.invoker.ApiException;
import com.exasol.extensionmanager.client.model.*;
import com.exasol.extensionmanager.itest.builder.ExtensionBuilder;
import com.exasol.extensionmanager.itest.process.SimpleProcess;
import com.exasol.udfdebugging.UdfTestSetup;

class ExtensionManagerSetupIT {
    private static final Path TESTING_EXTENSION_SOURCE_DIR = Paths.get("testing-extension").toAbsolutePath();
    private static ExasolTestSetup exasolTestSetup;

    @BeforeAll
    static void setupExasol() {
        exasolTestSetup = new ExasolTestSetupFactory(Path.of("dummy-config")).getTestSetup();
        SimpleProcess.start(TESTING_EXTENSION_SOURCE_DIR, List.of("npm", "install"), Duration.ofSeconds(60));
    }

    @AfterAll
    static void tearDownExasol() throws Exception {
        exasolTestSetup.close();
    }

    private Connection connection;
    private UdfTestSetup udfTestSetup;
    private ExtensionManagerClient client;
    private ExtensionManagerSetup extensionManager;

    @BeforeEach
    void setup() throws SQLException {
        this.connection = exasolTestSetup.createConnection();
        this.udfTestSetup = new UdfTestSetup(exasolTestSetup, connection);
        final ExasolObjectFactory exasolObjectFactory = new ExasolObjectFactory(connection,
                ExasolObjectConfiguration.builder().withJvmOptions(udfTestSetup.getJvmOptions()).build());
        extensionManager = ExtensionManagerSetup.create(exasolTestSetup, exasolObjectFactory,
                ExtensionBuilder.createDefaultNpmBuilder(TESTING_EXTENSION_SOURCE_DIR,
                        TESTING_EXTENSION_SOURCE_DIR.resolve("dist/testing-extension.js")));
        client = extensionManager.client();
    }

    @Test
    void listExtensions() {
        final ExtensionsResponseExtension expected = new ExtensionsResponseExtension().id("testing-extension.js")
                .name("Testing Extension").description("Extension for testing EM integration test setup")
                .addInstallableVersionsItem(new ExtensionVersion().name("0.0.0").latest(true).deprecated(false));
        assertThat(client.getExtensions(), contains(expected));
    }

    @Test
    void getInstallations() {
        final InstallationsResponseInstallation expected = new InstallationsResponseInstallation()
                .name("Testing Extension").version("0.0.0");
        assertThat(client.getInstallations(), contains(expected));
    }

    @Test
    void getDetails() {
        final ExtensionDetailsResponse expected = new ExtensionDetailsResponse().id("testing-extension.js")
                .version("0.0.0").addParameterDefinitionsItem(new ParamDefinition().id("param1").name("Param 1")
                        .definition(Map.of("id", "param1", "name", "Param 1", "type", "string", "required", true)));
        assertThat(client.getExtensionDetails("0.0.0"), equalTo(expected));
    }

    @Test
    void install() {
        assertDoesNotThrow(() -> client.install());
    }

    @Test
    void installVersion() {
        assertDoesNotThrow(() -> client.install("otherVersion"));
    }

    @Test
    void createInstanceFailsForMissingParam() {
        final List<ParameterValue> params = emptyList();
        final ApiException exception = assertThrows(ApiException.class, () -> client.createInstance(params));
        assertThat(exception.getMessage(), containsString(
                "invalid parameters: Failed to validate parameter 'Param 1': This is a required parameter."));
    }

    @Test
    void createInstanceSucceeds() {
        final List<ParameterValue> params = List.of(new ParameterValue().name("param1").value("value1"));
        assertDoesNotThrow(() -> client.createInstance(params));
    }

    @Test
    void assertRequestFails() {
        client.assertRequestFails(() -> client.createInstance(emptyList()),
                equalTo("invalid parameters: Failed to validate parameter 'Param 1': This is a required parameter."),
                equalTo(400));
    }

    @Test
    void listInstances() {
        assertThat(client.listInstances(), contains(new Instance().id("instance-1").name("Instance 1")));
    }

    @Test
    void deleteInstance() {
        assertDoesNotThrow(() -> client.deleteInstance("instance-id"));
    }

    @Test
    void deleteInstanceWithVersion() {
        assertDoesNotThrow(() -> client.deleteInstance("version", "instance-id"));
    }

    @AfterEach
    void teardown() throws SQLException {
        udfTestSetup.close();
        connection.close();
    }
}
