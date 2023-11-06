package com.exasol.extensionmanager.itest;

import static com.exasol.matcher.ResultSetStructureMatcher.table;
import static java.util.Collections.emptyList;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.assertAll;
import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;

import java.io.FileNotFoundException;
import java.net.URISyntaxException;
import java.util.*;
import java.util.concurrent.TimeoutException;

import org.junit.jupiter.api.Test;

import com.exasol.bucketfs.BucketAccessException;
import com.exasol.extensionmanager.client.model.*;

/**
 * This is a base class for Extension integration tests that already contains some basic tests for
 * installing/listing/uninstalling extensions and creating/listing/deleting instances.
 */
public abstract class AbstractExtensionIT {

    protected abstract ExtensionManagerSetup getSetup();

    /**
     * Get the ID of this extension, e.g. {@code s3-vs-extension.js}.
     * 
     * @return extension ID
     */
    protected abstract String getExtensionId();

    /**
     * Get the user visible name of this extension, e.g. {@code S3 Virtual Schema}.
     * 
     * @return extension name
     */
    protected abstract String getName();

    /**
     * Get the total number of parameters for this extension, incl. virtual schema name.
     * 
     * @return parameter count
     */
    protected abstract int getExpectedParameterCount();

    /**
     * Get the user visible description of this extension, e.g. {@code Virtual Schema for document files on AWS S3}.
     * 
     * @return extension description
     */
    protected abstract String getDescription();

    /**
     * Get the current version of this extension, e.g. {@code 1.2.3}.
     * 
     * @return current version
     */
    protected abstract String getCurrentVersion();

    /**
     * Get the previous version of this extension, e.g. {@code 1.2.2}.
     * <p>
     * This may be {@code null} if you are just creating the first version of the extension. Once you release a second
     * version, update this to return the previous version.
     * 
     * @return previous version
     */
    protected abstract String getPreviousVersion();

    /**
     * Get the previous version's JAR file name of this extension, e.g.
     * {@code document-files-virtual-schema-dist-7.3.6-s3-1.2.3.jar}.
     * <p>
     * This may be {@code null} if you are just creating the first version of the extension. Once you release a second
     * version, update this to return the JAR file name.
     * 
     * @return previous version's JAR file name
     */
    protected abstract String getPreviousVersionJarFile();

    /**
     * Get the project name, e.g. {@code s3-document-files-virtual-schema}.
     * 
     * @return project name
     */
    protected abstract String getProjectName();

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

    @Test
    void listExtensions() {
        final List<ExtensionsResponseExtension> extensions = getSetup().client().getExtensions();
        assertAll(() -> assertThat(extensions, hasSize(1)), //
                () -> assertThat(extensions.get(0).getName(), equalTo(getName())),
                () -> assertThat(extensions.get(0).getInstallableVersions().get(0).getName(),
                        equalTo(getCurrentVersion())),
                () -> assertThat(extensions.get(0).getInstallableVersions().get(0).isLatest(), is(true)),
                () -> assertThat(extensions.get(0).getInstallableVersions().get(0).isDeprecated(), is(false)),
                () -> assertThat(extensions.get(0).getDescription(), equalTo(getDescription())));
    }

    @Test
    void listInstallationsEmpty() {
        final List<InstallationsResponseInstallation> installations = getSetup().client().getInstallations();
        assertThat(installations, hasSize(0));
    }

    @Test
    void listInstallationsFindsMatchingScripts() {
        createScripts();
        final List<InstallationsResponseInstallation> installations = getSetup().client().getInstallations();
        assertAll(() -> assertThat(installations, hasSize(1)), //
                () -> assertThat(installations.get(0).getName(), equalTo(getName())),
                () -> assertThat(installations.get(0).getVersion(), equalTo(getCurrentVersion())));
    }

