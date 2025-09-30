package main

import (
	"context"
	"flowcargo/internal/app"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello from main!")
	ctx := context.Background()
	err := app.CreateAndRun(ctx, "./.env.dev")
	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
