package com.exasol.extensionmanager.itest.installer;

import java.nio.file.*;
import java.time.Duration;
import java.util.List;
import java.util.logging.Logger;

import com.exasol.extensionmanager.itest.ExtensionTestConfig;
import com.exasol.extensionmanager.itest.OsCheck;
import com.exasol.extensionmanager.itest.process.SimpleProcess;

class InstallerFromGitHub implements ExtensionManagerInstaller {
    private static final Logger LOGGER = Logger.getLogger(InstallerFromGitHub.class.getName());
    private final ExtensionTestConfig config;

    InstallerFromGitHub(final ExtensionTestConfig config) {
        this.config = config;
    }

    @Override
    public Path install() {
        if (this.config.buildExtensionManager()) {
            runGoInstall();
        } else {
            LOGGER.warning("Skipping installation of extension manager");
        }
        return getExtensionManagerExecutable();
    }

    private void runGoInstall() {
        final String version = this.config.getExtensionManagerVersion();
        LOGGER.info(() -> "Installing extension manager version '" + version + "'...");
        SimpleProcess.start(List.of("go", "install", "github.com/exasol/extension-manager/cmd@" + version),
                Duration.ofSeconds(60));
    }

    private Path getExtensionManagerExecutable() {
        final String executableName = "cmd" + OsCheck.getExecutableSuffix();
        final Path executablePath = getGoPath().resolve("bin").resolve(executableName);
        if (!Files.exists(executablePath)) {
            throw new IllegalStateException("Executable was not installed at '" + executablePath + "'");
        }
        return executablePath;
    }

    private Path getGoPath() {
        final String rawPath = SimpleProcess.start(List.of("go", "env", "GOPATH"), Duration.ofSeconds(1));
        final Path goPath = Paths.get(rawPath.trim());
        if (!Files.exists(goPath)) {
            throw new IllegalStateException("GOPATH does not exist: '" + goPath + "'");
        }
        LOGGER.info(() -> "Got GOPATH '" + goPath + "'");
        return goPath;
    }
}