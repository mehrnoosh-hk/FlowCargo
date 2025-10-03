package main

import (
	"context"
	"flowcargo/internal/app"
	"flowcargo/internal/shared/config"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello from main!")
	ctx := context.Background()
	configFile := "config" //TODO: Implement getting address from runtime flags
	err := app.CreateAndRun(ctx, config.Production, &configFile)
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
