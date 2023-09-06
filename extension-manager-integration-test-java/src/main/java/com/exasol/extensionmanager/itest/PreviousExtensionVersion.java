package com.exasol.extensionmanager.itest;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.net.URI;
import java.util.Objects;
import java.util.logging.Logger;

import com.exasol.errorreporting.ExaError;
import com.exasol.extensionmanager.client.model.UpgradeExtensionResponse;

/**
 * This represents a previous version of an extension.
 */
public class PreviousExtensionVersion {
    private static final Logger LOGGER = Logger.getLogger(PreviousExtensionVersion.class.getName());
    private final ExtensionManagerSetup setup;
    private final PreviousVersionManager previousVersionManager;
    private final String project;
    private final String extensionFileName;
    private final String currentVersion;
    private final String previousVersion;
    private final String adapterFileName;
    private String tempExtensionId;

    private PreviousExtensionVersion(final Builder builder) {
        this.setup = Objects.requireNonNull(builder.setup);
        this.previousVersionManager = Objects.requireNonNull(builder.previousVersionManager, "previousVersionManager");
        this.project = Objects.requireNonNull(builder.project, "project");
        this.extensionFileName = Objects.requireNonNull(builder.extensionFileName, "extensionFileName");
        this.currentVersion = Objects.requireNonNull(builder.currentVersion, "currentVersion");
        this.previousVersion = Objects.requireNonNull(builder.previousVersion, "previousVersion");
        this.adapterFileName = builder.adapterFileName;
    }

    /**
     * Prepare this previous version by downloading extension JavaScript definition and the adapter file.
     */
    public void prepare() {
        LOGGER.fine(() -> "Preparing extension file '" + extensionFileName + "' in BucketFS...");
        this.tempExtensionId = previousVersionManager.fetchExtension(getDownloadUrl(this.extensionFileName));
        if (this.adapterFileName != null) {
            LOGGER.fine(() -> "Preparing adapter file '" + adapterFileName + "' in BucketFS...");
            previousVersionManager.prepareBucketFsFile(getDownloadUrl(adapterFileName), adapterFileName);
        } else {
            LOGGER.fine("No adapter file name given, skipping adapter.");
        }
    }

    /**
     * Install this version by calling {@link ExtensionManagerClient#install()}.
     */
    public void install() {
        this.setup.client().install(getExtensionId(), previousVersion);
    }

    /**
     * Upgrade the previous version to the current version by calling {@link ExtensionManagerClient#upgrade()} and
     * verify that it returns the expected versions.
     */
    public void upgrade() {
        final UpgradeExtensionResponse upgradeResult = setup.client().upgrade(this.extensionFileName);
        assertEquals(
                new UpgradeExtensionResponse().previousVersion(this.previousVersion).newVersion(this.currentVersion),
                upgradeResult);
    }

    /**
     * The the temporary ID of the installed extension. This ID will only be valid for the currently running test. The
     * extension definition file will be automatically deleted after the test.
     * 
     * @return the temporary ID.
     */
    public String getExtensionId() {
        if (tempExtensionId == null) {
            throw new IllegalStateException(ExaError.messageBuilder("E-EMIT-37")
                    .message("Previous version not prepared.").mitigation("Call method prepare first.").toString());
        }
        return tempExtensionId;
    }

    private URI getDownloadUrl(final String fileName) {
        return URI.create("https://extensions-internal.exasol.com/com.exasol/" + project + "/" + previousVersion + "/"
                + fileName);
    }

    /**
     * Builder for {@link PreviousExtensionVersion} instances.
     */
    public static class Builder {
        private final ExtensionManagerSetup setup;
        private final PreviousVersionManager previousVersionManager;
        private String adapterFileName;
        private String currentVersion;
        private String previousVersion;
        private String extensionFileName;
        private String project;
        private PreviousExtensionVersion builtVersion;

        Builder(final ExtensionManagerSetup setup, final PreviousVersionManager previousVersionManager) {
            this.setup = setup;
            this.previousVersionManager = previousVersionManager;
        }

        /**
         * Set the adapter file name, e.g. {@code document-files-virtual-schema-dist-7.3.3-s3-2.6.2.jar}.
         * <p>
         * This is optional. If the adapter file name is {@code null}, no adapter file will be downloaded.
         * 
         * @param adapterFileName adapter file name
         * @return {@code this} for method chaining
         */
        public Builder adapterFileName(final String adapterFileName) {
            this.adapterFileName = adapterFileName;
            return this;
        }

        /**
         * The version currently under development, e.g. {@code 2.7.0}. This is used to verify that upgrading from the
         * previous version was successful.
         * 
         * @param version the current version
         * @return {@code this} for method chaining
         */
        public Builder currentVersion(final String version) {
            this.currentVersion = version;
            return this;
        }

        /**
         * The version of the adapter that was published previously. E.g. if the currently developed version specified
         * as {@link #currentVersion(String)} is {@code 2.7.0} then this could be {@code 2.6.0} or {@code 2.6.2}. This
         * version is used to build the download URLs for the extension repository, e.g.
         * {@code https://extensions-internal.exasol.com/com.exasol/$PROJECT/$VERSION/$FILENAME}
         * 
         * @param version the previous version
         * @return {@code this} for method chaining
         */
        public Builder previousVersion(final String version) {
            this.previousVersion = version;
            return this;
        }

        /**
         * The file number under which the extension JavaScript file is published, e.g. {@code s3-vs-extension.js}.
         * 
         * @param extensionFileName the extension file name
         * @return {@code this} for method chaining
         */
        public Builder extensionFileName(final String extensionFileName) {
            this.extensionFileName = extensionFileName;
            return this;
        }

        /**
         * The project name under which this project is published, e.g. {@code s3-document-files-virtual-schema}. This
         * name is used to build the download URLs for the extension repository, e.g.
         * {@code https://extensions-internal.exasol.com/com.exasol/$PROJECT/$VERSION/$FILENAME}.
         * 
         * @param project project name in the extension repository
         * @return {@code this} for method chaining
         */
        public Builder project(final String project) {
            this.project = project;
            return this;
        }

        /**
         * Build the instance or return the previously built instance.
         * 
         * @return the prepared previous version.
         */
        public PreviousExtensionVersion build() {
            if (builtVersion != null) {
                throw new IllegalStateException(
                        ExaError.messageBuilder("E-EMIT-36").message("Version was already prepared").toString());
            }
            this.builtVersion = new PreviousExtensionVersion(this);
            return this.builtVersion;
        }
    }
}
