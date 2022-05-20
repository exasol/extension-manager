package main

import (
	"backend/respApi"
)

func main() {
	restApi := respApi.RestApi{}
	restApi.Serve()
}
