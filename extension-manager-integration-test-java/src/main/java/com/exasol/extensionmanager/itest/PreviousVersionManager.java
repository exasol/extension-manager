package com.exasol.extensionmanager.itest;

import java.net.http.HttpClient;
import java.net.http.HttpClient.Redirect;
import java.nio.file.Path;
import java.util.logging.Logger;

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

    public PreviousVersion.Builder create() {
        if (previousVersion != null) {
            throw new IllegalStateException(
                    ExaError.messageBuilder("E-EMIT-38").message("Previous version already prepared.")
                            .mitigation("Call this only once per test case").toString());
        }
        this.previousVersion = new PreviousVersion.Builder(setup, exasolTestSetup, httpClient, this,
                this.extensionFolder);
        return this.previousVersion;
    }

    void cleanup() {
        if (previousVersion != null) {
            previousVersion.close();
            previousVersion = null;
        }
    }

    void close() {
    }
}
