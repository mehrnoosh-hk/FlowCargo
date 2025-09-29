package main

import (
	"context"
	"fmt"
	"flowcargo/internal/app"
)

func main() {
	fmt.Println("Hello from main!")
	ctx := context.Background()
	err := app.CreateAndRun(ctx, "./.env.dev")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}