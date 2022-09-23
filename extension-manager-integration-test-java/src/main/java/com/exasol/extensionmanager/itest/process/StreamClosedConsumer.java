package com.exasol.extensionmanager.itest.process;

import java.io.IOException;
import java.time.Duration;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

class StreamClosedConsumer implements ProcessStreamConsumer {

    private final CountDownLatch latch = new CountDownLatch(1);

    void waitUntilStreamClosed(final Duration timeout) {
        if (!await(timeout)) {
            throw new IllegalStateException("Stream was not closed within timeout of " + timeout);
        }
    }

    private boolean await(final Duration timeout) {
        try {
            return this.latch.await(timeout.toMillis(), TimeUnit.MILLISECONDS);
        } catch (final InterruptedException e) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException("Interrupted while waiting for stream to be closed");
        }
    }

    @Override
    public void accept(final String line) {
        // ignore
    }

    @Override
    public void readFinished() {
        this.latch.countDown();
    }

    @Override
    public void readFailed(final IOException ioException) {
        this.latch.countDown();
    }
}
