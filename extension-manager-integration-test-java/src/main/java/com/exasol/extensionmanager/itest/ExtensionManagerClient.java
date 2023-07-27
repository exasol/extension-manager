package com.exasol.extensionmanager.itest;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.junit.jupiter.api.Assertions.assertAll;
import static org.junit.jupiter.api.Assertions.assertThrows;

import java.util.List;
import java.util.logging.Logger;

import org.hamcrest.Matcher;
import org.junit.jupiter.api.function.Executable;
import org.opentest4j.MultipleFailuresError;

import com.exasol.errorreporting.ExaError;
import com.exasol.exasoltestsetup.SqlConnectionInfo;
import com.exasol.extensionmanager.client.api.*;
import com.exasol.extensionmanager.client.invoker.ApiClient;
import com.exasol.extensionmanager.client.invoker.ApiException;
import com.exasol.extensionmanager.client.model.*;

import jakarta.json.JsonObject;
import jakarta.json.bind.Jsonb;
import jakarta.json.bind.JsonbBuilder;

/**
 * This class wraps the extension manager REST API and simplifies calling its endpoint methods:
 * <ul>
 * <li>Adds authentication header value and passes database connection parameters required by most requests.
 * <li>Adds extension ID parameter for the extension under test to requests when required.
 * <li>Optionally adds the current version of the extension under test to requests when required.
 * </ul>
 * The class also provides a convenient method {@link #assertRequestFails(Executable, Matcher, Matcher)} for verifying
 * that an API fails with expected error message and HTTP status code.
 */
public class ExtensionManagerClient {
    private static final Logger LOGGER = Logger.getLogger(ExtensionManagerClient.class.getName());
    private final ExtensionApi extensionClient;
    private final InstallationApi installationApi;
    private final InstanceApi instanceClient;
    private final SqlConnectionInfo dbConnectionInfo;

    private ExtensionManagerClient(final ExtensionApi extensionClient, final InstallationApi installationApi,
            final InstanceApi instanceClient, final SqlConnectionInfo dbConnectionInfo) {
        this.extensionClient = extensionClient;
        this.installationApi = installationApi;
        this.instanceClient = instanceClient;
        this.dbConnectionInfo = dbConnectionInfo;
    }

    static ExtensionManagerClient create(final String serverBasePath, final SqlConnectionInfo connectionInfo) {
        final ApiClient apiClient = createApiClient(serverBasePath, connectionInfo);
        return new ExtensionManagerClient(new ExtensionApi(apiClient), new InstallationApi(apiClient),
                new InstanceApi(apiClient), connectionInfo);
    }

    private static ApiClient createApiClient(final String serverBasePath, final SqlConnectionInfo connectionInfo) {
        final ApiClient apiClient = new ApiClient().setBasePath(serverBasePath);
        apiClient.setUsername(connectionInfo.getUser());
        apiClient.setPassword(connectionInfo.getPassword());
        return apiClient;
    }

    /**
     * Calls {@link ExtensionApi#listAvailableExtensions(String, Integer)}.
     * 
     * @return list of available extensions
     */
    public List<ExtensionsResponseExtension> getExtensions() {
        return this.extensionClient.listAvailableExtensions(getDbHost(), getDbPort()).getExtensions();
    }

    /**
     * Calls {@link InstallationApi#listInstalledExtensions(String, Integer)}.
     * 
     * @return list of installed extensions
     */
    public List<InstallationsResponseInstallation> getInstallations() {
        return this.installationApi.listInstalledExtensions(getDbHost(), getDbPort()).getInstallations();
    }

    /**
     * Calls {@link ExtensionApi#getExtensionDetails(String, String, String, Integer)}.
     * 
     * @param extensionVersion extension version
     * @return extension details
     */
    public ExtensionDetailsResponse getExtensionDetails(final String extensionVersion) {
        return getExtensionDetails(getExtension().getId(), extensionVersion);
    }

    private ExtensionDetailsResponse getExtensionDetails(final String extensionId, final String extensionVersion) {
        return this.extensionClient.getExtensionDetails(extensionId, extensionVersion, getDbHost(), getDbPort());
    }

