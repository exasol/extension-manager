package com.exasol.extensionmanager.client.model;

import static org.junit.jupiter.api.Assertions.assertEquals;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.stream.Stream;

import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.MethodSource;

import nl.jqno.equalsverifier.EqualsVerifier;

class EqualsContractTest {
    private static final Path GENERATED_SOURCES = Path
            .of("../extension-manager-client-java/target/generated-sources/swagger/extension-manager-client/");
    private static final String MODEL_PACKAGE = "com.exasol.extensionmanager.client.model";

    @ParameterizedTest
    @MethodSource("findModelClasses")
    void testEqualsContract(final Class<?> modelClass) {
        EqualsVerifier.simple().forClass(modelClass).verify();
    }

    @Test
    void testModelClassesFound() throws IOException {
        assertEquals(15, findModelClasses().count(), "model class count");
    }

    private static Stream<Class<?>> findModelClasses() throws IOException {
        final String modelPath = MODEL_PACKAGE.replace(".", "/");
        final Path sourcePath = GENERATED_SOURCES.toAbsolutePath().resolve(modelPath);
        return Files.list(sourcePath).map(Path::getFileName).map(Path::toString) //
                .filter(fileName -> fileName.endsWith(".java"))
                .map(fileName -> fileName.substring(0, fileName.lastIndexOf(".")))
                .map(className -> MODEL_PACKAGE + "." + className)
                .map(qualifiedClassName -> loadClass(qualifiedClassName));
    }

    private static Class<?> loadClass(final String qualifiedClassName) {
        try {
            return Class.forName(qualifiedClassName);
        } catch (final ClassNotFoundException exception) {
            throw new IllegalStateException("Failed to create class for name '" + qualifiedClassName + "'", exception);
        }
    }
}
