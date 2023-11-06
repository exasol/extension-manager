package com.exasol.extensionmanager.itest.base;

import static com.exasol.matcher.ResultSetStructureMatcher.table;
import static java.util.Collections.emptyList;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.assertAll;
import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;

import java.util.*;
import java.util.logging.Logger;

import org.junit.jupiter.api.*;

import com.exasol.extensionmanager.client.model.*;
import com.exasol.extensionmanager.itest.*;

/**
 * This is a base class for Extension integration tests that already contains some basic tests for
 * installing/listing/uninstalling extensions and creating/listing/deleting instances.
 */
public abstract class AbstractExtensionIT {
    private static final Logger LOG = Logger.getLogger(AbstractExtensionIT.class.getName());
    private final ExtensionITConfig config;

    protected AbstractExtensionIT() {
        this.config = createConfig();
    }

    /**
     * Creates a new configuration for the integration tests.
     * 
     * @return new configuration
     */
    protected abstract ExtensionITConfig createConfig();

    protected abstract ExtensionManagerSetup getSetup();

    protected abstract void assertScriptsExist();

    protected abstract void prepareInstance();

    protected abstract void verifyVirtualTableContainsData(final String virtualSchemaName);

    /**
     * Create the same {@code SCRIPT}s as the extension would do. This is used to check that the extension also detects
     * manually created scripts.
     */
    protected abstract void createScripts();

    /**
     * Create valid parameters for a new instance. The instance/virtual schema name will be added automatically and is
     * not required here.
     *
     * @return valid parameters for a new instance
     */
    protected abstract Collection<ParameterValue> createValidParameterValues();

    @BeforeEach
    void logTestName(final TestInfo testInfo) {
        LOG.info(">>> " + testInfo.getDisplayName());
    }

    @AfterEach
    void cleanup() {
        getSetup().cleanup();
    }

    @Test
    public void listExtensions() {
        final List<ExtensionsResponseExtension> extensions = getSetup().client().getExtensions();
        assertAll(() -> assertThat(extensions, hasSize(1)), //
                () -> assertThat(extensions.get(0).getName(), equalTo(config.getExtensionName())),
                () -> assertThat(extensions.get(0).getInstallableVersions().get(0).getName(),
                        equalTo(config.getCurrentVersion())),
                () -> assertThat(extensions.get(0).getInstallableVersions().get(0).isLatest(), is(true)),
                () -> assertThat(extensions.get(0).getInstallableVersions().get(0).isDeprecated(), is(false)),
                () -> assertThat(extensions.get(0).getDescription(), equalTo(config.getExtensionDescription())));
    }

    @Test
    public void listInstallationsEmpty() {
        final List<InstallationsResponseInstallation> installations = getSetup().client().getInstallations();
        assertThat(installations, hasSize(0));
    }

    @Test
    public void listInstallationsFindsMatchingScripts() {
        createScripts();
        final List<InstallationsResponseInstallation> installations = getSetup().client().getInstallations();
        assertAll(() -> assertThat(installations, hasSize(1)), //
                () -> assertThat(installations.get(0).getName(), equalTo(config.getExtensionName())),
                () -> assertThat(installations.get(0).getVersion(), equalTo(config.getCurrentVersion())));
    }

    @Test
    public void listInstallationsFindsOwnInstallation() {
        getSetup().client().install();
        final List<InstallationsResponseInstallation> installations = getSetup().client().getInstallations();
        assertAll(() -> assertThat(installations, hasSize(1)), //
                () -> assertThat(installations.get(0).getName(), equalTo(config.getExtensionName())),
                () -> assertThat(installations.get(0).getVersion(), equalTo(config.getCurrentVersion())));
    }

    @Test
    public void getExtensionDetailsFailsForUnknownVersion() {
        getSetup().client().assertRequestFails(() -> getSetup().client().getExtensionDetails("unknownVersion"),
                equalTo("Version 'unknownVersion' not supported, can only use '" + config.getCurrentVersion() + "'."),
                equalTo(404));
    }

    @Test
    public void getExtensionDetailsSuccess() {
        final ExtensionDetailsResponse extensionDetails = getSetup().client()
                .getExtensionDetails(config.getCurrentVersion());
        final List<ParamDefinition> parameters = extensionDetails.getParameterDefinitions();
        final ParamDefinition param1 = new ParamDefinition().id("base-vs.virtual-schema-name")
                .name("Virtual Schema name").definition(Map.of( //
                        "id", "base-vs.virtual-schema-name", //
                        "name", "Virtual Schema name", //
                        "description", "Name for the new virtual schema", //
                        "placeholder", "MY_VIRTUAL_SCHEMA", //
                        "regex", "[a-zA-Z_]+", //
                        "required", true, //
                        "type", "string"));
        assertAll(() -> assertThat(extensionDetails.getId(), equalTo(config.getExtensionId())),
                () -> assertThat(extensionDetails.getVersion(), equalTo(config.getCurrentVersion())),
                () -> assertThat(parameters, hasSize(config.getExpectedParameterCount())),
                () -> assertThat(parameters.get(0), equalTo(param1)));
    }

