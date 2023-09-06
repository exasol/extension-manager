package com.exasol.extensionmanager.itest;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.junit.jupiter.api.Assertions.assertDoesNotThrow;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.ArgumentMatchers.any;
import static org.mockito.Mockito.*;

import java.net.URI;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import com.exasol.extensionmanager.client.model.UpgradeExtensionResponse;

@ExtendWith(MockitoExtension.class)
class PreviousExtensionVersionTest {

    @Mock
    private ExtensionManagerSetup setupMock;
    @Mock
    private PreviousVersionManager previousVersionManagerMock;
    @Mock
    private ExtensionManagerClient clientMock;

    @BeforeEach
    private void setup() {
        lenient().when(setupMock.client()).thenReturn(clientMock);
    }

    private PreviousExtensionVersion.Builder builder() {
        return new PreviousExtensionVersion.Builder(setupMock, previousVersionManagerMock).adapterFileName("adapter")
                .currentVersion("currentVersion").previousVersion("previousVersion")
                .extensionFileName("extensionFileName").project("project");
    }

    @Test
    void prepare() {
        when(previousVersionManagerMock.fetchExtension(URI.create("https://extensions-internal.exasol.com/com.exasol/project/previousVersion/extensionFileName"))).thenReturn("extId");
        final PreviousExtensionVersion testee = builder().build();
        testee.prepare();
        verify(previousVersionManagerMock).prepareBucketFsFile(URI.create("https://extensions-internal.exasol.com/com.exasol/project/previousVersion/adapter"), "adapter");
        assertThat(testee.getExtensionId(),equalTo("extId"));
    }

    @Test
    void prepareWithoutAdapter() {
        when(previousVersionManagerMock.fetchExtension(URI.create("https://extensions-internal.exasol.com/com.exasol/project/previousVersion/extensionFileName"))).thenReturn("extId");
        final PreviousExtensionVersion testee = builder().adapterFileName(null).build();
        testee.prepare();
        verify(previousVersionManagerMock, never()).prepareBucketFsFile(any(), any());
        assertThat(testee.getExtensionId(),equalTo("extId"));
    }

    @Test
    void getExtensionIdFailsWhenNotPrepared() {
        final PreviousExtensionVersion testee = builder().build();
        final IllegalStateException exception = assertThrows(IllegalStateException.class, testee::getExtensionId);
        assertThat(exception.getMessage(),
                equalTo("E-EMIT-37: Previous version not prepared. Call method prepare first."));
    }

    @Test
    void upgrade() {
        when(clientMock.upgrade("extensionFileName")).thenReturn(new UpgradeExtensionResponse().previousVersion("previousVersion").newVersion("currentVersion"));
        final PreviousExtensionVersion testee = builder().build();
        assertDoesNotThrow(testee::upgrade);
    }
}
