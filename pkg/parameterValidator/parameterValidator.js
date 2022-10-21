(() => {
  // node_modules/@exasol/extension-parameter-validator/dist/extensionParameterValidator.js
  const SUCCESS_RESULT = { success: true, message: "" };
  function validationError(errorMessage) {
    return { success: false, message: errorMessage };
  }
  function validateParameter(definition, value) {
    if (value === void 0 || value === null || value === "") {
      if (definition.required) {
        return validationError("This is a required parameter.");
      } else {
        return SUCCESS_RESULT;
      }
    } else {
      switch (definition.type) {
        case "string":
          return validateStringParameter(definition, value);
        case "boolean":
          return validateBooleanParameter(value);
        default:
          return validationError("unsupported parameter type '" + definition.type + "'");
      }
    }
  }
  function validateStringParameter(definition, value) {
    if (definition.regex) {
      if (!new RegExp(definition.regex).test(value)) {
        return validationError("The value has an invalid format.");
      }
    }
    return SUCCESS_RESULT;
  }
  function validateBooleanParameter(value) {
    if (value === "true" || value === "false") {
      return SUCCESS_RESULT;
    }
    return validationError("Boolean value must be 'true' or 'false'.");
  }

  // dist/parameterValidatorWrapper.js
  global.validateParameter = validateParameter;
})();