    /**
     * Calls {@link ExtensionApi#installExtension(InstallExtensionRequest, String, Integer, String, String)}.
     * 
     * @param version extension version
     */
    public void install(final String version) {
        final Extension extension = getExtension();
        LOGGER.fine(() -> "Installing extension " + extension.getId() + " in version " + version);
        install(extension.getId(), version);
    }

    /**
     * Calls {@link ExtensionApi#installExtension(InstallExtensionRequest, String, Integer, String, String)}.
     */
    public void install() {
        install(getExtension().getCurrentVersion());
    }

    /**
     * Calls {@link ExtensionApi#installExtension(InstallExtensionRequest, String, Integer, String, String)}.
     * 
     * @param extensionId extension id
     * @param version     extension version
     */
    public void install(final String extensionId, final String extensionVersion) {
        this.extensionClient.installExtension(new InstallExtensionRequest(), getDbHost(), getDbPort(), extensionId,
                extensionVersion);
    }

    /**
     * Calls {@link InstallationApi#uninstallExtension(String, String, String, Integer)}.
     * 
     * @param extensionVersion the extension version
     */
    public void uninstall(final String extensionVersion) {
        this.uninstall(getExtension().getId(), extensionVersion);
    }

    /**
     * Calls {@link InstallationApi#uninstallExtension(String, String, String, Integer)} with the current version.
     */
    public void uninstall() {
        final Extension extension = getExtension();
        this.uninstall(extension.getId(), extension.getCurrentVersion());
    }

    private void uninstall(final String extensionId, final String extensionVersion) {
        this.installationApi.uninstallExtension(extensionId, extensionVersion, getDbHost(), getDbPort());
    }

    /**
     * Calls {@link InstallationApi#upgradeExtension(String, String, Integer)}.
     * 
     * @return upgrade response
     */
    public UpgradeExtensionResponse upgrade() {
        return this.upgrade(getExtension().getId());
    }

    /**
     * Calls {@link InstallationApi#upgradeExtension(String, String, Integer)}.
     * 
     * @param extensionId extension id
     * @return upgrade response
     */
    public UpgradeExtensionResponse upgrade(final String extensionId) {
        return this.installationApi.upgradeExtension(extensionId, getDbHost(), getDbPort());
    }

    /**
     * Calls {@link InstanceApi#createInstance(CreateInstanceRequest, String, Integer, String, String)}.
     * 
     * @param parameterValues parameter values for creating the instance
     * @return name of the new instance
     */
    public String createInstance(final List<ParameterValue> parameterValues) {
        final Extension extension = getExtension();
        return createInstance(extension.getId(), extension.getCurrentVersion(), parameterValues);
    }

    /**
     * Calls {@link InstanceApi#createInstance(CreateInstanceRequest, String, Integer, String, String)}.
     * 
     * @param extensionId      extension id
     * @param extensionVersion extension version
     * @param parameterValues  parameter name
     * @return name of the new instance
     */
    public String createInstance(final String extensionId, final String extensionVersion,
            final List<ParameterValue> parameterValues) {
        final CreateInstanceRequest request = new CreateInstanceRequest().parameterValues(parameterValues);
        return this.instanceClient.createInstance(request, getDbHost(), getDbPort(), extensionId, extensionVersion)
                .getInstanceName();
    }

    /**
     * Calls {@link InstanceApi#listInstances(String, String, String, Integer)}.
     * 
     * @return list of available instances
     */
    public List<Instance> listInstances() {
        final Extension extension = getExtension();
        return listInstances(extension.getCurrentVersion());
    }

    /**
     * Calls {@link InstanceApi#listInstances(String, String, String, Integer)}.
     * 
     * @param version extension version
     * @return list of available instances
     */
    public List<Instance> listInstances(final String version) {
        final Extension extension = getExtension();
        return listInstances(extension.getId(), version).getInstances();
    }

    private ListInstancesResponse listInstances(final String extensionId, final String extensionVersion) {
        return this.instanceClient.listInstances(extensionId, extensionVersion, getDbHost(), getDbPort());
    }

