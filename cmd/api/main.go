package main

import (
	"context"
	"flowcargo/internal/app"
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello from main!")
	parsedFlags, err := ParseFlags()
	if err != nil {
		fmt.Fprint(os.Stderr, "Error:", err)
		os.Exit(1)
	}
	ctx := context.Background()
	err = app.CreateAndRun(ctx, parsedFlags.Environment, parsedFlags.ConfigPath)
	if err != nil {
		fmt.Fprint(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}
