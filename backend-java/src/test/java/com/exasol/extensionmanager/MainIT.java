package com.exasol.extensionmanager;

import com.exasol.containers.ExasolContainer;
import org.junit.jupiter.api.Test;
import org.testcontainers.junit.jupiter.Container;
import org.testcontainers.junit.jupiter.Testcontainers;

import java.io.IOException;
import java.nio.charset.StandardCharsets;

import static org.hamcrest.MatcherAssert.assertThat;
import static org.hamcrest.Matchers.equalTo;

@Testcontainers
class MainIT {
    @Container
    private static final ExasolContainer<? extends ExasolContainer<?>> EXASOL = new ExasolContainer<>().withReuse(true);

    @Test
    void test234(){
        Main.main(new String[]{EXASOL.getJdbcUrl(), EXASOL.getUsername(), EXASOL.getPassword()});
    }

    @Test
    void testNativeImage() throws InterruptedException, IOException {
        final Process command = Runtime.getRuntime().exec(new String[]{"target/extension-manager-backend", EXASOL.getJdbcUrl(), EXASOL.getUsername(), EXASOL.getPassword()});
        final int exitCode = command.waitFor();
        final String outout = new String(command.getInputStream().readAllBytes(), StandardCharsets.UTF_8);
        final String err = new String(command.getErrorStream().readAllBytes(), StandardCharsets.UTF_8);
        System.out.println(outout);
        System.err.println(err);
        assertThat(exitCode, equalTo(0));
    }
}