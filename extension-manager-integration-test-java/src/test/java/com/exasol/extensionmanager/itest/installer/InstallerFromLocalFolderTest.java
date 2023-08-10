package com.exasol.extensionmanager.itest.installer;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.*;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.when;

import java.io.IOException;
import java.nio.file.*;
import java.time.Duration;
import java.util.List;
import java.util.Optional;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.api.io.TempDir;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import com.exasol.extensionmanager.itest.ExtensionTestConfig;
import com.exasol.extensionmanager.itest.process.SimpleProcess;

@ExtendWith(MockitoExtension.class)
class InstallerFromLocalFolderTest {

    @Mock
    ExtensionTestConfig configMock;
    private InstallerFromLocalFolder installer;

    @BeforeEach
    void setup() {
        installer = new InstallerFromLocalFolder(configMock);
    }

    @Test
    void extensionManagerBinaryMissing() {
        when(configMock.getLocalExtensionManagerProject()).thenReturn(Optional.of(Paths.get("localExtManager")));
        when(configMock.buildExtensionManager()).thenReturn(false);
        final IllegalStateException exception = assertThrows(IllegalStateException.class, installer::install);
        assertThat(exception.getMessage(), equalTo("E-EMIT-5: Extension manager executable not found at 'localExtManager/extension-manager' after build. This is an internal error that should not happen. Please report it by opening a GitHub issue."));
    }

    @Test
    void extensionManagerBinaryExists(@TempDir final Path tempDir) throws IOException {
        final Path binary = Files.createFile(tempDir.resolve("extension-manager"));
        when(configMock.getLocalExtensionManagerProject()).thenReturn(Optional.of(tempDir));
        when(configMock.buildExtensionManager()).thenReturn(false);
        assertThat(installer.install(), equalTo(binary));
    }

    @Test
    void buildingFails(@TempDir final Path tempDir) throws IOException {
        when(configMock.getLocalExtensionManagerProject()).thenReturn(Optional.of(tempDir));
        when(configMock.buildExtensionManager()).thenReturn(true);
        final IllegalStateException exception = assertThrows(IllegalStateException.class, installer::install);
        assertThat(exception.getMessage(), startsWith("E-EMIT-12: Command 'go build -o extension-manager cmd/main.go' in working dir '" + tempDir + "' failed"));
    }

    @Test
    void buildCurrentExtensionManager() throws IOException {
        final Path extensionManagerDir = Paths.get("..").toAbsolutePath();
        final Path binary = extensionManagerDir.resolve("extension-manager");
        Files.deleteIfExists(binary);
        when(configMock.getLocalExtensionManagerProject()).thenReturn(Optional.of(extensionManagerDir));
        when(configMock.buildExtensionManager()).thenReturn(true);
        assertThat(installer.install(), equalTo(binary));
        assertBinaryWorks(binary);
    }

    private void assertBinaryWorks(final Path binary) {
        assertThat("Binary file " + binary + " exits", Files.exists(binary), is(true));
        final String output = SimpleProcess.start(List.of(binary.toString(), "--help"), Duration.ofSeconds(3));
        assertThat(output, allOf(startsWith("Usage of"), containsString("-extensionRegistryURL string"),
                containsString("-serverAddress string")));
    }
}
