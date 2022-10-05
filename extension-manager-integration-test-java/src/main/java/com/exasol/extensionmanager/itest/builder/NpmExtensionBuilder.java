package com.exasol.extensionmanager.itest.builder;

import java.nio.file.Files;
import java.nio.file.Path;
import java.time.Duration;
import java.util.List;

import com.exasol.extensionmanager.itest.process.SimpleProcess;

class NpmExtensionBuilder implements ExtensionBuilder {

    private final Path sourceDir;
    private final Path builtJsExtension;

    NpmExtensionBuilder(final Path sourceDir, final Path builtJsExtension) {
        this.sourceDir = sourceDir;
        this.builtJsExtension = builtJsExtension;
    }

    @Override
    public void build() {
        if (!Files.exists(sourceDir)) {
            throw new IllegalArgumentException("Extension source dir " + sourceDir + " does not exist");
        }
        SimpleProcess.start(sourceDir, List.of("npm", "run", "build"), Duration.ofSeconds(30));
    }

    @Override
    public Path getExtensionFile() {
        return builtJsExtension;
    }
}
