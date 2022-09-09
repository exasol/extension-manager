import {
    BadRequestError, Context, ExaMetadata, ExasolExtension,
    Installation,
    Instance, InternalServerError, ParameterValues,
    registerExtension
} from "@exasol/extension-manager-interface";

function createExtension(): ExasolExtension {
    return {
        name: "MyDemoExtension",
        description: "An extension for testing.",
        installableVersions: ["0.1.0"],
        bucketFsUploads: $UPLOADS$,
        install(context: Context, version: string) {
            $INSTALL$
        },
        addInstance(context: Context, version: string, params: ParameterValues): Instance {
            $ADD_INSTANCE$
        },
        findInstallations(context: Context, metadata: ExaMetadata): Installation[] {
            $FIND_INSTALLATIONS$
        },
        findInstances(context: Context, version: string): Instance[] {
            $FIND_INSTANCES$
        },
        uninstall(context: Context, installation: Installation): void {
            //empty on purpose
        },
        deleteInstance(context: Context, instanceId: string): void {
            $DELETE_INSTANCE$
        },
        readInstanceParameters(context: Context, metadata: ExaMetadata, instanceId: string): ParameterValues {
            return undefined;
        }
    }
}

if(false) {
    // dummy to keep import
    throw new BadRequestError("dummy");
}
if(false) {
    // dummy to keep import
    throw new InternalServerError("dummy");
}

registerExtension(createExtension())
