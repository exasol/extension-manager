package com.exasol.extensionmanager.itest;

import static org.hamcrest.MatcherAssert.assertThat;

import java.sql.*;

import org.hamcrest.Matcher;

import com.exasol.matcher.ResultSetStructureMatcher;

/**
 * This class simplifies verifying the content of Exasol's metadata tables, e.g. scripts or virtual schemas.
 */
public class ExasolMetadata {

    private static final String VARCHAR_TYPE = "VARCHAR";
    private final Connection connection;
    private final String extensionSchemaName;

    ExasolMetadata(final Connection connection, final String extensionSchemaName) {
        this.connection = connection;
        this.extensionSchemaName = extensionSchemaName;
    }

    /**
     * Verify the content of the {@code SYS.EXA_ALL_SCRIPTS} table.
     * 
     * @param matcher matcher for verifying the table content
     */
    public void assertScript(final Matcher<ResultSet> matcher) {
        try (final PreparedStatement statement = this.connection.prepareStatement(
                "SELECT SCRIPT_NAME, SCRIPT_TYPE, SCRIPT_INPUT_TYPE, SCRIPT_RESULT_TYPE, SCRIPT_TEXT, SCRIPT_COMMENT "
                        + " FROM SYS.EXA_ALL_SCRIPTS " //
                        + " WHERE SCRIPT_SCHEMA=?" //
                        + " ORDER BY SCRIPT_NAME")) {
            statement.setString(1, this.extensionSchemaName);
            assertThat(statement.executeQuery(), matcher);
        } catch (final SQLException exception) {
            throw new IllegalStateException(exception);
        }
    }

    /**
     * Verify that the {@code SYS.EXA_ALL_SCRIPTS} table is empty.
     */
    public void assertNoScripts() {
        assertScript(ResultSetStructureMatcher
                .table(VARCHAR_TYPE, VARCHAR_TYPE, VARCHAR_TYPE, VARCHAR_TYPE, VARCHAR_TYPE, VARCHAR_TYPE).matches());
    }

    /**
     * Verify the content of the {@code SYS.EXA_ALL_CONNECTIONS} table.
     * 
     * @param matcher matcher for verifying the table content
     */
    public void assertConnection(final Matcher<ResultSet> matcher) {
        assertResult(
                "SELECT CONNECTION_NAME, CONNECTION_COMMENT FROM SYS.EXA_ALL_CONNECTIONS ORDER BY CONNECTION_NAME ASC",
                matcher);
    }

    /**
     * Verify that the {@code SYS.EXA_ALL_CONNECTIONS} table is empty.
     */
    public void assertNoConnections() {
        assertConnection(ResultSetStructureMatcher.table(VARCHAR_TYPE, VARCHAR_TYPE).matches());
    }

    /**
     * Verify the content of the {@code SYS.EXA_ALL_VIRTUAL_SCHEMAS} table.
     * 
     * @param matcher matcher for verifying the table content
     */
    public void assertVirtualSchema(final Matcher<ResultSet> matcher) {
        assertResult(
                "SELECT SCHEMA_NAME, SCHEMA_OWNER, ADAPTER_SCRIPT, ADAPTER_NOTES FROM SYS.EXA_ALL_VIRTUAL_SCHEMAS ORDER BY SCHEMA_NAME ASC",
                matcher);
    }

    /**
     * Verify that the {@code SYS.EXA_ALL_VIRTUAL_SCHEMAS} table is empty.
     */
    public void assertNoVirtualSchema() {
        assertVirtualSchema(
                ResultSetStructureMatcher.table(VARCHAR_TYPE, VARCHAR_TYPE, VARCHAR_TYPE, VARCHAR_TYPE).matches());
    }

    private void assertResult(final String sql, final Matcher<ResultSet> matcher) {
        try (final PreparedStatement statement = this.connection.prepareStatement(sql)) {
            assertThat(statement.executeQuery(), matcher);
        } catch (final SQLException exception) {
            throw new IllegalStateException(exception);
        }
    }
}
