package com.exasol.extensionmanager.itest;

import static java.util.stream.Collectors.toList;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.*;

import java.io.IOException;
import java.io.UncheckedIOException;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpClient.Redirect;
import java.nio.file.Files;
import java.nio.file.Path;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.api.io.TempDir;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import com.exasol.exasoltestsetup.ExasolTestSetup;

@ExtendWith(MockitoExtension.class)
class PreviousVersionTest {

    @Mock
    private ExtensionManagerSetup setupMock;
    @Mock
    private ExasolTestSetup exasolTestSetupMock;
    @Mock
    private PreviousVersionManager previousVersionManagerMock;
    @TempDir
    private Path tempDir;

    PreviousVersion.Builder builder() {
        final HttpClient httpClient = HttpClient.newBuilder().followRedirects(Redirect.NORMAL).build();
        return new PreviousVersion.Builder(setupMock, exasolTestSetupMock, httpClient, previousVersionManagerMock,
                tempDir).adapterFileName("adapter").currentVersion("currentVersion").previousVersion("previousVersion")
                .extensionFileName("extensionFileName").project("project");
    }

    @Test
    void fetchExtension() throws IOException {
        final PreviousVersion testee = builder().build();
        testee.fetchExtension(URI.create(
                "https://extensions-internal.exasol.com/com.exasol/s3-document-files-virtual-schema/2.6.2/s3-vs-extension.js"));
        final Path file = tempDir.resolve(testee.getExtensionId());
        assertAll(() -> assertTrue(Files.exists(file), "file downloaded"),
                () -> assertThat(Files.size(file), equalTo(20875L)));
        testee.close();
        assertFalse(Files.exists(file), "file was deleted during cleanup");
    }

    @Test
    void fetchExtensionFailsForMissingFile() throws IOException {
        final URI uri = URI.create("https://extensions-internal.exasol.com/no-such-file");
        final PreviousVersion testee = builder().build();
        final IllegalStateException exception = assertThrows(IllegalStateException.class,
                () -> testee.fetchExtension(uri));
        assertThat(exception.getMessage(),
                equalTo("E-EMIT-39: Download of '" + uri + "' failed with non-OK status 404"));
        assertThat(Files.list(tempDir).collect(toList()), empty());
    }

    @Test
    void fetchExtensionFails() throws IOException {
        final URI uri = URI.create("https://invalid-url");
        final PreviousVersion testee = builder().build();
        final UncheckedIOException exception = assertThrows(UncheckedIOException.class,
                () -> testee.fetchExtension(uri));
        assertThat(exception.getMessage(), startsWith("E-EMIT-29: Failed to download"));
        assertThat(Files.list(tempDir).collect(toList()), empty());
    }
}
