import { validateParameter } from "@exasol/extension-parameter-validator";

// [impl -> dsn~reuse-parameter-validation-rules~1]
// [impl -> dsn~parameter-validation-rules-simple~1]
// @ts-ignore global is defined globally in the VM
global.validateParameter = validateParameter
