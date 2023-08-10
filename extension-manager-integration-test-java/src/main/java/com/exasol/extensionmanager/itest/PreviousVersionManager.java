package com.exasol.extensionmanager.itest;

import java.io.*;
import java.net.URI;
import java.net.http.*;
import java.net.http.HttpClient.Redirect;
import java.net.http.HttpResponse.BodyHandlers;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.concurrent.TimeoutException;
import java.util.logging.Logger;

import com.exasol.bucketfs.BucketAccessException;
import com.exasol.errorreporting.ExaError;
import com.exasol.exasoltestsetup.ExasolTestSetup;
import com.exasol.extensionmanager.itest.PreviousVersion.Builder;

/**
 * This class manages previous extension versions. This is useful for testing the upgrade process.
 */
public class PreviousVersionManager {
    private static final Logger LOGGER = Logger.getLogger(PreviousVersionManager.class.getName());
    private final HttpClient httpClient;
    private final Path extensionFolder;
    private final ExtensionManagerSetup setup;
    private final ExasolTestSetup exasolTestSetup;
    private Builder previousVersion;

    PreviousVersionManager(final ExtensionManagerSetup setup, final ExasolTestSetup exasolTestSetup,
            final HttpClient httpClient, final Path extensionFolder) {
        this.setup = setup;
        this.exasolTestSetup = exasolTestSetup;
        this.httpClient = httpClient;
        this.extensionFolder = extensionFolder;
    }

    static PreviousVersionManager create(final ExtensionManagerSetup setup, final ExasolTestSetup exasolTestSetup,
            final Path extensionFolder) {
        final HttpClient httpClient = HttpClient.newBuilder().followRedirects(Redirect.NORMAL).build();
        return new PreviousVersionManager(setup, exasolTestSetup, httpClient, extensionFolder);
    }

    /**
     * Create a new previous version.
     * 
     * @return a new previous version
     */
    public PreviousExtensionVersion.Builder newVersion() {
        return new PreviousExtensionVersion.Builder(setup, this);
    }

    /**
     * Downloads the file from the given URL and uploads it to BucketFS under the specified file name.
     * 
     * @param url      URL to download
     * @param fileName BucketFS file name
     */
    public void prepareBucketFsFile(final URI url, final String fileName) {
        final Path tempFile = createTempFile();
        try {
            downloadToFile(url, tempFile);
            uploadToBucketFs(fileName, tempFile);
        } finally {
            deleteFile(tempFile);
        }
    }

    private void uploadToBucketFs(final String fileName, final Path adapterTempFile) {
        try {
            this.exasolTestSetup.getDefaultBucket().uploadFile(adapterTempFile, fileName);
        } catch (FileNotFoundException | BucketAccessException | TimeoutException exception) {
            throw new IllegalStateException(ExaError.messageBuilder("E-EMIT-35")
                    .message("Failed to upload file {{local file}} to {{bucket file}} in default bucket",
                            adapterTempFile, adapterTempFile)
                    .toString(), exception);
        }
    }

    private void downloadToFile(final URI url, final Path file) {
        final HttpRequest request = HttpRequest.newBuilder(url).GET().build();
        try {
            final HttpResponse<Path> response = httpClient.send(request, BodyHandlers.ofFile(file));
            final long fileSize = Files.size(file);
            LOGGER.fine(() -> "Downloaded " + url + " with response status " + response.statusCode() + " to " + file
                    + " with file size " + fileSize + " bytes");
            if (response.statusCode() / 100 != 2) {
                deleteFile(file);
                throw new IllegalStateException(ExaError.messageBuilder("E-EMIT-39")
                        .message("Download of {{url}} failed with non-OK status {{status code}}", url,
                                response.statusCode())
                        .toString());
            }
        } catch (final IOException exception) {
            deleteFile(file);
            throw new UncheckedIOException(ExaError.messageBuilder("E-EMIT-42")
                    .message("Failed to download {{url}} to {{target file}}", url, file).toString(), exception);
        } catch (final InterruptedException exception) {
            deleteFile(file);
            Thread.currentThread().interrupt();
            throw new IllegalStateException(
                    ExaError.messageBuilder("E-EMIT-32").message("Download of {{url}} was interrupted", url).toString(),
                    exception);
        }
    }

    /**
     * Downloads an additional extension definition (e.g. the previous version of the extension under test).
     * <p>
     * This will allow installing a previous version of the extension and use it for testing the upgrade process. The
     * extension will be deleted during cleanup so that the following test won't be affected.
     * 
     * @param url URL of the extension file to download, e.g. from a GitHub release
     * @return the ID of the downloaded extension
     */
    String fetchExtension(final URI url) {
        final Path extensionFile = getExtensionFile();
        downloadToFile(url, extensionFile);
        setup.addFileToCleanupQueue(extensionFile);
        return extensionFile.getFileName().toString();
    }

    private Path getExtensionFile() {
        try {
            return Files.createTempFile(extensionFolder, "ext-", ".js");
        } catch (final IOException exception) {
            throw new UncheckedIOException(
                    ExaError.messageBuilder("E-EMIT-40")
                            .message("Failed to create a temp file in {{directory}}", extensionFolder).toString(),
                    exception);
        }
    }

    private Path createTempFile() {
        try {
            return Files.createTempFile("adapter-", ".tmp");
        } catch (final IOException exception) {
            throw new UncheckedIOException(
                    ExaError.messageBuilder("E-EMIT-41").message("Failed to create a temp file").toString(), exception);
        }
    }

    static void deleteFile(final Path file) {
        try {
            Files.delete(file);
        } catch (final IOException exception) {
            throw new UncheckedIOException(
                    ExaError.messageBuilder("E-EMIT-34").message("Error deleting file {{file}}", file).toString(),
                    exception);
        }
    }
}
