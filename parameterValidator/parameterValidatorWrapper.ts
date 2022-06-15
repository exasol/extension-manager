import {validateParameter} from "@exasol/extension-manager-parameter-validator/dist/extensionParameterValidator";

// @ts-ignore global is defined globally in the VM
global.validateParameter = validateParameter