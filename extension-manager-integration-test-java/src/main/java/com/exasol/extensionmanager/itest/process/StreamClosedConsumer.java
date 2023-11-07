package com.exasol.extensionmanager.itest.process;

import java.io.IOException;
import java.time.Duration;
import java.util.concurrent.CountDownLatch;
import java.util.concurrent.TimeUnit;

import com.exasol.errorreporting.ExaError;

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
        } catch (final InterruptedException exception) {
            Thread.currentThread().interrupt();
            throw new IllegalStateException(ExaError.messageBuilder("E-EITFJ-11")
                    .message("Interrupted while waiting for stream to be closed").ticketMitigation().toString(),
                    exception);
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
