package com.exasol.extensionmanager.itest;

import static java.util.stream.Collectors.toList;
import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.*;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.ArgumentMatchers.eq;
import static org.mockito.Mockito.lenient;
import static org.mockito.Mockito.verify;

import java.io.*;
import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpClient.Redirect;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.concurrent.TimeoutException;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.api.io.TempDir;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import com.exasol.bucketfs.Bucket;
import com.exasol.bucketfs.BucketAccessException;
import com.exasol.exasoltestsetup.ExasolTestSetup;

@ExtendWith(MockitoExtension.class)
class PreviousVersionManagerTest {
    private static final String BASE_URL = "https://extensions-internal.exasol.com";
    @Mock
    private ExtensionManagerSetup setupMock;
    @Mock
    private ExasolTestSetup exasolTestSetupMock;
    @Mock
    private Bucket bucketMock;
    @TempDir
    private Path tempDir;
    private PreviousVersionManager testee;

    @BeforeEach
    void setup() {
        final HttpClient httpClient = HttpClient.newBuilder().followRedirects(Redirect.NORMAL).build();
        lenient().when(exasolTestSetupMock.getDefaultBucket()).thenReturn(bucketMock);
        testee = new PreviousVersionManager(setupMock, exasolTestSetupMock, httpClient, tempDir);
    }

    @Test
    void newVersion() {
        assertThat(testee.newVersion().adapterFileName("adapter").currentVersion("current").previousVersion("previous")
                .project("project").extensionFileName("extensionFilename").build(), notNullValue());
    }

    @Test
    void fetchExtension() throws IOException {
        final String extensionId = testee.fetchExtension(
                URI.create(BASE_URL + "/com.exasol/s3-document-files-virtual-schema/2.6.2/s3-vs-extension.js"));
        final Path file = tempDir.resolve(extensionId);
        assertAll(() -> assertTrue(Files.exists(file), "file downloaded"),
                () -> assertThat(Files.size(file), equalTo(20875L)));
        verify(setupMock).addFileToCleanupQueue(file);
    }

    @Test
    void fetchExtensionFailsForMissingFile() throws IOException {
        final URI uri = URI.create(BASE_URL + "/no-such-file");
        final IllegalStateException exception = assertThrows(IllegalStateException.class,
                () -> testee.fetchExtension(uri));
        assertThat(exception.getMessage(),
                equalTo("E-EITFJ-39: Download of '" + uri + "' failed with non-OK status 404"));
        assertThat(Files.list(tempDir).collect(toList()), empty());
    }

    @Test
    void fetchExtensionFails() throws IOException {
        final URI url = URI.create("https://invalid-url");
        final UncheckedIOException exception = assertThrows(UncheckedIOException.class,
                () -> testee.fetchExtension(url));
        assertThat(exception.getMessage(), startsWith("E-EITFJ-42: Failed to download '" + url + "' to"));
        assertThat(Files.list(tempDir).collect(toList()), empty());
    }

    @Test
    void prepareBucketFsFile() throws FileNotFoundException, BucketAccessException, TimeoutException {
        testee.prepareBucketFsFile(
                URI.create(BASE_URL + "/com.exasol/s3-document-files-virtual-schema/2.6.2/s3-vs-extension.js"),
                "filename");
        verify(bucketMock).uploadFile(any(), eq("filename"));
    }

    @Test
    void prepareBucketFsFileFailsForMissingFile()
            throws FileNotFoundException, BucketAccessException, TimeoutException {
        final URI url = URI.create(BASE_URL + "/missing-file");
        final IllegalStateException exception = assertThrows(IllegalStateException.class,
                () -> testee.prepareBucketFsFile(url, "filename"));
        assertThat(exception.getMessage(), equalTo(
                "E-EITFJ-39: Download of 'https://extensions-internal.exasol.com/missing-file' failed with non-OK status 404"));
    }
}
