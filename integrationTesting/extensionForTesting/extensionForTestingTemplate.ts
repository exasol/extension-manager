import {
    Context, ExaMetadata, ExasolExtension,
    Installation,
    Instance, ParameterValues,
    registerExtension
} from "@exasol/extension-manager-interface";

function createExtension(): ExasolExtension {
    return {
        name: "MyDemoExtension",
        description: "An extension for testing.",
        installableVersions: ["0.1.0"],
        bucketFsUploads: $UPLOADS$,
        install(context: Context) {
            $INSTALL$
        },
        addInstance(_context: Context, _installation: Installation, _params: ParameterValues): Instance {
            return undefined;
        },
        findInstallations(_context: Context, metadata: ExaMetadata): Installation[] {
            $FIND_INSTALLATIONS$
        },
        findInstances(_context: Context, _installation: Installation): Instance[] {
            return [];
        },
        uninstall(_context: Context, _installation: Installation): void {
            //empty on purpose
        },
        deleteInstance(_context: Context, _instance: Instance): void {
            //empty on purpose
        },
        readInstanceParameters(_context: Context, _installation: Installation, _instance: Instance): ParameterValues {
            return undefined;
        }
    }
}

registerExtension(createExtension())
