package com.exasol.extensionmanager.restapi;

import com.exasol.errorreporting.ExaError;
import io.javalin.plugin.json.JsonMapper;
import jakarta.json.bind.Jsonb;
import jakarta.json.bind.JsonbBuilder;
import jakarta.json.bind.spi.JsonbProvider;
import org.jetbrains.annotations.NotNull;

import java.io.*;
import java.nio.charset.StandardCharsets;

class JsonBJsonMapper implements JsonMapper {
    private static final JsonbProvider JSONB_PROVIDER = JsonbProvider.provider();

    @NotNull
    @Override
    public String toJsonString(@NotNull final Object obj) {
        try (final Jsonb jsonb = JsonbBuilder.create()) {
            return jsonb.toJson(obj);
        } catch (final Exception exception) {
            throw new IllegalStateException(ExaError.messageBuilder("F-EMB-1").message("Failed to JSON serialize response.").ticketMitigation().toString());
        }
    }

    @NotNull
    @Override
    public InputStream toJsonStream(@NotNull final Object obj) {
        return new ByteArrayInputStream(toJsonString(obj).getBytes(StandardCharsets.UTF_8));
    }

    @NotNull
    @Override
    public <T> T fromJsonString(@NotNull final String json, @NotNull final Class<T> targetClass) {
        try (final Jsonb jsonb = JSONB_PROVIDER.create().build()) {
            return jsonb.fromJson(json, targetClass);
        } catch (final Exception exception) {
            throw new IllegalStateException(ExaError.messageBuilder("F-EMB-2").message("Failed to JSON deserialize message.").ticketMitigation().toString());
        }
    }

    @NotNull
    @Override
    public <T> T fromJsonStream(@NotNull final InputStream json, @NotNull final Class<T> targetClass) {
        final String jsonString;
        try {
            jsonString = new String(json.readAllBytes(), StandardCharsets.UTF_8);
            return this.fromJsonString(jsonString, targetClass);
        } catch (final IOException e) {
            throw new IllegalStateException(ExaError.messageBuilder("F-EMB-3").message("Failed to read JSON message.").ticketMitigation().toString());
        }
    }
}
