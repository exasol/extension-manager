package com.exasol.extensionmanager.itest;

import java.io.*;
import java.net.URI;
import java.net.http.*;
import java.net.http.HttpResponse.BodyHandlers;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.concurrent.TimeoutException;
import java.util.logging.Logger;

import com.exasol.bucketfs.BucketAccessException;
import com.exasol.errorreporting.ExaError;
import com.exasol.exasoltestsetup.ExasolTestSetup;

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
    private final String version;
    private final Path extensionFolder;
    private final String adapterFileName;
    private Path extensionFile;

    private PreviousVersion(final Builder builder) {
        this.setup = builder.setup;
        this.exasolTestSetup = builder.exasolTestSetup;
        this.httpClient = builder.httpClient;
        this.extensionFolder = builder.extensionFolder;
        this.previousVersionManager = builder.previousVersionManager;
        this.project = builder.project;
        this.extensionFileName = builder.extensionFileName;
        this.version = builder.version;
        this.adapterFileName = builder.adapterFileName;
    }

    /**
     * Prepare the previous version by downloading extension JavaScript definition and the adapter file.
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
            LOGGER.info(() -> "Got response status " + response.statusCode() + ", file size: " + fileSize + " bytes");
        } catch (final IOException exception) {
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

    String getExtensionId() {
        if (extensionFile == null) {
            throw new IllegalStateException(
                    ExaError.messageBuilder("E-EMIT-37").message("Previous version not prepared").toString());
        }
        return this.extensionFile.getFileName().toString();
    }

    private URI getDownloadUrl(final String fileName) {
        return URI.create(
                "https://extensions-internal.exasol.com/com.exasol/" + project + "/" + version + "/" + fileName);
    }

    /**
     * A builder for {@link PreviousVersion} instances.
     */
    public static class Builder {
        private final ExtensionManagerSetup setup;
        private final PreviousVersionManager previousVersionManager;
        private final Path extensionFolder;
        private final HttpClient httpClient;
        private final ExasolTestSetup exasolTestSetup;
        private String adapterFileName;
        private String version;
        private String extensionFileName;
        private String project;
        private PreviousVersion builtVersion;

        Builder(final ExtensionManagerSetup setup, final ExasolTestSetup exasolTestSetup, final HttpClient httpClient,
                final PreviousVersionManager previousVersionManager, final Path extensionFolder) {
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
         * The version of the adapter that was published previously. E.g. if the currently developed version is
         * {@code 2.7.0} then this should be {@code 2.6.0} or {@code 2.6.2}. This version is used to build the download
         * URLs for the extension repository, e.g.
         * {@code https://extensions-internal.exasol.com/com.exasol/$PROJECT/$VERSION/$FILENAME}
         * 
         * @param version the previous version
         * @return {@code this} for method chaining
         */
        public Builder version(final String version) {
            this.version = version;
            return this;
        }

        /**
         * The file number under which the extension JavaScript file is published, e.g. {@code s3-vs-extension.js}.
         * 
         * @param extensionFileName the extension file name
         * @return {@code this} for method chaining
         */
        public Builder setExtensionFileName(final String extensionFileName) {
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
        public Builder setProject(final String project) {
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
        try {
            Files.delete(this.extensionFile);
        } catch (final IOException exception) {
            throw new UncheckedIOException(ExaError.messageBuilder("E-EMIT-34")
                    .message("Error deleting file {{file}}", this.extensionFile).toString(), exception);
        }
    }
}
