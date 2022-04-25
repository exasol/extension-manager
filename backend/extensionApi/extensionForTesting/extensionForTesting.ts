import {
    ExasolExtension,
    Installation,
    Instance,
    Operators,
    ParameterValues,
    registerExtension,
    SqlClient
} from "../../../interface/src/api";


function createExtension(): ExasolExtension {
    return {
        name: "MyDemoExtension",
        description: "An extension for testing.",
        bucketFsUploads: [{
            name: "My Extension JAR",
            downloadUrl: "https://my.download.de/demo.jar",
            licenseUrl: "https://my.download.de/LICENSE",
            bucketFsFilename: "my-extension.1.2.3.jar", // do we need this?,  maybe we can autogenerate it my extension-id + version  + artifactid
            licenseAgreementRequired: false
        }
        ],
        instanceParameters: [
            {
                id: "direction",
                name: "Direction",
                type: "select",
                options: {
                    import: "Import",
                    export: "Export"
                }
            },
            {
                id: "",
                name: "",
                type: "string",
                regex: /\d+/,
                condition: {
                    and: [
                        {
                            parameter: "direction",
                            operator: Operators.EQ,
                            value: "Import"
                        },
                        {
                            parameter: "amount",
                            operator: Operators.LESS,
                            value: 2
                        },
                    ]
                }
            },
            {
                id: "",
                name: "",
                type: "string"
            }],
        install(sqlClient) {
            sqlClient.runQuery("CREATE ADAPTER SCRIPT ...")
        },
        addInstance(_installation: Installation, _params: ParameterValues, _sql: SqlClient): Instance {
            return undefined;
        },
        findInstallations(_sqlClient: SqlClient): Installation[] {
            return [];
        },
        findInstances(_installation: Installation, _sql: SqlClient): Instance[] {
            return [];
        },
        uninstall(_installation: Installation, _sql: SqlClient): void {
            //empty on purpose
        },
        deleteInstance(_instance: Instance): void {
            //empty on purpose
        },
        readInstanceParameters(_installation: Installation, _instance: Instance, _sqlClient: SqlClient): ParameterValues {
            return undefined;
        }
    }
}

registerExtension(createExtension())
