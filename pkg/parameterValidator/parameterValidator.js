(() => {
  // node_modules/@exasol/extension-parameter-validator/dist/extensionParameterValidator.js
  var SUCCESS_RESULT = { success: true };
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
      const definitionType = definition.type;
      switch (definition.type) {
        case "string":
          return validateStringParameter(definition, value);
        case "boolean":
          return validateBooleanParameter(value);
        case "select":
          return validateSelectParameter(definition, value);
        default:
          return validationError(`unsupported parameter type '${definitionType}'`);
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
  function validateSelectParameter(definition, value) {
    const possibleValues = definition.options.map((option) => option.id);
    if (possibleValues.length === 0) {
      return validationError("No option available for this parameter.");
    }
    if (possibleValues.includes(value)) {
      return SUCCESS_RESULT;
    }
    const quotedValues = possibleValues.map((value2) => `'${value2}'`).join(", ");
    return validationError(`The value is not allowed. Possible values are ${quotedValues}.`);
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
