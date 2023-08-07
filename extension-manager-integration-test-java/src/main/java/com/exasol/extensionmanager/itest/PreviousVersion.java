package com.exasol.extensionmanager.itest;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.io.*;
import java.net.URI;
import java.net.http.*;
import java.net.http.HttpResponse.BodyHandlers;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.Objects;
import java.util.concurrent.TimeoutException;
import java.util.logging.Logger;

import com.exasol.bucketfs.BucketAccessException;
import com.exasol.errorreporting.ExaError;
import com.exasol.exasoltestsetup.ExasolTestSetup;
import com.exasol.extensionmanager.client.model.UpgradeExtensionResponse;

/**
 * This represents a previous version of an extension.
 */
public class PreviousVersion {
    private static final Logger LOGGER = Logger.getLogger(PreviousVersion.class.getName());
    private final ExtensionManagerSetup setup;
    private final ExasolTestSetup exasolTestSetup;
    private final HttpClient httpClient;
    private final PreviousVersionManager previousVersionManager;
    private final String project;
    private final String extensionFileName;
    private final String currentVersion;
    private final String previousVersion;
    private final Path extensionFolder;
    private final String adapterFileName;
    private Path extensionFile;

    private PreviousVersion(final Builder builder) {
        this.setup = Objects.requireNonNull(builder.setup);
        this.exasolTestSetup = Objects.requireNonNull(builder.exasolTestSetup);
        this.httpClient = Objects.requireNonNull(builder.httpClient);
        this.extensionFolder = Objects.requireNonNull(builder.extensionFolder);
        this.previousVersionManager = Objects.requireNonNull(builder.previousVersionManager, "previousVersionManager");
        this.project = Objects.requireNonNull(builder.project, "project");
        this.extensionFileName = Objects.requireNonNull(builder.extensionFileName, "extensionFileName");
        this.currentVersion = Objects.requireNonNull(builder.currentVersion, "currentVersion");
        this.previousVersion = Objects.requireNonNull(builder.previousVersion, "previousVersion");
        this.adapterFileName = Objects.requireNonNull(builder.adapterFileName, "adapterFileName");
    }

    /**
     * Prepare this previous version by downloading extension JavaScript definition and the adapter file.
     */
    public void prepare() {
        fetchExtension(getDownloadUrl(this.extensionFileName));
        final Path adapterTempFile = downloadToTemp(getDownloadUrl(adapterFileName));
        try {
            this.exasolTestSetup.getDefaultBucket().uploadFile(adapterTempFile, adapterFileName);
        } catch (FileNotFoundException | BucketAccessException | TimeoutException exception) {
            throw new IllegalStateException(ExaError.messageBuilder("E-EMIT-35")
                    .message("Failed to upload file {{local file}} to {{bucket file}} in default bucket",
                            adapterTempFile, adapterTempFile)
                    .toString(), exception);
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
     * Downloads an additional extension definition (e.g. the previous version of the extension under test).
     * <p>
     * This will allow installing a previous version of the extension and use it for testing the upgrade process.
     * 
     * @param url URL of the extension file to download, e.g. from a GitHub release
     * @return the ID of the downloaded extension
     */
    void fetchExtension(final URI url) {
        final HttpRequest request = HttpRequest.newBuilder(url).GET().build();
        try {
            this.extensionFile = Files.createTempFile(extensionFolder, "ext-", ".js");
            LOGGER.info(() -> "Downloading " + url + " to " + extensionFile + "...");
            final HttpResponse<Path> response = httpClient.send(request, BodyHandlers.ofFile(extensionFile));
            final long fileSize = Files.size(extensionFile);
            if (response.statusCode() / 100 != 2) {
                deleteFile(this.extensionFile);
                this.extensionFile = null;
                throw new IllegalStateException(ExaError.messageBuilder("E-EMIT-39")
                        .message("Download of {{url}} failed with non-OK status {{status code}}", url,
                                response.statusCode())
                        .toString());
            }
            LOGGER.info(() -> "Got response status " + response.statusCode() + ", file size: " + fileSize + " bytes");
        } catch (final IOException exception) {
            deleteFile(this.extensionFile);
            this.extensionFile = null;
            throw new UncheckedIOException(ExaError.messageBuilder("E-EMIT-29")
                    .message("Failed to download {{url}} to local folder {{folder}}", url, extensionFolder).toString(),
                    exception);
        } catch (final InterruptedException exception) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException(
                    ExaError.messageBuilder("E-EMIT-30").message("Download of {{url}} was interrupted", url).toString(),
                    exception);
        }
    }

    private Path downloadToTemp(final URI url) {
        final HttpRequest request = HttpRequest.newBuilder(url).GET().build();
        try {
            final Path tempFile = Files.createTempFile("adapter-", ".tmp");
            final HttpResponse<Path> response = httpClient.send(request, BodyHandlers.ofFile(tempFile));
            final long fileSize = Files.size(tempFile);
            LOGGER.fine(() -> "Downloaded " + url + " with response status " + response.statusCode() + " to " + tempFile
                    + " with file size " + fileSize + " bytes");
            return tempFile;
        } catch (final IOException exception) {
            throw new UncheckedIOException(ExaError.messageBuilder("E-EMIT-31")
                    .message("Failed to download {{url}} to temp file", url).toString(), exception);
        } catch (final InterruptedException exception) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException(
                    ExaError.messageBuilder("E-EMIT-32").message("Download of {{url}} was interrupted", url).toString(),
                    exception);
        }
    }