    /**
     * Calls {@link InstanceApi#deleteInstance(String, String, String, String, Integer)}.
     * 
     * @param version    extension version
     * @param instanceId instance id to delete
     */
    public void deleteInstance(final String version, final String instanceId) {
        final Extension extension = getExtension();
        deleteInstance(extension.getId(), version, instanceId);
    }

    /**
     * Calls {@link InstanceApi#deleteInstance(String, String, String, String, Integer)}.
     * 
     * @param instanceId instance id to delete
     */
    public void deleteInstance(final String instanceId) {
        final Extension extension = getExtension();
        deleteInstance(extension.getId(), extension.getCurrentVersion(), instanceId);
    }

    private void deleteInstance(final String extensionId, final String extensionVersion, final String instanceId) {
        this.instanceClient.deleteInstance(extensionId, extensionVersion, instanceId, getDbHost(), getDbPort());
    }

    /**
     * Verify that the given executable throws an {@link ApiException} with a given error message and HTTP status code.
     * 
     * @param executable     executable to run
     * @param messageMatcher {@link Matcher} for the error message field of the response
     * @param statusMatcher  {@link Matcher} for the HTTP status code of the response
     */
    public void assertRequestFails(final Executable executable, final Matcher<String> messageMatcher,
            final Matcher<Integer> statusMatcher) {
        final ApiException exception = assertThrows(ApiException.class, executable);
        final String errorMessage = exception.getMessage();
        final JsonObject error = parseErrorMessageJson(errorMessage);
        assertAll(() -> assertThat(error.getJsonString("message").getString(), messageMatcher),
                () -> assertThat(error.getJsonNumber("code").intValue(), statusMatcher));
    }

    private JsonObject parseErrorMessageJson(final String errorMessage) throws MultipleFailuresError {
        try (Jsonb jsonb = JsonbBuilder.create()) {
            return jsonb.fromJson(errorMessage, JsonObject.class);
        } catch (final Exception jsonbException) {
            throw new IllegalStateException(ExaError.messageBuilder("E-EMIT-15")
                    .message("Failed to parse error message {{error message}} as JSON")
                    .parameter("error message", errorMessage, "messaged to be parsed as JSON").ticketMitigation()
                    .toString(), jsonbException);
        }
    }

    /**
     * Verify that the given executable throws an {@link ApiException} with a given error message and HTTP status code.
     * 
     * @param executable      executable to run
     * @param expectedMessage expected response error message
     * @param expectedStatus  expected response status code
     */
    public void assertRequestFails(final Executable executable, final String expectedMessage,
            final int expectedStatus) {
        this.assertRequestFails(executable, equalTo(expectedMessage), equalTo(expectedStatus));
    }

    private ExtensionsResponseExtension getSingleExtension() {
        final List<ExtensionsResponseExtension> extensions = this.getExtensions();
        if (extensions.size() != 1) {
            throw new IllegalStateException(ExaError.messageBuilder("E-EMIT-28")
                    .message(
                            "Expected exactly one extension but found {{actual count}}: {{actual list of extensions}}.",
                            extensions.size(), extensions)
                    .mitigation("Check the extension manager log for errors loading the extension.").toString());
        }
        return extensions.get(0);
    }

    private Extension getExtension() {
        final ExtensionsResponseExtension extension = getSingleExtension();
        if (extension.getInstallableVersions().size() != 1) {
            throw new IllegalStateException(ExaError.messageBuilder("E-EMIT-16").message(
                    "Expected exactly one installable version for extension {{extension id}} but got {{actual list of versions}}.",
                    extension.getId(), extension.getInstallableVersions()).toString());
        }
        return new Extension(extension.getId(), extension.getInstallableVersions().get(0).getName());
    }

    private static class Extension {
        private final String id;
        private final String currentVersion;

        private Extension(final String id, final String currentVersion) {
            this.id = id;
            this.currentVersion = currentVersion;
        }

        public String getId() {
            return this.id;
        }

        public String getCurrentVersion() {
            return this.currentVersion;
        }
    }

    private String getDbHost() {
        return this.dbConnectionInfo.getHost();
    }

    private int getDbPort() {
        return this.dbConnectionInfo.getPort();
    }
}
