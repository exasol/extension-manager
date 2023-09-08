package restAPI

import (
	"github.com/exasol/extension-manager/pkg/extensionController"
)

func NewApiContext(controller extensionController.TransactionController, addCauseToInternalServerError bool) *ApiContext {
	return &ApiContext{
		Controller:                    controller,
		addCauseToInternalServerError: addCauseToInternalServerError,
	}
}

type ApiContext struct {
	Controller                    extensionController.TransactionController
	addCauseToInternalServerError bool
}
