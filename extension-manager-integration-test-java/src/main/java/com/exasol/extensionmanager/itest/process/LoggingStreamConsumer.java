package com.exasol.extensionmanager.itest.process;

import java.io.IOException;
import java.util.logging.Level;
import java.util.logging.Logger;

/**
 * This {@link ProcessStreamConsumer} logs all content with a configurable log level using a {@link Logger}.
 */
public class LoggingStreamConsumer implements ProcessStreamConsumer {
    private static final Logger LOGGER = Logger.getLogger(LoggingStreamConsumer.class.getName());
    private final String prefix;
    private final Level logLevel;

    /**
     * Create a new {@link LoggingStreamConsumer}.
     * 
     * @param prefix   prefix for all log messages
     * @param logLevel log level for all log messages
     */
    public LoggingStreamConsumer(final String prefix, final Level logLevel) {
        this.prefix = prefix;
        this.logLevel = logLevel;
    }

    @Override
    public void accept(final String line) {
        LOGGER.log(this.logLevel, () -> this.prefix + line);
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