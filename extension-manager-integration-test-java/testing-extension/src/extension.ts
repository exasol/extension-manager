import {
    Context, ExaMetadata,
    ExasolExtension,
    Installation,
    Instance, Parameter, ParameterValues,
    registerExtension
} from "@exasol/extension-manager-interface";


export function createExtension(): ExasolExtension {
    return {
        name: "Testing Extension",
        description: "Extension for testing EM integration test setup",
        installableVersions: [{ name: "0.0.0", latest: true, deprecated: false }],
        bucketFsUploads: [],
        install(context: Context, version: string) {
            context.sqlClient.execute("select 1")
        },
        addInstance(context: Context, version: string, params: ParameterValues): Instance {
            return { id: "new-instance", name: "New instance" };
        },
        findInstallations(_context: Context, metadata: ExaMetadata): Installation[] {
            return [{ name: "Testing Extension", version: "0.0.0" }];
        },
        findInstances(context: Context, version: string): Instance[] {
            return [{ id: "instance-1", name: "Instance 1" }];
        },
        uninstall(context: Context, version: string): void {
            context.sqlClient.execute("select 1")
        },
        deleteInstance(context: Context, version: string, instanceId: string): void {
            context.sqlClient.execute("select 1")
        },
        getInstanceParameters(context: Context, version: string): Parameter[] {
            return [{ id: "param1", name: "Param 1", type: "string", required: true }]
        },
        readInstanceParameterValues(_context: Context, _version: string, _instanceId: string): ParameterValues {
            return { values: [] };
        }
    }
}

registerExtension(createExtension())
