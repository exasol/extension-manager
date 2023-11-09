package com.exasol.extensionmanager.itest.base;

import static java.util.Collections.emptyList;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.assertAll;
import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assumptions.assumeTrue;

import java.util.List;
import java.util.logging.Logger;

import org.junit.jupiter.api.*;

import com.exasol.extensionmanager.client.model.ExtensionsResponseExtension;
import com.exasol.extensionmanager.client.model.InstallationsResponseInstallation;
import com.exasol.extensionmanager.itest.ExtensionManagerSetup;
import com.exasol.extensionmanager.itest.PreviousExtensionVersion;

/**
 * This is a integration tests base class for {@code SCRIPT} based Extensions that already contains some basic tests for
 * installing/listing/uninstalling extensions.
 */
public abstract class AbstractScriptExtensionIT {

    private static final Logger LOG = Logger.getLogger(AbstractScriptExtensionIT.class.getName());
    private final ExtensionITConfig config;

    /**
     * Create a new base integration test.
     */
    protected AbstractScriptExtensionIT() {
        this.config = createConfig();
    }

    /**
     * Get the {@link ExtensionManagerSetup extension manager setup}.
     * 
     * @return extension manager setup
     */
    protected abstract ExtensionManagerSetup getSetup();

    /**
     * Create a new configuration for the integration tests.
     * 
     * @return new configuration
     */
    protected abstract ExtensionITConfig createConfig();

    /**
     * Assert that the installed {@code SCRIPT}s work as expected.
     */
    protected abstract void assertScriptsWork();

    /**
     * Assert that the expected {@code SCRIPT}s exist after installing the extension.
     */
    protected abstract void assertScriptsExist();

    @BeforeEach
    void logTestName(final TestInfo testInfo) {
        LOG.info(">>> " + testInfo.getDisplayName());
    }

    @AfterEach
    void cleanup() {
        getSetup().cleanup();
    }

    /**
     * Verify that extension is listed with expected properties.
     */
    @Test
    void listExtensions() {
        final List<ExtensionsResponseExtension> extensions = getSetup().client().getExtensions();
        assertAll(() -> assertThat(extensions, hasSize(1)), //
                () -> assertThat(extensions.get(0).getName(), equalTo(config.getExtensionName())),
                () -> assertThat(extensions.get(0).getInstallableVersions().get(0).getName(),
                        equalTo(config.getCurrentVersion())),
                () -> assertThat(extensions.get(0).getInstallableVersions().get(0).isLatest(), is(true)),
                () -> assertThat(extensions.get(0).getInstallableVersions().get(0).isDeprecated(), is(false)),
                () -> assertThat(extensions.get(0).getDescription(), equalTo(config.getExtensionDescription())));
    }

    /**
     * Verify that listing installations returns an empty list when there is no installation.
     */
    @Test
    void getInstallationsReturnsEmptyList() {
        assertThat(getSetup().client().getInstallations(), hasSize(0));
    }

    /**
     * Verify that listing installations finds {@code SCRIPT}s created by the extension.
     */
    @Test
    void getInstallationsReturnsResult() {
        getSetup().client().install();
        assertThat(getSetup().client().getInstallations(), contains(new InstallationsResponseInstallation() //
                .id(config.getExtensionId()) //
                .name(config.getExtensionName()) //
                .version(config.getCurrentVersion())));
    }

    /**
     * Verify that installing an unknown version fails.
     */
    @Test
    void installingWrongVersionFails() {
        getSetup().client().assertRequestFails(() -> getSetup().client().install("wrongVersion"),
                equalTo("Version 'wrongVersion' not supported, can only use '" + config.getCurrentVersion() + "'."),
                equalTo(404));
        getSetup().exasolMetadata().assertNoScripts();
    }

    /**
     * Verify that installing the extension creates expected {@code SCRIPT}s.
     */
    @Test
    void installCreatesScripts() {
        getSetup().client().install();
        assertScriptsExist();
    }

    /**
     * Verify that installing the extension twice creates expected {@code SCRIPT}s.
     */
    @Test
    void installingTwiceCreatesScripts() {
        getSetup().client().install();
        getSetup().client().install();
        assertScriptsExist();
    }

