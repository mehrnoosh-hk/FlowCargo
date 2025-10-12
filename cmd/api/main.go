package main

import (
	"context"
	"fmt"
	"os"

	"flowcargo/internal/app"

	_ "flowcargo/docs" // Import generated docs
)

// @title           FlowCargo API
// @version         1.0
// @description     Multi-tenant cargo management system API
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.flowcargo.com/support
// @contact.email  support@flowcargo.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /

// @schemes http https

func main() {
	fmt.Println("Hello from main!")
	parsedFlags, err := ParseFlags()
	if err != nil {
		fmt.Fprint(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	ctx := context.Background()
	err = app.CreateAndRun(
		ctx,
		parsedFlags.Environment,
		parsedFlags.ConfigPath,
		app.NewWire(),
	)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
