package com.exasol.extensionmanager.itest;

import java.nio.file.Path;

import org.junit.jupiter.api.extension.ExtendWith;
import org.junit.jupiter.api.io.TempDir;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import com.exasol.exasoltestsetup.ExasolTestSetup;

@ExtendWith(MockitoExtension.class)
class PreviousExtensionVersionTest {

    @Mock
    private ExtensionManagerSetup setupMock;
    @Mock
    private ExasolTestSetup exasolTestSetupMock;
    @Mock
    private PreviousVersionManager previousVersionManagerMock;
    @TempDir
    private Path tempDir;

    PreviousExtensionVersion.Builder builder() {

        return new PreviousExtensionVersion.Builder(setupMock, previousVersionManagerMock).adapterFileName("adapter")
                .currentVersion("currentVersion").previousVersion("previousVersion")
                .extensionFileName("extensionFileName").project("project");
    }

}
