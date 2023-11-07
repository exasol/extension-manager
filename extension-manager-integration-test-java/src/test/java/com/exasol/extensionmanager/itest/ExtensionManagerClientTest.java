package com.exasol.extensionmanager.itest;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.hamcrest.Matchers.startsWith;
import static org.junit.jupiter.api.Assertions.assertAll;
import static org.junit.jupiter.api.Assertions.assertThrows;
import static org.mockito.Mockito.when;

import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import com.exasol.exasoltestsetup.SqlConnectionInfo;
import com.exasol.extensionmanager.client.api.*;
import com.exasol.extensionmanager.client.model.*;
import com.exasol.extensionmanager.itest.ExtensionManagerClient.Extension;

import jakarta.json.JsonObject;

// [utest -> dsn~eitfj-access-extension-manager-rest-interface~1]
@ExtendWith(MockitoExtension.class)
class ExtensionManagerClientTest {

    private static final String DB_HOST = "host";
    private static final int DB_PORT = 1234;
    private static final String DB_USER = "user";
    private static final String DB_PASSWORD = "pass";
    @Mock
    ExtensionApi extensionClientMock;
    @Mock
    InstallationApi installationApiMock;
    @Mock
    InstanceApi instanceClientMock;

    SqlConnectionInfo dbConnectionInfoMock;
    private ExtensionManagerClient testee;

    @BeforeEach
    void setup() {
        dbConnectionInfoMock = new SqlConnectionInfo(DB_HOST, DB_PORT, DB_USER, DB_PASSWORD);
        testee = new ExtensionManagerClient(extensionClientMock, installationApiMock, instanceClientMock,
                dbConnectionInfoMock);
    }

    @Test
    void parseErrorMessageJson() {
        final JsonObject json = testee.parseErrorMessageJson("{\"field\":42}");
        assertThat(json.getInt("field"), equalTo(42));
    }

    @Test
    void parseErrorMessageJsonFails() {
        final IllegalArgumentException exception = assertThrows(IllegalArgumentException.class,
                () -> testee.parseErrorMessageJson("invalid json"));
        assertThat(exception.getMessage(), equalTo("E-EITFJ-15: Failed to parse error message 'invalid json' as JSON"));
    }

    @Test
    void getSingleExtensionNoExtensionFound() {
        when(extensionClientMock.listAvailableExtensions(DB_HOST, DB_PORT)).thenReturn(new ExtensionsResponse());
        final IllegalStateException exception = assertThrows(IllegalStateException.class,  testee::getSingleExtension);
        assertThat(exception.getMessage(), equalTo("E-EITFJ-28: Expected exactly one extension but found 0: []. Check the extension manager log for errors loading the extension."));
    }

    @Test
    void getSingleExtensionMultipleExtensionFound() {
        when(extensionClientMock.listAvailableExtensions(DB_HOST, DB_PORT)).thenReturn(new ExtensionsResponse().addExtensionsItem(new ExtensionsResponseExtension()).addExtensionsItem(new ExtensionsResponseExtension()));
        final IllegalStateException exception = assertThrows(IllegalStateException.class,  testee::getSingleExtension);
        assertThat(exception.getMessage(), startsWith("E-EITFJ-28: Expected exactly one extension but found 2"));
    }

    @Test
    void getSingleExtension() {
        when(extensionClientMock.listAvailableExtensions(DB_HOST, DB_PORT)).thenReturn(new ExtensionsResponse().addExtensionsItem(new ExtensionsResponseExtension().id("id")));
          final ExtensionsResponseExtension ext = testee.getSingleExtension();
        assertThat(ext.getId(), equalTo("id"));
    }

    @Test
    void getExtension() {
        when(extensionClientMock.listAvailableExtensions(DB_HOST, DB_PORT)).thenReturn(new ExtensionsResponse().addExtensionsItem(new ExtensionsResponseExtension().id("id").addInstallableVersionsItem(new ExtensionVersion().name("ver"))));
        final Extension ext = testee.getExtension();
        assertAll(() -> assertThat(ext.getId(), equalTo("id")),
            () -> assertThat(ext.getCurrentVersion(), equalTo("ver")));
    }
}
