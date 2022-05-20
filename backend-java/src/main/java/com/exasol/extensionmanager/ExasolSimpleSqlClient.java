package com.exasol.extensionmanager;

import com.exasol.extensionmanager.extensionapi.SimpleSqlClient;

import java.sql.*;

public class ExasolSimpleSqlClient implements SimpleSqlClient {

    private Connection sqlConnection;

    public ExasolSimpleSqlClient(Connection sqlConnection) {

        this.sqlConnection = sqlConnection;
    }

    @Override
    public void runQuery(String query) {
        try(final Statement statement = sqlConnection.createStatement()){
            statement.executeUpdate(query);
        } catch (SQLException e) {
            throw new RuntimeException(e);//todo
        }
    }
}
