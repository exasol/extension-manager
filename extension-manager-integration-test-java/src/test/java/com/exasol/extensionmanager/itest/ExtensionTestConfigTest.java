package com.exasol.extensionmanager.itest;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.assertAll;
import static org.junit.jupiter.api.Assertions.assertThrows;

import java.io.IOException;
import java.io.UncheckedIOException;
import java.nio.file.*;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.io.TempDir;

class ExtensionTestConfigTest {

    @Test
    void readNonExistingDefaultFile() {
        final ExtensionTestConfig config = ExtensionTestConfig.read();
        assertAll(() -> assertThat(config.buildExtension(), is(true)),
                () -> assertThat(config.buildExtensionManager(), is(true)),
                () -> assertThat(config.getExtensionManagerVersion(), is("latest")),
                () -> assertThat(config.getLocalExtensionManagerProject().isPresent(), is(false)));
    }

    @Test
    void loadNonExistingFileFails() {
        final Path missingFile = Path.of("missing-file");
        final UncheckedIOException exception = assertThrows(UncheckedIOException.class,
                () -> ExtensionTestConfig.loadProperties(missingFile));
        assertThat(exception.getMessage(), equalTo("E-EMIT-26: Error reading config file missing-file"));
    }

    @Test
    void loadEmptyFile(@TempDir final Path tempDir) throws IOException {
        final Path tempFile = tempDir.resolve("emptyFile");
        Files.createFile(tempFile);
        assertThat(ExtensionTestConfig.loadProperties(tempFile), notNullValue());
    }

    @Test
    void getConfigFile() {
        assertThat(ExtensionTestConfig.read().getConfigFile(),
                equalTo(Paths.get("extension-test.properties").toAbsolutePath()));
    }
}
