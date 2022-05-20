package com.exasol.extensionmanager.extensionapi;

import lombok.Data;

import java.util.List;

@Data
public class ExaAllScriptsTable {
    private final List<Row> rows;

    @Data
    public static class Row {
        public final String name;
        public final String text;
    }
}