    @Test
    public void installCreatesScripts() {
        getSetup().client().install();
        assertScriptsExist();
    }

    @Test
    public void installWorksIfCalledTwice() {
        getSetup().client().install();
        getSetup().client().install();
        assertScriptsExist();
    }

    @Test
    public void createInstanceFailsWithoutRequiredParameters() {
        final ExtensionManagerClient client = getSetup().client();
        client.install();
        client.assertRequestFails(() -> client.createInstance(emptyList()), startsWith(
                "invalid parameters: Failed to validate parameter 'Virtual Schema name' (base-vs.virtual-schema-name): This is a required parameter."),
                equalTo(400));
    }

    @Test
    public void uninstallSucceedsForNonExistingInstallation() {
        assertDoesNotThrow(() -> getSetup().client().uninstall());
    }

    @Test
    public void uninstallRemovesAdapters() {
        getSetup().client().install();
        assertAll(() -> assertScriptsExist(), //
                () -> assertThat(getSetup().client().getInstallations(), hasSize(1)));
        getSetup().client().uninstall(config.getCurrentVersion());
        assertAll(() -> assertThat(getSetup().client().getInstallations(), is(empty())),
                () -> getSetup().exasolMetadata().assertNoScripts());
    }

    @Test
    public void upgradeFailsWhenNotInstalled() {
        getSetup().client().assertRequestFails(() -> getSetup().client().upgrade(), //
                allOf(startsWith("Not all required scripts are installed: Validation failed: Script"),
                        endsWith("is missing")),
                equalTo(412));
    }

    @Test
    public void upgradeFailsWhenAlreadyUpToDate() {
        getSetup().client().install();
        getSetup().client().assertRequestFails(() -> getSetup().client().upgrade(),
                "Extension is already installed in latest version " + config.getCurrentVersion(), 412);
    }

    @Test
    public void upgradeFromPreviousVersion() {
        final PreviousExtensionVersion previousVersion = createPreviousVersion();
        previousVersion.prepare();
        previousVersion.install();
        prepareInstance();
        final String virtualSchemaName = "my_VS";
        createInstance(previousVersion.getExtensionId(), config.getPreviousVersion(), virtualSchemaName);
        verifyVirtualTableContainsData("my_VS");
        assertInstalledVersion(config.getPreviousVersion(), previousVersion);
        previousVersion.upgrade();
        assertInstalledVersion(config.getCurrentVersion(), previousVersion);
        verifyVirtualTableContainsData("my_VS");
    }

    private PreviousExtensionVersion createPreviousVersion() {
        return getSetup().previousVersionManager().newVersion().currentVersion(config.getCurrentVersion()) //
                .previousVersion(config.getPreviousVersion()) //
                .adapterFileName(config.getPreviousVersionJarFile()) //
                .extensionFileName(config.getExtensionId()) //
                .project(config.getProjectName()) //
                .build();
    }

    private void assertInstalledVersion(final String expectedVersion, final PreviousExtensionVersion previousVersion) {
        // The extension is installed twice (previous and current version), so each one returns one installation.
        assertThat(getSetup().client().getInstallations(),
                containsInAnyOrder(
                        new InstallationsResponseInstallation().name(config.getExtensionName()).version(expectedVersion)
                                .id(config.getExtensionId()), //
                        new InstallationsResponseInstallation().name(config.getExtensionName()).version(expectedVersion)
                                .id(previousVersion.getExtensionId())));
    }

    @Test
    public void virtualSchemaWorks() {
        getSetup().client().install();
        prepareInstance();
        createInstance("my_VS");
        verifyVirtualTableContainsData("my_VS");
    }

    @Test
    public void listingInstancesNoVSExists() {
        assertThat(getSetup().client().listInstances(), hasSize(0));
    }

    @Test
    public void listInstances() {
        getSetup().client().install();
        final String name = "my_virtual_SCHEMA";
        createInstance(name);
        assertThat(getSetup().client().listInstances(config.getCurrentVersion()),
                allOf(hasSize(1), equalTo(List.of(new Instance().id(name).name(name)))));
    }

