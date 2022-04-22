const CURRENT_API_VERSION = "0.1.0";

/**
 * This class represents an extension that can be installed with the extension-manager.
 *
 * Wondering why we picked TypeScript as an interface? Check the design.md / Extension API
 */
export interface ExasolExtension {
    /** Name of the extension */
    name: string;
    /** Description of the extension */
    description: string;
    /** Files that this extension requires in BucketFS. */
    bucketFsUploads?: BucketFSUpload[];
    /**
     * Install this extension.
     *
     * Installing means creating the adapter scripts / UDF definitions.
     *
     * @param sqlClient client for running SQL queries
     */
    install: (sqlClient: SqlClient) => void
    /**
     * Find installations of this extension independent of the version.
     *
     * @param sqlClient client for running SQL queries
     */
    findInstallations: (sqlClient: SqlClient) => Installation[]
    /**
     * Uninstall this extension. (Delete adapter scripts / udf definitions)
     *
     * This method does not delete the instances first. The caller takes care of this.
     *
     * @param installation the installation to uninstall
     * @param sqlClient client for running SQL queries
     */
    uninstall: (installation: Installation, sqlClient: SqlClient) => void
    /** Parameter definitions for an instance of this extension. */
    instanceParameters: Parameter[];
    /**
     * Add an instance of this extension
     *
     * An instance of an extension is for example a Virtual Schema.
     *
     * @param installation installation
     * @param params parameter values
     * @param sqlClient client for running SQL queries
     */
    addInstance: (installation: Installation, params: ParameterValues, sqlClient: SqlClient) => Instance
    /**
     * Find instances of this extension.
     *
     * @param installation installation
     * @param sqlClient client for running SQL queries
     */
    findInstances: (installation: Installation, sqlClient: SqlClient) => Instance[]
    /**
     * Read the parameter values of an instance.
     *
     * @param installation installation
     * @param instance instance
     * @param sqlClient client for running SQL queries
     */
    readInstanceParameters: (installation: Installation, instance: Instance, sqlClient: SqlClient) => ParameterValues
    /**
     * Delete an instance.
     *
     * @param installation installation
     * @param instance instance to delete
     * @param sqlClient client for running SQL queries
     */
    deleteInstance: (installation: Installation, instance: Instance, sqlClient: SqlClient) => void
}

/**
 * Reference to an installation of this extension.
 */
export interface Installation {

}
/**
 * Reference to an instance of this extension.
 */
export interface Instance {
    name: string
}

/**
 * Map of parameter name -> parameter value.
 */
export interface ParameterValues {
    [index: string]: string;
}

/**
 * Simple SQL client.
 */
export interface SqlClient {
    /**
     * Run a SQL query.
     * @param query sql query string
     */
    runQuery: (query: string) => void
}


/**
 * Description of a file that needs to be uploaded to BucketFS.
 */
export interface BucketFSUpload {
    /** Human-readable name or short description of the file */
    name: string
    downloadUrl: string
    licenseUrl: string
    /** Default: false */
    licenseAgreementRequired?: boolean
    bucketFsFilename: string
}

/**
 * Abstract base for parameters.
 */
interface BaseParameter {
    id: string
    name: string
    type: string
    condition?: Condition
    default?: string
    placeholder?: string
    readOnly?: boolean
    required?: boolean
}

/**
 * String parameter.
 */
export interface StringParameter extends BaseParameter {
    type: "string"
    regex?: RegExp
}

/**
 * Parameter type.
 */
export type Parameter = StringParameter | SelectParameter;

/**
 * Type for a map for select options.
 * Map: value to select -> display name
 */
export interface OptionsType{
    [index:string]: string
}

/**
 * Parameter that allows to select a value from a list.
 */
export interface SelectParameter extends BaseParameter {
    options: OptionsType
    type: "select"
}

/**
 * Condition for conditional parameters.
 */
export type Condition = Comparison | And | Or;

/**
 * Comparison operators
 */
export enum Operators {
    EQ, LESS, GREATER, LESS_EQUAL, GREATER_EQ
}

/**
 * Comparison of a parameter value and given value.
 */
export interface Comparison {
    /** parameter name */
    parameter: string
    /** Value to compare with */
    value: string | number
    operator: Operators
}

/** And predicate */
export interface And {
    and: Condition[]
}

/** Or predicate */
export interface Or {
    or: Condition[]
}

/** Not predicate */
export interface Not {
    not: Condition
}

/**
 * This method registers an extension at the GO JS runtime.
 *
 * @param extensionToRegister extension to register
 */
export function registerExtension(extensionToRegister: ExasolExtension): void {
    // @ts-ignore //this is a global variable defined in the nested JS VM in the backend
    installedExtension = {
        extension: extensionToRegister,
        apiVersion: CURRENT_API_VERSION
    };
}
