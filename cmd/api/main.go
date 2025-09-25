package main

import (
	"fmt"
	"flowcargo/internal/app"
)

func main() {
	fmt.Println("Hello from main!")
	err := app.CreateAndRun("./.env.dev")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}