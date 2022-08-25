package requests

import "github.com/exasol/extension-manager/restAPI/core"

var authentication = map[string][]string{core.BasicAuth: {}, core.BearerAuth: {}}
