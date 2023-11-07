package com.exasol.extensionmanager.itest.builder;

import java.nio.file.Files;
import java.nio.file.Path;
import java.time.Duration;
import java.util.List;

import com.exasol.errorreporting.ExaError;
import com.exasol.extensionmanager.itest.process.SimpleProcess;

class NpmExtensionBuilder implements ExtensionBuilder {

    private final Path sourceDir;
    private final Path builtJsExtension;

    NpmExtensionBuilder(final Path sourceDir, final Path builtJsExtension) {
        verifySourceDirExists(sourceDir);
        this.sourceDir = sourceDir;
        this.builtJsExtension = builtJsExtension;
    }

    private static void verifySourceDirExists(final Path sourceDir) {
        if (!Files.exists(sourceDir)) {
            throw new IllegalArgumentException(ExaError.messageBuilder("E-EITFJ-2")
                    .message("Extension source dir {{source directory}} does not exist", sourceDir).toString());
        }
    }

    @Override
    public void build() {
        SimpleProcess.start(sourceDir, List.of("npm", "run", "build"), Duration.ofSeconds(30));
    }

    @Override
    public Path getExtensionFile() {
        return builtJsExtension;
    }
}
