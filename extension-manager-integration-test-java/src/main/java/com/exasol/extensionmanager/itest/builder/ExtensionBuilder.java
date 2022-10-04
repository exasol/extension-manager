package com.exasol.extensionmanager.itest.builder;

import java.nio.file.Path;

/**
 * This interface allows customizing how extensions are built before running integration tests.
 */
public interface ExtensionBuilder {

    /**
     * Create a default builder that builds the extension by executing {@code npm run build}.
     * 
     * @param sourceDir        source directory containing {@code package.json} and other project files
     * @param builtJsExtension path to the built JS extension file
     * @return new default extension builder
     */
    public static ExtensionBuilder createDefaultNpmBuilder(final Path sourceDir, final Path builtJsExtension) {
        return new NpmExtensionBuilder(sourceDir, builtJsExtension);
    }

    /**
     * Build the extension JS file.
     */
    void build();

    /**
     * Get the path to the extension JS file. The file may be missing or outdated if {@link #build()} was not called
     * before.
     * 
     * @return path to the extension JS file.
     */
    Path getExtensionFile();
}