    @Test
    void listInstallationsFindsOwnInstallation() {
        getSetup().client().install();
        final List<InstallationsResponseInstallation> installations = getSetup().client().getInstallations();
        assertAll(() -> assertThat(installations, hasSize(1)), //
                () -> assertThat(installations.get(0).getName(), equalTo(getName())),
                () -> assertThat(installations.get(0).getVersion(), equalTo(getCurrentVersion())));
    }

    @Test
    void getExtensionDetailsFailsForUnknownVersion() {
        getSetup().client().assertRequestFails(() -> getSetup().client().getExtensionDetails("unknownVersion"),
                equalTo("Version 'unknownVersion' not supported, can only use '" + getCurrentVersion() + "'."),
                equalTo(404));
    }

    @Test
    void getExtensionDetailsSuccess() {
        final ExtensionDetailsResponse extensionDetails = getSetup().client().getExtensionDetails(getCurrentVersion());
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
        assertAll(() -> assertThat(extensionDetails.getId(), equalTo(getExtensionId())),
                () -> assertThat(extensionDetails.getVersion(), equalTo(getCurrentVersion())),
                () -> assertThat(parameters, hasSize(getExpectedParameterCount())),
                () -> assertThat(parameters.get(0), equalTo(param1)));
    }

    @Test
    void installCreatesScripts() {
        getSetup().client().install();
        assertScriptsExist();
    }

    @Test
    void installWorksIfCalledTwice() {
        getSetup().client().install();
        getSetup().client().install();
        assertScriptsExist();
    }

    @Test
    void createInstanceFailsWithoutRequiredParameters() {
        final ExtensionManagerClient client = getSetup().client();
        client.install();
        client.assertRequestFails(() -> client.createInstance(emptyList()), startsWith(
                "invalid parameters: Failed to validate parameter 'Virtual Schema name' (base-vs.virtual-schema-name): This is a required parameter."),
                equalTo(400));
    }

    @Test
    void uninstallSucceedsForNonExistingInstallation() {
        assertDoesNotThrow(() -> getSetup().client().uninstall());
    }

    @Test
    void uninstallRemovesAdapters() {
        getSetup().client().install();
        assertAll(() -> assertScriptsExist(), //
                () -> assertThat(getSetup().client().getInstallations(), hasSize(1)));
        getSetup().client().uninstall(getCurrentVersion());
        assertAll(() -> assertThat(getSetup().client().getInstallations(), is(empty())),
                () -> getSetup().exasolMetadata().assertNoScripts());
    }

    @Test
    void upgradeFailsWhenNotInstalled() {
        getSetup().client().assertRequestFails(() -> getSetup().client().upgrade(), //
                allOf(startsWith("Not all required scripts are installed: Validation failed: Script"),
                        endsWith("is missing")),
                equalTo(412));
    }

    @Test
    void upgradeFailsWhenAlreadyUpToDate() {
        getSetup().client().install();
        getSetup().client().assertRequestFails(() -> getSetup().client().upgrade(),
                "Extension is already installed in latest version " + getCurrentVersion(), 412);
    }

    @Test
    void upgradeFromPreviousVersion() throws InterruptedException, BucketAccessException, TimeoutException,
            FileNotFoundException, URISyntaxException {
        final PreviousExtensionVersion previousVersion = createPreviousVersion();
        previousVersion.prepare();
        previousVersion.install();
        prepareInstance();
        createInstance(getExtensionId(), getPreviousVersion(), "my_VS");
        verifyVirtualTableContainsData("my_VS");
        assertInstalledVersion(getPreviousVersion(), previousVersion);
        previousVersion.upgrade();
        assertInstalledVersion(getCurrentVersion(), previousVersion);
        verifyVirtualTableContainsData("my_VS");
    }

    private PreviousExtensionVersion createPreviousVersion() {
        return getSetup().previousVersionManager().newVersion().currentVersion(getCurrentVersion()) //
                .previousVersion(getPreviousVersion()) //
                .adapterFileName(getPreviousVersionJarFile()) //
                .extensionFileName(getExtensionId()) //
                .project(getProjectName()) //
                .build();
    }

