package main

import "github.com/elspasial/cmd"

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
func main() {
	cmd.Execute()
}
