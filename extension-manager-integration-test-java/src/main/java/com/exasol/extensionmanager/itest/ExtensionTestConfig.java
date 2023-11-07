package com.exasol.extensionmanager.itest;

import java.io.*;
import java.nio.file.*;
import java.util.Optional;
import java.util.Properties;
import java.util.logging.Logger;

import com.exasol.errorreporting.ExaError;

/**
 * Configuration for integration tests of extensions.
 */
public class ExtensionTestConfig {
    private static final Logger LOGGER = Logger.getLogger(ExtensionTestConfig.class.getName());
    private static final Path CONFIG_FILE = Paths.get("extension-test.properties").toAbsolutePath();
    private final Properties properties;

    private ExtensionTestConfig(final Properties properties) {
        this.properties = properties;
    }

    /**
     * Read the configuration file from the default location. If the file does not exist this returns the default
     * configuration.
     * 
     * @return configuration read from the config file
     */
    static ExtensionTestConfig read() {
        final Path file = CONFIG_FILE;
        if (!Files.exists(file)) {
            LOGGER.info(() -> "Extension test config file " + file + " not found. Using defaults.");
            return new ExtensionTestConfig(new Properties());
        }
        return new ExtensionTestConfig(loadProperties(file));
    }

    static Properties loadProperties(final Path configFile) {
        LOGGER.info(() -> "Reading config file " + configFile);
        try (InputStream stream = Files.newInputStream(configFile)) {
            final Properties props = new Properties();
            props.load(stream);
            return props;
        } catch (final IOException e) {
            throw new UncheckedIOException(ExaError.messageBuilder("E-EITFJ-26")
                    .message("Error reading config file {{config file path}}", configFile).toString(), e);
        }
    }

    /**
     * Get the configured path to the local extension manager project or an empty {@link Optional} if it is not
     * configured.
     * 
     * @return configured path to the local extension manager
     */
    public Optional<Path> getLocalExtensionManagerProject() {
        return getOptionalValue("localExtensionManager") //
                .map(path -> Paths.get(path).toAbsolutePath()) //
                .map(path -> {
                    if (!Files.exists(path) || !Files.isDirectory(path)) {
                        throw new IllegalStateException(
                                ExaError.messageBuilder("E-EITFJ-27")
                                        .message("Invalid extension manager {{path}} in config file {{config file}}.",
                                                path, getConfigFile())
                                        .mitigation("Specify a valid directory.").toString());
                    }
                    return path;
                });
    }

    /**
     * Get the extension manager version to use for the tests. If the POM file of the extension does not specify a
     * version, this defaults to version of the Manifest in the JAR or {@code "latest"} if no Manifest exists.
     * 
     * @return extension manager version
     */
    public String getExtensionManagerVersion() {
        return getOptionalValue("extensionManagerVersion") //
                .or(this::getJarFileVersion).map(version -> "v" + version) //
                .orElse("latest");
    }

    private Optional<String> getJarFileVersion() {
        final String version = this.getClass().getPackage().getImplementationVersion();
        LOGGER.finest(() -> "Found version " + version + " for JAR");
        return Optional.ofNullable(version);
    }

    /**
     * Check if the extension should be built before running the tests. This is useful for speeding up tests when there
     * are no changes to the extension.
     * 
     * @return {@code true} if the extension should be built before running the tests
     */
    public boolean buildExtension() {
        return getOptionalValue("buildExtension").map(Boolean::valueOf).orElse(true);
    }

    /**
     * Check if the extension manager should be built before running the tests. This is useful for speeding up tests
     * when there are no changes to the extension manager.
     * 
     * @return {@code true} if the extension manager should be built before running the tests
     */
    public boolean buildExtensionManager() {
        return getOptionalValue("buildExtensionManager").map(Boolean::valueOf).orElse(true);
    }

    private Optional<String> getOptionalValue(final String param) {
        return Optional.ofNullable(this.properties.getProperty(param));
    }

    /**
     * Get the path of the config file.
     * 
     * @return path of the config file
     */
    public Path getConfigFile() {
        return CONFIG_FILE;
    }
}