    private void assertInstalledVersion(final String expectedVersion, final PreviousExtensionVersion previousVersion) {
        // The extension is installed twice (previous and current version), so each one returns one installation.
        assertThat(getSetup().client().getInstallations(),
                containsInAnyOrder(
                        new InstallationsResponseInstallation().name(getName()).version(expectedVersion)
                                .id(getExtensionId()), //
                        new InstallationsResponseInstallation().name(getName()).version(expectedVersion)
                                .id(previousVersion.getExtensionId())));
    }

    @Test
    void virtualSchemaWorks() {
        getSetup().client().install();
        prepareInstance();
        createInstance("my_VS");
        verifyVirtualTableContainsData("my_VS");
    }

    @Test
    void listingInstancesNoVSExists() {
        assertThat(getSetup().client().listInstances(), hasSize(0));
    }

    @Test
    void listInstances() {
        getSetup().client().install();
        final String name = "my_virtual_SCHEMA";
        createInstance(name);
        assertThat(getSetup().client().listInstances(getCurrentVersion()),
                allOf(hasSize(1), equalTo(List.of(new Instance().id(name).name(name)))));
    }

    @Test
    void createInstanceCreatesDbObjects() {
        getSetup().client().install();
        final String name = "my_virtual_SCHEMA";
        createInstance(name);

        getSetup().exasolMetadata().assertConnection(table().row("MY_VIRTUAL_SCHEMA_CONNECTION",
                "Created by Extension Manager for " + getName() + " v" + getCurrentVersion() + " my_virtual_SCHEMA")
                .matches());
        getSetup().exasolMetadata().assertVirtualSchema(table()
                .row("my_virtual_SCHEMA", "SYS", "EXA_EXTENSIONS", not(emptyOrNullString()), not(emptyOrNullString()))
                .matches());
        assertThat(getSetup().client().listInstances(),
                allOf(hasSize(1), equalTo(List.of(new Instance().id(name).name(name)))));
    }

    @Test
    void createTwoInstances() {
        getSetup().client().install();
        createInstance("vs1");
        createInstance("vs2");

        assertAll(
                () -> getSetup().exasolMetadata()
                        .assertConnection(table()
                                .row("VS1_CONNECTION",
                                        "Created by Extension Manager for " + getName() + " v" + getCurrentVersion()
                                                + " vs1")
                                .row("VS2_CONNECTION",
                                        "Created by Extension Manager for " + getName() + " v" + getCurrentVersion()
                                                + " vs2")
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
    void createInstanceWithSingleQuote() {
        getSetup().client().install();
        createInstance("Quoted'schema");
        assertAll(
                () -> getSetup().exasolMetadata()
                        .assertConnection(table().row("QUOTED'SCHEMA_CONNECTION",
                                "Created by Extension Manager for S3 Virtual Schema v" + getCurrentVersion()
                                        + " Quoted'schema")
                                .matches()),
                () -> getSetup().exasolMetadata().assertVirtualSchema(table()
                        .row("Quoted'schema", "SYS", "EXA_EXTENSIONS", "S3_FILES_ADAPTER", not(emptyOrNullString()))
                        .matches()));
    }

    @Test
    void deleteNonExistingInstance() {
        assertDoesNotThrow(() -> getSetup().client().deleteInstance("no-such-instance"));
    }

    @Test
    void deleteFailsForUnknownVersion() {
        getSetup().client().assertRequestFails(
                () -> getSetup().client().deleteInstance("unknownVersion", "no-such-instance"),
                equalTo("Version 'unknownVersion' not supported, can only use '" + getCurrentVersion() + "'."),
                equalTo(404));
    }

    @Test
    void deleteExistingInstance() {
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
        createInstance(getExtensionId(), getCurrentVersion(), virtualSchemaName);
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
