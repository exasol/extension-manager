import {
    BadRequestError, Context, ExaMetadata, ExasolExtension,
    Installation,
    Instance, InternalServerError, Parameter, ParameterValues,
    registerExtension
} from "@exasol/extension-manager-interface";

function createExtension(): ExasolExtension {
    return {
        name: "MyDemoExtension",
        description: "An extension for testing.",
        category: "Demo category",
        installableVersions: [{ name: "0.1.0", latest: true, deprecated: false }],
        bucketFsUploads: $UPLOADS$,
        install(context: Context, version: string) {
            $INSTALL_EXTENSION$
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
        uninstall(context: Context, version: string): void {
            $UNINSTALL_EXTENSION$
        },
        deleteInstance(context: Context, extensionVersion: string, instanceId: string): void {
            $DELETE_INSTANCE$
        },
        getInstanceParameters(context: Context, version: string): Parameter[] {
            $GET_INSTANCE_PARAMETER_DEFINITIONS$
        },
        readInstanceParameterValues(context: Context, extensionVersion: string, instanceId: string): ParameterValues {
            return undefined;
        }
    }
}

if (false) {
    // dummy to keep import
    throw new BadRequestError("dummy");
}
if (false) {
    // dummy to keep import
    throw new InternalServerError("dummy");
}

registerExtension(createExtension())
