import {validateParameter, validateParameters} from "./extensionParameterValidator";

describe("extensionParameterValidator", () => {
    describe("validateParameter", () => {
        it.each`
        parameter  | value |  expectedResult
        ${{type: "string", regex: "^a+$"}}   |${"test"} | ${{success: false, message: "The value has an invalid format."} }
        ${{type: "string", regex: "^t+$"}} |${"test"}  | ${{success: false, message: "The value has an invalid format."} }
        ${{type: "string", regex: "^test$"}}  |${"test"} | ${{success: true, message: ""} }
        ${{type: "string", regex: "^.*$"}} |${"test"}  | ${{success: true, message: ""} }
        ${{type: "string"}}  |${"test"} | ${{success: true, message: ""} }
        ${{type: "string", required: true}}  |${""} | ${{success: false, message: "This is a required filed."} } 
        `('validates $parameter as $result', ({parameter, value, expectedResult}) => {
            let result = validateParameter(parameter, value);
            expect(result).toEqual(expectedResult)
        })
    })

    describe("validateParameters", () => {
        it("detects a missing parameter", () => {
            let result = validateParameters([{id: "param1", type: "string", name: "Parameter 1", required: true}], {});
            expect(result).toEqual({success: false, message: "This is a required filed."})
        })

        it("accepts a valid parameter", () => {
            let result = validateParameters([{id: "param1", type: "string", name: "Parameter 1", required: true}], {param1: "test"});
            expect(result).toEqual({success: true, message: ""})
        })

        it("rejects invalid parameters", () => {
            let result = validateParameters([{id: "param1", type: "string", name: "Parameter 1", regex: "^a+$"},
                {id: "param2", type: "string", name: "Parameter 2", regex: "^a+$"}], {param1: "test", param2: "test"});
            expect(result).toEqual({success: false, message: "Parameter 1: The value has an invalid format.\n" +
                    "Parameter 2: The value has an invalid format."})
        })
    })
})
