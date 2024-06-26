package com.exasol.extensionmanager.itest;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;

import org.junit.jupiter.api.Test;

import com.exasol.mavenprojectversiongetter.MavenProjectVersionGetter;

class VersionConsistencyTest {

    @Test
    void verifyConsistentVersions() throws IOException {
        assertThat(versionFromPom(), equalTo(versionFromProjectKeeperConfig()));
    }

    private String versionFromProjectKeeperConfig() throws IOException {
        final Path projectKeeperConfig = Path.of("../.project-keeper.yml").toAbsolutePath();
        final String linePrefix = "version: ";
        return Files.readAllLines(projectKeeperConfig).stream() //
                .filter(line -> line.startsWith(linePrefix)) //
                .map(line -> line.replace(linePrefix, "")) //
                .map(String::trim) //
                .findFirst() //
                .orElseThrow(() -> new AssertionError(
                        "Did not find '" + linePrefix + "' entry in file " + projectKeeperConfig));
    }

    private String versionFromPom() {
        return MavenProjectVersionGetter.getProjectRevision(Path.of("../pom.xml"));
    }
}
