package com.exasol.extensionmanager.itest.process;

import java.io.IOException;
import java.nio.file.Path;
import java.time.Duration;
import java.time.Instant;
import java.util.List;
import java.util.concurrent.TimeUnit;
import java.util.logging.Level;
import java.util.logging.Logger;

import com.exasol.errorreporting.ExaError;

/**
 * This is a convenient wrapper for {@link ProcessBuilder} and {@link Process} that simplifies starting a process,
 * waiting for it to finish and getting its stdout.
 */
public class SimpleProcess {
    private static final Logger LOGGER = Logger.getLogger(SimpleProcess.class.getName());

    private final Process process;
    private final Path workingDirectory;
    private final List<String> command;
    private final Instant startTime;

    private SimpleProcess(final Process process, final Path workingDirectory, final List<String> command,
            final Instant startTime) {
        this.process = process;
        this.workingDirectory = workingDirectory;
        this.command = command;
        this.startTime = startTime;
    }

    /**
     * Starts a new process using the working directory of the current Java process and waits until it terminates
     * successfully.
     *
     * @param command          command to execute
     * @param executionTimeout execution timeout for the process
     * @return combined stdout and stderr from the process
     */
    public static String start(final List<String> command, final Duration executionTimeout) {
        return start(null, command, executionTimeout);
    }

    /**
     * Starts a new process using the given working directory and waits until it terminates successfully.
     *
     * @param workingDirectory directory in which to start the process. Use the working directory of the current Java
     *                         process if {@code null}.
     * @param command          command to execute
     * @param executionTimeout execution timeout for the process
     * @return combined stdout and stderr from the process
     */
    public static String start(final Path workingDirectory, final List<String> command,
            final Duration executionTimeout) {
        final StringBuilder stringBuilder = new StringBuilder();
        final CollectingStreamConsumer collectingStreamConsumer = new CollectingStreamConsumer(stringBuilder);
        final StreamClosedConsumer stdoutStreamClosed = new StreamClosedConsumer();
        final StreamClosedConsumer stderrStreamClosed = new StreamClosedConsumer();
        final SimpleProcess process = start(workingDirectory, command,
                new DelegatingStreamConsumer(collectingStreamConsumer, new LoggingStreamConsumer("stdout>", Level.FINE),
                        stdoutStreamClosed),
                new DelegatingStreamConsumer(collectingStreamConsumer, new LoggingStreamConsumer("stderr>", Level.FINE),
                        stderrStreamClosed));
        process.waitUntilTerminatedSuccessfully(executionTimeout);
        stdoutStreamClosed.waitUntilStreamClosed(Duration.ofMillis(100));
        stderrStreamClosed.waitUntilStreamClosed(Duration.ofMillis(100));
        return stringBuilder.toString();
    }

    /**
     * Starts a new process using the working directory of the current Java process.
     *
     * @param command              command to execute
     * @param outputStreamConsumer {@link ProcessStreamConsumer} for the output stream
     * @param errorStreamConsumer  {@link ProcessStreamConsumer} for the error stream
     * @return a new {@link SimpleProcess} you can use to wait for the process to finish and retrieve its output
     */
    public static SimpleProcess start(final List<String> command, final ProcessStreamConsumer outputStreamConsumer,
            final ProcessStreamConsumer errorStreamConsumer) {
        return start(null, command, outputStreamConsumer, errorStreamConsumer);
    }

    /**
     * Starts a new process.
     *
     * @param workingDirectory     directory in which to start the process. Use the working directory of the current
     *                             Java process if {@code null}
     * @param command              command to execute
     * @param outputStreamConsumer {@link ProcessStreamConsumer} for the output stream
     * @param errorStreamConsumer  {@link ProcessStreamConsumer} for the error stream
     * @return a new {@link SimpleProcess} you can use to wait for the process to finish and retrieve its output
     */
    public static SimpleProcess start(final Path workingDirectory, final List<String> command,
            final ProcessStreamConsumer outputStreamConsumer, final ProcessStreamConsumer errorStreamConsumer) {
        LOGGER.fine(() -> "Executing command '" + formatCommand(command) + "' in working dir " + workingDirectory);
        try {
            final Process process = new ProcessBuilder(command)
                    .directory(workingDirectory == null ? null : workingDirectory.toFile()) //
                    .redirectErrorStream(false) //
                    .start();
            final Instant startTime = Instant.now();
            new AsyncStreamReader().startCollectingConsumer(process.getInputStream(), outputStreamConsumer);
            new AsyncStreamReader().startCollectingConsumer(process.getErrorStream(), errorStreamConsumer);
            return new SimpleProcess(process, workingDirectory, command, startTime);
        } catch (final IOException exception) {
            throw new IllegalStateException(ExaError.messageBuilder("E-EITFJ-8")
                    .message("Error executing command {{command}}.", String.join(" ", command))
                    .mitigation("Verify that the executable {{executable}} is on the PATH.", command.get(0)).toString(),
                    exception);
        }
    }

    /**
     * Wait for the process to terminate successfully.
     *
     * @param executionTimeout maximum time to wait until the process finishes
     * @throws IllegalStateException if the process did not finish within the given timeout or returned an exit code
     *                               other than 0
     */
    public void waitUntilTerminatedSuccessfully(final Duration executionTimeout) {
        waitForProcessTerminated(executionTimeout);
        final Duration duration = Duration.between(this.startTime, Instant.now());
        final int exitCode = this.process.exitValue();
        if (exitCode != 0) {
            throw new IllegalStateException(ExaError.messageBuilder("E-EITFJ-12").message(
                    "Command {{command}} in working dir {{working dir}} failed with exit code {{exit code|u}} after {{duration|u}}.",
                    formatCommand(), workingDirectory, exitCode, duration).mitigation("See log output for details.")
                    .toString());
        }
        LOGGER.fine(() -> "Command '" + formatCommand() + "' finished successfully after " + duration);
    }

    /**
     * Destroy a running process.
     */
    public void stop() {
        this.process.destroy();
    }

    private void waitForProcessTerminated(final Duration executionTimeout) {
        try {
            if (!this.process.waitFor(executionTimeout.toMillis(), TimeUnit.MILLISECONDS)) {
                throw new IllegalStateException(ExaError.messageBuilder("E-EITFJ-9")
                        .message("Timeout while waiting {{timeout duration|u}} for command {{command}}",
                                executionTimeout, formatCommand())
                        .toString());
            }
        } catch (final InterruptedException exception) {
            throw handleInterruptedException(exception);
        }
    }

    private RuntimeException handleInterruptedException(final InterruptedException exception) {
        Thread.currentThread().interrupt();
        return new IllegalStateException(ExaError.messageBuilder("E-EITFJ-10")
                .message("Interrupted while waiting for command {{command}} to terminate.", formatCommand())
                .ticketMitigation().toString(), exception);
    }

    private static String formatCommand(final List<String> command) {
        return String.join(" ", command);
    }

    private String formatCommand() {
        return formatCommand(this.command);
    }
}
