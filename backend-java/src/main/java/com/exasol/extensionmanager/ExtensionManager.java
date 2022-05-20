package com.exasol.extensionmanager;

import com.exasol.extensionmanager.extensionapi.*;

import java.sql.Connection;
import java.util.List;
import java.util.stream.Collectors;

public class ExtensionManager {
    private final List<ExtensionDefinition> extensionDefinitions;
    private final Connection sqlConnection;

    public ExtensionManager(final Connection sqlConnection) {
        this.sqlConnection = sqlConnection;
        this.extensionDefinitions = new ExtensionDefinitionProvider().getExtensionDefinitions();
    }

    public List<Installation> getInstallations() {
        final ExaAllScriptsTable allScriptsTable = new ExaAllScriptsTable(List.of(new ExaAllScriptsTable.Row("test.my-adapter", "some SQL")));
        return this.extensionDefinitions.stream().flatMap(ext -> ext.findInstallations(new ExasolSimpleSqlClient(this.sqlConnection), allScriptsTable).stream()).collect(Collectors.toList());
    }
}
