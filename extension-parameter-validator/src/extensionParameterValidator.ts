import {Parameter, ParameterValues, StringParameter} from "@exasol/extension-manager-interface";

const SUCCESS_RESULT = {success: true, message: ""};

function validationError(errorMessage: string): ValidationResult {
    return {success: false, message: errorMessage}
}

export function validateParameter(definition: Parameter, value: string): ValidationResult {
    if (value === undefined || value === null || value == "") {
        if (definition.required) {
            return validationError("This is a required filed.")
        } else {
            return SUCCESS_RESULT;
        }
    } else {
        switch (definition.type) {
            case "string":
                return validateStringParameter(definition, value);
            default:
                return validationError("unsupported parameter type '" + definition.type + "'");
        }
    }
}

export function validateParameters(definitions: Parameter[], values: ParameterValues): ValidationResult {
    let findings: string[] = []
    for (const key in definitions) {
        let definition = definitions[key];
        let singleResult = validateParameter(definition, values[definition.id])
        if (!singleResult.success) {
            findings.push(definition.name + ": " + singleResult.message)
        }
    }
    if (findings.length == 0) {
        return SUCCESS_RESULT
    } else {
        return validationError(findings.join("\n"))
    }
}

function validateStringParameter(definition: StringParameter, value: string) {
    if (definition.regex !== null) {
        if (!new RegExp(definition.regex).test(value)) {
            return validationError("The value has an invalid format.")
        }
    }
    return SUCCESS_RESULT
}


export interface ValidationResult {
    /** true of the validation passed with no findings. */
    success: boolean
    /** Validation error description. If multiples errors were found they are separated by \n. */
    message: string
}
