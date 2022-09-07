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
        findInstances(context: Context, installation: Installation): Instance[] {
            return [];
        },
        uninstall(context: Context, installation: Installation): void {
            //empty on purpose
        },
        deleteInstance(context: Context, instance: Instance): void {
            //empty on purpose
        },
        readInstanceParameters(context: Context, installation: Installation, instance: Instance): ParameterValues {
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
