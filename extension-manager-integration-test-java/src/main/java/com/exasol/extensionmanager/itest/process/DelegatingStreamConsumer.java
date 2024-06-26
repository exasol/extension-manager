package com.exasol.extensionmanager.itest.process;

import static java.util.Arrays.asList;

import java.io.IOException;
import java.util.List;

/**
 * This {@link ProcessStreamConsumer} forwards all events to the given delegates.
 */
public class DelegatingStreamConsumer implements ProcessStreamConsumer {

    private final List<ProcessStreamConsumer> delegates;

    /**
     * Create a new {@link DelegatingStreamConsumer}.
     * 
     * @param delegates delegates to which events should be forwarded
     */
    public DelegatingStreamConsumer(final ProcessStreamConsumer... delegates) {
        this.delegates = asList(delegates);
    }

    @Override
    public void accept(final String line) {
        this.delegates.forEach(delegate -> delegate.accept(line));
    }

    @Override
    public void readFinished() {
        this.delegates.forEach(ProcessStreamConsumer::readFinished);
    }

    @Override
    public void readFailed(final IOException exception) {
        this.delegates.forEach(delegate -> delegate.readFailed(exception));
    }
}