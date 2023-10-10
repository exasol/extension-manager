package com.exasol.extensionmanager.itest;

import static com.exasol.extensionmanager.itest.IntegrationTestCommon.*;
import static java.util.Collections.emptyList;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertThrows;

import java.io.IOException;
import java.nio.file.*;
import java.sql.Connection;
import java.sql.SQLException;
import java.util.List;
import java.util.Map;

import org.junit.jupiter.api.*;

import com.exasol.exasoltestsetup.ExasolTestSetup;
import com.exasol.extensionmanager.client.invoker.ApiException;
import com.exasol.extensionmanager.client.model.*;
import com.exasol.extensionmanager.itest.builder.ExtensionBuilder;

// [itest -> dsn~eitfj-access-extension-manager-rest-interface~1]
// [itest -> dsn~eitfj-start-extension-manager~1]
class ExtensionManagerClientIT {

    private static final String EXTENSION_VERSION = "0.0.0";
    private static ExasolTestSetup exasolTestSetup;
    private static Connection connection;
    private static ExtensionManagerClient client;
    private static ExtensionManagerSetup extensionManager;

    @BeforeAll
    static void setup() throws SQLException, IOException {
        exasolTestSetup = IntegrationTestCommon.createExasolTestSetup();
        connection = exasolTestSetup.createConnection();
        createTestConfigFile();
        extensionManager = ExtensionManagerSetup.create(exasolTestSetup,
                ExtensionBuilder.createDefaultNpmBuilder(TESTING_EXTENSION_SOURCE_DIR, BUILT_EXTENSION_JS));
        client = extensionManager.client();
    }

    private static void createTestConfigFile() throws IOException {
        final Path extensionManagerProjectPath = Paths.get("..").toAbsolutePath().normalize();
        Files.writeString(IntegrationTestCommon.CONFIG_FILE,
                "localExtensionManager = " + extensionManagerProjectPath.toString() + "\n");
    }

    @AfterAll
    static void tearDown() throws Exception {
        connection.close();
        extensionManager.close();
        exasolTestSetup.close();
        Files.delete(IntegrationTestCommon.CONFIG_FILE);
    }

    @Test
    void listExtensions() {
        final ExtensionsResponseExtension expected = new ExtensionsResponseExtension().id(EXTENSION_ID)
                .name("Testing Extension").category("testing")
                .description("Extension for testing EM integration test setup").addInstallableVersionsItem(
                        new ExtensionVersion().name(EXTENSION_VERSION).latest(true).deprecated(false));
        assertThat(client.getExtensions(), contains(expected));
    }

    @Test
    void getInstallations() {
        final InstallationsResponseInstallation expected = new InstallationsResponseInstallation().id(EXTENSION_ID)
                .name("Testing Extension").version(EXTENSION_VERSION);
        assertThat(client.getInstallations(), contains(expected));
    }

    @Test
    void getDetails() {
        final ExtensionDetailsResponse expected = new ExtensionDetailsResponse().id(EXTENSION_ID)
                .version(EXTENSION_VERSION)
                .addParameterDefinitionsItem(new ParamDefinition().id("param1").name("Param 1")
                        .definition(Map.of("id", "param1", "name", "Param 1", "type", "string", "required", true)));
        assertThat(client.getExtensionDetails(EXTENSION_VERSION), equalTo(expected));
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
    void uninstallFailsBecauseInstanceExists() {
        client.assertRequestFails(() -> client.uninstall(),
                "cannot uninstall extension because 1 instance(s) still exist: Instance 1", 400);
    }

    @Test
    void upgrade() {
        assertDoesNotThrow(() -> client.upgrade());
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
                "invalid parameters: Failed to validate parameter 'Param 1': This is a required parameter.", 400);
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
}
