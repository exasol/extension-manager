package core

import (
	"github.com/exasol/extension-manager/extensionController"
)

func NewApiContext(controller extensionController.TransactionController) *ApiContext {
	return &ApiContext{Controller: controller}
}

type ApiContext struct {
	Controller extensionController.TransactionController
}