    @Test
    public void createInstanceCreatesDbObjects() {
        getSetup().client().install();
        final String name = "my_virtual_SCHEMA";
        createInstance(name);

        getSetup().exasolMetadata()
                .assertConnection(table()
                        .row("MY_VIRTUAL_SCHEMA_CONNECTION", "Created by Extension Manager for "
                                + config.getExtensionName() + " v" + config.getCurrentVersion() + " my_virtual_SCHEMA")
                        .matches());
        getSetup().exasolMetadata().assertVirtualSchema(table()
                .row("my_virtual_SCHEMA", "SYS", "EXA_EXTENSIONS", not(emptyOrNullString()), not(emptyOrNullString()))
                .matches());
        assertThat(getSetup().client().listInstances(),
                allOf(hasSize(1), equalTo(List.of(new Instance().id(name).name(name)))));
    }

    @Test
    public void createTwoInstances() {
        getSetup().client().install();
        createInstance("vs1");
        createInstance("vs2");

        assertAll(
                () -> getSetup().exasolMetadata()
                        .assertConnection(table()
                                .row("VS1_CONNECTION",
                                        "Created by Extension Manager for " + config.getExtensionName() + " v"
                                                + config.getCurrentVersion() + " vs1")
                                .row("VS2_CONNECTION",
                                        "Created by Extension Manager for " + config.getExtensionName() + " v"
                                                + config.getCurrentVersion() + " vs2")
                                .matches()),
                () -> getSetup().exasolMetadata()
                        .assertVirtualSchema(table()
                                .row("vs1", "SYS", "EXA_EXTENSIONS", not(emptyOrNullString()), not(emptyOrNullString()))
                                .row("vs2", "SYS", "EXA_EXTENSIONS", not(emptyOrNullString()), not(emptyOrNullString()))
                                .matches()),

                () -> assertThat(getSetup().client().listInstances(), allOf(hasSize(2),
                        equalTo(List.of(new Instance().id("vs1").name("vs1"), new Instance().id("vs2").name("vs2"))))));
    }

    @Test
    public void createInstanceWithSingleQuote() {
        getSetup().client().install();
        createInstance("Quoted'schema");
        assertAll(
                () -> getSetup().exasolMetadata()
                        .assertConnection(table().row("QUOTED'SCHEMA_CONNECTION",
                                "Created by Extension Manager for S3 Virtual Schema v" + config.getCurrentVersion()
                                        + " Quoted'schema")
                                .matches()),
                () -> getSetup().exasolMetadata().assertVirtualSchema(table()
                        .row("Quoted'schema", "SYS", "EXA_EXTENSIONS", "S3_FILES_ADAPTER", not(emptyOrNullString()))
                        .matches()));
    }

    @Test
    public void deleteNonExistingInstance() {
        assertDoesNotThrow(() -> getSetup().client().deleteInstance("no-such-instance"));
    }

    @Test
    public void deleteFailsForUnknownVersion() {
        getSetup().client().assertRequestFails(
                () -> getSetup().client().deleteInstance("unknownVersion", "no-such-instance"),
                equalTo("Version 'unknownVersion' not supported, can only use '" + config.getCurrentVersion() + "'."),
                equalTo(404));
    }

    @Test
    public void deleteExistingInstance() {
        getSetup().client().install();
        createInstance("vs1");
        final List<Instance> instances = getSetup().client().listInstances();
        assertThat(instances, hasSize(1));
        getSetup().client().deleteInstance(instances.get(0).getId());

        assertAll(() -> assertThat(getSetup().client().listInstances(), is(empty())),
                () -> getSetup().exasolMetadata().assertNoConnections(),
                () -> getSetup().exasolMetadata().assertNoVirtualSchema());
    }

    private void createInstance(final String virtualSchemaName) {
        createInstance(config.getExtensionId(), config.getCurrentVersion(), virtualSchemaName);
    }

    private void createInstance(final String extensionId, final String extensionVersion,
            final String virtualSchemaName) {
        getSetup().addVirtualSchemaToCleanupQueue(virtualSchemaName);
        getSetup().addConnectionToCleanupQueue(virtualSchemaName.toUpperCase() + "_CONNECTION");
        final String instanceName = getSetup().client().createInstance(extensionId, extensionVersion,
                createValidParameters(virtualSchemaName));
        assertThat(instanceName, equalTo(virtualSchemaName));
    }

    private List<ParameterValue> createValidParameters(final String virtualSchemaName) {
        final List<ParameterValue> parameters = new ArrayList<>();
        parameters.add(param("base-vs.virtual-schema-name", virtualSchemaName));
        parameters.addAll(createValidParameterValues());
        return parameters;
    }

    protected ParameterValue param(final String name, final String value) {
        return new ParameterValue().name(name).value(value);
    }
}
