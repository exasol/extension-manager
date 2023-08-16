import {
    Context, ExaMetadata,
    ExasolExtension,
    Installation,
    Instance, Parameter, ParameterValues,
    UpgradeResult,
    registerExtension
} from "@exasol/extension-manager-interface";


export function createExtension(): ExasolExtension {
    return {
        name: "Testing Extension",
        category: "testing",
        description: "Extension for testing EM integration test setup",
        installableVersions: [{ name: "0.0.0", latest: true, deprecated: false }],
        bucketFsUploads: [],
        install(context: Context, version: string) {
            console.log(`Install version ${version}`)
            context.sqlClient.execute("select 1")
        },
        addInstance(context: Context, version: string, params: ParameterValues): Instance {
            console.log(`Add instance for version ${version}`)
            return { id: "new-instance", name: "New instance" };
        },
        findInstallations(_context: Context, metadata: ExaMetadata): Installation[] {
            console.log(`Find installations`)
            return [{ name: "Testing Extension", version: "0.0.0" }];
        },
        findInstances(context: Context, version: string): Instance[] {
            console.log(`Find instances of version ${version}`)
            return [{ id: "instance-1", name: "Instance 1" }];
        },
        uninstall(context: Context, version: string): void {
            console.log(`Uninstall version ${version}`)
            context.sqlClient.execute("select 1")
        },
        upgrade(context): UpgradeResult {
            const result = context.sqlClient.query("select '0.2.0'")
            const newVersion: string = result.rows[0][0]
            console.log(`Upgrading to version ${newVersion}`)
            return { previousVersion: "0.1.0", newVersion }
        },
        deleteInstance(context: Context, version: string, instanceId: string): void {
            console.log(`Delete instance ${instanceId} of version ${version}`)
            context.sqlClient.execute("select 1")
        },
        getInstanceParameters(context: Context, version: string): Parameter[] {
            console.log(`Get instance parameters for version ${version}`)
            return [{ id: "param1", name: "Param 1", type: "string", required: true }]
        },
        readInstanceParameterValues(_context: Context, version: string, instanceId: string): ParameterValues {
            console.log(`Read instance parameters for instance ${instanceId} in version ${version}`)
            return { values: [] };
        }
    }
}

registerExtension(createExtension())
