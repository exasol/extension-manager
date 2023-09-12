import {
    BadRequestError,
    InternalServerError,
    registerExtension
} from "@exasol/extension-manager-interface";

function createExtension() {
    return {
        name: "MyDemoExtension",
        description: "An extension for testing.",
        category: "Demo category",
        installableVersions: [{ name: "0.1.0", latest: true, deprecated: false }],
        bucketFsUploads: $UPLOADS$,
        install(context, version) {
            $INSTALL_EXTENSION$
        },
        addInstance(context, version, params) {
            $ADD_INSTANCE$
        },
        findInstallations(context, metadata) {
            $FIND_INSTALLATIONS$
        },
        findInstances(context, version) {
            $FIND_INSTANCES$
        },
        uninstall(context, version) {
            $UNINSTALL_EXTENSION$
        },
        upgrade(context) {
            $UPGRADE_EXTENSION$
        },
        deleteInstance(context, extensionVersion, instanceId) {
            $DELETE_INSTANCE$
        },
        getInstanceParameters(context, version) {
            $GET_INSTANCE_PARAMETER_DEFINITIONS$
        },
        readInstanceParameterValues(context, extensionVersion, instanceId) {
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
