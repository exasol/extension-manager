package com.exasol.extensionmanager.restapi;

import com.exasol.extensionmanager.ExtensionManager;
import com.exasol.extensionmanager.extensionapi.Installation;
import io.javalin.Javalin;
import io.javalin.http.Context;
import lombok.Data;
import org.jetbrains.annotations.NotNull;

import java.util.List;
import java.util.stream.Collectors;

public class RestApi {

    public RestApi() {
        final Javalin app = Javalin.create(config -> config.jsonMapper(new JsonBJsonMapper())).start(7070);
        app.get("/", ctx -> ctx.result("Exasol Extension Store"));
        app.get("/installations", this::handleGetInstallations);
    }

    @NotNull
    private Context handleGetInstallations(final Context ctx) {
        final List<Installation> installations = new ExtensionManager(null).getInstallations();
        final List<InstallationsResponse.Installation> installationResponses = installations.stream().map(inst -> new InstallationsResponse.Installation(inst.getName())).collect(Collectors.toList());
        return ctx.json(new InstallationsResponse(installationResponses));
    }

    @Data
    public static class InstallationsResponse {
        private final List<Installation> installations;

        @Data
        public static class Installation {
            private final String name;
        }
    }
}
