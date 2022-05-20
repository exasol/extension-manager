package com.exasol.extensionmanager.extensionapi;

import java.io.IOException;
import java.nio.file.Files;
import java.nio.file.Path;
import java.util.List;
import java.util.stream.Collectors;

import javax.script.*;

import org.jetbrains.annotations.NotNull;
import org.openjdk.nashorn.api.scripting.JSObject;

import com.exasol.errorreporting.ExaError;

import jakarta.json.bind.Jsonb;
import jakarta.json.bind.spi.JsonbProvider;
import lombok.*;

public class ExtensionDefinitionApi {
    private static final JsonbProvider JSONB_PROVIDER = JsonbProvider.provider();
    private final ScriptEngine jsEngine;
    private final CompiledScript compiledScript;

    public ExtensionDefinitionApi(final Path extensionPath) {
        final ScriptEngineManager scriptEngineManager = new ScriptEngineManager();
        this.jsEngine = scriptEngineManager.getEngineByName("nashorn");
        try {
            this.compiledScript = ((Compilable) this.jsEngine).compile(Files.readString(extensionPath));
        } catch (final ScriptException | IOException e) {
            throw new RuntimeException(e);
        }
    }

    public ExtensionDefinition loadExtension() {
        try {
            final SimpleScriptContext scriptContext = new SimpleScriptContext();
            final Bindings scopedBindings = this.jsEngine.createBindings();
            scriptContext.setBindings(scopedBindings, ScriptContext.ENGINE_SCOPE);
            scopedBindings.put("installedExtension", "");
            this.compiledScript.eval(scriptContext);

            final JSObject installedExtension = (JSObject) scopedBindings.get("installedExtension");
            final JSObject extension = (JSObject) installedExtension.getMember("extension");
            final String extensionName = (String) extension.getMember("name");
            return new ExtensionDefinitionImpl(scopedBindings, extension, extensionName);
        } catch (final RuntimeException | ScriptException e) {
            throw new RuntimeException(e);
        }
    }

    /**
     * This function converts an object into an JS object in the VM. This approach does only work for data objects with
     * not callback functions.
     * <p>
     * In contrast to directly passing the Java object this approach has the advantage that it creates native JS
     * objects, while graal.js just creates objects that imitate the interface. The problem there is that the interface
     * is incomplete and does for example not support calling map on lists.
     * </p>
     *
     * @param object  object to pass to the JS VM
     * @param context JS VM context
     * @return
     */
    private static JSObject toJsDataObj(final Object object, final Bindings context) {
        final JSObject json = (JSObject) context.get("JSON");
        final JSObject parseFunction = (JSObject) json.getMember("parse");
        return (JSObject) parseFunction.call(context, toJson(object));
    }

    private static String toJson(final Object object) {
        try (final Jsonb jsonb = JSONB_PROVIDER.create().build()) {
            return jsonb.toJson(object);
        } catch (final Exception exception) {
            throw new IllegalStateException(ExaError.messageBuilder("F-EMB-4")
                    .message("Failed to serialize data object.").ticketMitigation().toString());
        }
    }

    @RequiredArgsConstructor
    private static class ExtensionDefinitionImpl implements ExtensionDefinition {
        private final Bindings context;
        private final JSObject extension;

        @Getter
        private final String name;

        @Override
        public void install(final SimpleSqlClient sqlClient) {
            try {
                final JSObject installFunction = (JSObject) this.extension.getMember("install");
                installFunction.call(this.extension, sqlClient);
            } catch (final ClassCastException e) {
                throw new RuntimeException(e);
            }
        }

        public List<Installation> findInstallations(final SimpleSqlClient sqlClient,
                final ExaAllScriptsTable exaAllScriptsTable) {
            final JSObject findInstallationsFunction = (JSObject) this.extension.getMember("findInstallations");
            final JSObject result = (JSObject) findInstallationsFunction.call(this.extension, sqlClient,
                    toJsDataObj(exaAllScriptsTable, this.context));
            return result.values().stream().map(this::readInstallation).collect(Collectors.toList());
        }

        @NotNull
        private InstallationImpl readInstallation(final Object installationJs) {
            final JSObject installationJsObj = (JSObject) installationJs;
            return new InstallationImpl(installationJsObj.getMember("name").toString());
        }

        @Override
        public void close() {

        }
    }

    @Data
    private static class InstallationImpl implements Installation {
        private final String name;

        private InstallationImpl(final String name) {
            this.name = name;
        }
    }
}
