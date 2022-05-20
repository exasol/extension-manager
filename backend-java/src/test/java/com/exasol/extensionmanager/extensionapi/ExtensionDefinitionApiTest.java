package com.exasol.extensionmanager.extensionapi;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;
import static org.mockito.Mockito.mock;
import static org.mockito.Mockito.verify;

import java.nio.file.Path;
import java.util.List;

import org.junit.jupiter.api.Test;

class ExtensionDefinitionApiTest {
    // TODO build extensionForTesting automatically. Here or in maven?
    private static final Path PATH_OF_DEMO_EXTENSION = Path.of("../backend/extensionApi/extensionForTesting/dist.js");

    @Test
    void testLoadExtension() {
        try (final ExtensionDefinition extensionDef = new ExtensionDefinitionApi(PATH_OF_DEMO_EXTENSION)
                .loadExtension()) {
            assertThat(extensionDef.getName(), equalTo("MyDemoExtension"));
        }
    }

    @Test
    void testInstall() {
        final ExtensionDefinitionApi api = new ExtensionDefinitionApi(PATH_OF_DEMO_EXTENSION);
        for (int i = 0; i < 100; i++) {
            final SimpleSqlClient sqlClient = mock(SimpleSqlClient.class);
            final long start = System.currentTimeMillis();
            try (final ExtensionDefinition extensionDef = api.loadExtension()) {
                extensionDef.install(sqlClient);
                System.out.println(System.currentTimeMillis() - start);
                verify(sqlClient).runQuery("CREATE ADAPTER SCRIPT ...");
            }
        }
    }

    @Test
    void testFindInstallations() {
        try (final ExtensionDefinition extensionDef = new ExtensionDefinitionApi(PATH_OF_DEMO_EXTENSION)
                .loadExtension()) {
            final SimpleSqlClient sqlClient = mock(SimpleSqlClient.class);
            mock(SimpleSqlClient.class);
            final ExaAllScriptsTable allScriptsTable = new ExaAllScriptsTable(
                    List.of(new ExaAllScriptsTable.Row("test.my-adapter", "some SQL")));
            final List<Installation> installations = extensionDef.findInstallations(sqlClient, allScriptsTable);
            assertThat(installations.get(0).getName(), equalTo("test.my-adapter"));
        }
    }
}