    /**
     * Verify that installed {@code SCRIPT}s work as expected.
     */
    @Test
    void installedScriptsWork() {
        getSetup().client().install();
        assertScriptsWork();
    }

    /**
     * Verify that uninstalling an extension that is not yet install does not fail.
     */
    @Test
    void uninstallExtensionWithoutInstallation() {
        assertDoesNotThrow(() -> getSetup().client().uninstall());
    }

    /**
     * Verify that uninstalling the extension removes all {@code SCRIPT}s.
     */
    @Test
    void uninstallExtensionRemovesScripts() {
        getSetup().client().install();
        assertScriptsExist();
        getSetup().client().uninstall();
        getSetup().exasolMetadata().assertNoScripts();
    }

    /**
     * Verify that uninstalling an unknown version fails.
     */
    @Test
    void uninstallWrongVersionFails() {
        getSetup().client().assertRequestFails(() -> getSetup().client().uninstall("wrongVersion"),
                equalTo("Version 'wrongVersion' not supported, can only use '" + config.getCurrentVersion() + "'."),
                equalTo(404));
    }

    /**
     * Verify that listing instances is not supported.
     */
    @Test
    void listingInstancesNotSupported() {
        getSetup().client().assertRequestFails(() -> getSetup().client().listInstances(),
                equalTo("Finding instances not supported"), equalTo(404));
    }

    /**
     * Verify that creating instances is not supported.
     */
    @Test
    void creatingInstancesNotSupported() {
        getSetup().client().assertRequestFails(() -> getSetup().client().createInstance(emptyList()),
                equalTo("Creating instances not supported"), equalTo(404));
    }

    /**
     * Verify that deleting instances is not supported.
     */
    @Test
    void deletingInstancesNotSupported() {
        getSetup().client().assertRequestFails(() -> getSetup().client().deleteInstance("inst"),
                equalTo("Deleting instances not supported"), equalTo(404));
    }

    /**
     * Verify that getting extension details is not supported.
     */
    @Test
    void getExtensionDetailsInstancesNotSupported() {
        getSetup().client().assertRequestFails(
                () -> getSetup().client().getExtensionDetails(config.getCurrentVersion()),
                equalTo("Creating instances not supported"), equalTo(404));
    }

    /**
     * Verify that upgrading the extension fails when it is not installed.
     */
    @Test
    public void upgradeFailsWhenNotInstalled() {
        getSetup().client().assertRequestFails(() -> getSetup().client().upgrade(), //
                allOf(startsWith("Not all required scripts are installed: Validation failed: Script"),
                        endsWith("is missing")),
                equalTo(412));
    }

    /**
     * Verify that upgrading fails when the latest version is already installed.
     */
    @Test
    public void upgradeFailsWhenAlreadyUpToDate() {
        getSetup().client().install();
        getSetup().client().assertRequestFails(() -> getSetup().client().upgrade(),
                "Extension is already installed in latest version " + config.getCurrentVersion(), 412);
    }

    /**
     * Verify that upgrading from the previous version works and the scripts continue working.
     */
    @Test
    void upgradeFromPreviousVersion() {
        assumeTrue(config.getPreviousVersion() != null, "No previous version available for testing");
        final PreviousExtensionVersion previousVersion = createPreviousVersion();
        previousVersion.prepare();
        previousVersion.install();
        assertScriptsWork();
        assertInstalledVersion(config.getPreviousVersion(), previousVersion);
        previousVersion.upgrade();
        assertInstalledVersion(config.getCurrentVersion(), previousVersion);
        assertScriptsWork();
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

    private PreviousExtensionVersion createPreviousVersion() {
        return getSetup().previousVersionManager().newVersion().currentVersion(config.getCurrentVersion()) //
                .previousVersion(config.getPreviousVersion()) //
                .adapterFileName(config.getPreviousVersionJarFile()) //
                .extensionFileName(config.getExtensionId()) //
                .project("cloud-storage-extension") //
                .build();
    }
}
