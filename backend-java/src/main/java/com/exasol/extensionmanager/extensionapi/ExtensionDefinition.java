package com.exasol.extensionmanager.extensionapi;

import java.util.List;

public interface ExtensionDefinition extends AutoCloseable {
    void install(SimpleSqlClient sqlClient);

    List<Installation> findInstallations(SimpleSqlClient sqlClient, ExaAllScriptsTable exaAllScriptsTable);

    @Override
    void close();

    String getName();
}
