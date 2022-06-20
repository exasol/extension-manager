package main

//go:generate swag init -g restAPI/restApi.go -o generatedApiDocs
//go:generate npm --prefix parameterValidator ci
//go:generate npm --prefix parameterValidator run build