    /**
     * The the temporary ID of the installed extension. This ID will only be valid until {@link #close()} is called.
     * 
     * @return the temporary ID.
     */
    public String getExtensionId() {
        if (extensionFile == null) {
            throw new IllegalStateException(
                    ExaError.messageBuilder("E-EMIT-37").message("Previous version not prepared").toString());
        }
        return this.extensionFile.getFileName().toString();
    }

    private URI getDownloadUrl(final String fileName) {
        return URI.create("https://extensions-internal.exasol.com/com.exasol/" + project + "/" + previousVersion + "/"
                + fileName);
    }

    public static class Builder {

        private final ExtensionManagerSetup setup;
        private final PreviousVersionManager previousVersionManager;
        private final Path extensionFolder;
        private final HttpClient httpClient;
        private final ExasolTestSetup exasolTestSetup;
        private String adapterFileName;
        private String currentVersion;
        private String previousVersion;
        private String extensionFileName;
        private String project;
        private PreviousVersion builtVersion;

        public Builder(final ExtensionManagerSetup setup, final ExasolTestSetup exasolTestSetup,
                final HttpClient httpClient, final PreviousVersionManager previousVersionManager,
                final Path extensionFolder) {
            this.setup = setup;
            this.exasolTestSetup = exasolTestSetup;
            this.httpClient = httpClient;
            this.previousVersionManager = previousVersionManager;
            this.extensionFolder = extensionFolder;
        }

        /**
         * Set the adapter file name, e.g. {@code document-files-virtual-schema-dist-7.3.3-s3-2.6.2.jar}.
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
         * @return
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
        public PreviousVersion build() {
            if (builtVersion != null) {
                throw new IllegalStateException(
                        ExaError.messageBuilder("E-EMIT-36").message("Version was already prepared").toString());
            }
            this.builtVersion = new PreviousVersion(this);
            return this.builtVersion;
        }

        /**
         * Called by {@link PreviousVersionManager#cleanup()}.
         */
        void close() {
            if (builtVersion != null) {
                builtVersion.close();
            }
        }
    }

    /**
     * Called by {@link Builder#close()}.
     */
    void close() {
        deleteFile(this.extensionFile);
    }

    private static void deleteFile(final Path file) {
        try {
            Files.delete(file);
        } catch (final IOException exception) {
            throw new UncheckedIOException(
                    ExaError.messageBuilder("E-EMIT-34").message("Error deleting file {{file}}", file).toString(),
                    exception);
        }
    }
}
