package com.exasol.extensionmanager;

import java.nio.file.Path;
import java.util.List;

import com.exasol.extensionmanager.extensionapi.ExtensionDefinition;
import com.exasol.extensionmanager.extensionapi.ExtensionDefinitionApi;

public class ExtensionDefinitionProvider {
    public List<ExtensionDefinition> getExtensionDefinitions() {
        return List.of(new ExtensionDefinitionApi(Path.of("../backend/extensionApi/extensionForTesting/dist.js"))
                .loadExtension());
    }
}
