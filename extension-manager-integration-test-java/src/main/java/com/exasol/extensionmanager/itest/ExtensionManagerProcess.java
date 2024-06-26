package com.exasol.extensionmanager.itest;

import java.io.IOException;
import java.net.ServerSocket;
import java.nio.file.Path;
import java.time.Duration;
import java.util.ArrayList;
import java.util.List;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;
import java.util.logging.Level;
import java.util.logging.Logger;
import java.util.regex.Pattern;

import com.exasol.errorreporting.ExaError;
import com.exasol.extensionmanager.itest.process.*;

/**
 * This class allows starting and stopping an extension manager process.
 */
// [impl -> dsn~eitfj-start-extension-manager~1]
class ExtensionManagerProcess implements AutoCloseable {
    private static final Logger LOGGER = Logger.getLogger(ExtensionManagerProcess.class.getName());
    private static final Duration SERVER_STARTUP_TIMEOUT = Duration.ofSeconds(5);
    private final SimpleProcess process;
    private final int port;

    private ExtensionManagerProcess(final SimpleProcess process, final int port) {
        this.process = process;
        this.port = port;
    }

    static ExtensionManagerProcess start(final Path extensionManagerBinary, final Path extensionFolder) {
        final int port = findOpenPort();
        LOGGER.info(() -> "Starting extension manager " + extensionManagerBinary + " on port " + port
                + " with extension folder " + extensionFolder + "...");
        final List<String> command = new ArrayList<>(List.of(extensionManagerBinary.toString(), "-extensionRegistryURL",
                extensionFolder.toString(), "-serverAddress", "localhost:" + port));
        addFlagIfSupported(extensionManagerBinary, "-addCauseToInternalServerError", command);

        final ServerStartupConsumer serverPortConsumer = new ServerStartupConsumer();
        final SimpleProcess process = SimpleProcess.start(command,
                new DelegatingStreamConsumer(new LoggingStreamConsumer("server stdout>", Level.FINE)),
                new DelegatingStreamConsumer(new LoggingStreamConsumer("server stderr>", Level.FINE),
                        serverPortConsumer));
        if (!serverPortConsumer.isStartupFinished(SERVER_STARTUP_TIMEOUT)) {
            process.stop();
            throw new IllegalStateException(ExaError.messageBuilder("E-EITFJ-17")
                    .message("Extension manager did not log server port after {{timeout}}.", SERVER_STARTUP_TIMEOUT)
                    .mitigation("Check log output for error messages.").toString());
        }
        return new ExtensionManagerProcess(process, port);
    }

    private static void addFlagIfSupported(final Path extensionManagerBinary, final String flag,
            final List<String> command) {
        if (supportsFlag(extensionManagerBinary, flag)) {
            command.add(flag);
        }
    }

    static boolean supportsFlag(final Path extensionManagerBinary, final String flag) {
        final String helpContent = getHelpContent(extensionManagerBinary);
        return helpContent.contains(flag);
    }

    static String getHelpContent(final Path extensionManagerBinary) {
        return SimpleProcess.start(List.of(extensionManagerBinary.toString(), "-help"), Duration.ofSeconds(3));
    }

    private static int findOpenPort() {
        try (ServerSocket socket = new ServerSocket(0)) {
            return socket.getLocalPort();
        } catch (final IOException exception) {
            throw new IllegalStateException(ExaError.messageBuilder("E-EITFJ-18")
                    .message("Failed to find an open port: {{error message}}", exception.getMessage()).toString(),
                    exception);
        }
    }

    @Override
    public void close() {
        LOGGER.fine("Stopping extension manager process");
        this.process.stop();
    }

    String getServerBasePath() {
        return "http://localhost:" + this.port;
    }

    private static class ServerStartupConsumer implements ProcessStreamConsumer {
        @SuppressWarnings("java:S5852") // Accepting potential denial of service as this class will be used only for
                                        // tests
        private static final Pattern STARTUP_FINISHED = Pattern.compile(".*Starting server on localhost:\\d+.*");
        private final CountDownLatch startupFinishedLatch = new CountDownLatch(1);

        @Override
        public void accept(final String line) {
            if (STARTUP_FINISHED.matcher(line).matches()) {
                this.startupFinishedLatch.countDown();
                LOGGER.info(() -> "Found server startup message in line '" + line + "'");
            }
        }

        boolean isStartupFinished(final Duration timeout) {
            return awaitStartupFinished(timeout);
        }

        private boolean awaitStartupFinished(final Duration timeout) {
            try {
                return this.startupFinishedLatch.await(timeout.toMillis(), TimeUnit.MILLISECONDS);
            } catch (final InterruptedException exception) {
                Thread.currentThread().interrupt();
                throw new IllegalStateException(
                        ExaError.messageBuilder("E-EITFJ-19")
                                .message("Interrupted while waiting for server startup to finish").toString(),
                        exception);
            }
        }

        @Override
        public void readFinished() {
            // Ignore
        }

        @Override
        public void readFailed(final IOException exception) {
            // Ignore
        }
    }
}
