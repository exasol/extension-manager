package com.exasol.extensionmanager;

import java.nio.file.Path;

import com.exasol.extensionmanager.extensionapi.*;

public class Main {

    public static void main(final String[] args) {
        // new RestApi();

        try (final ExtensionDefinition extensionDef = new ExtensionDefinitionApi(
                Path.of("../backend/extensionApi/extensionForTesting/dist.js")).loadExtension()) {
            System.out.println(extensionDef.getName());
            extensionDef.install(new SimpleSqlClient() {
                @Override
                public void runQuery(final String query) {

                }
            });
        }
        /*
         * try { final Connection connection = DriverManager.getConnection(args[0], args[1], args[2]); final ResultSet
         * resultSet = connection.createStatement().executeQuery("SELECT NOW()"); resultSet.next();
         * System.out.println(resultSet.getString(1)); } catch (SQLException e) { throw new RuntimeException(e); }
         */
    }
}
