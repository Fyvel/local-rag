package main

import (
	"fmt"
	"local-ai/api"

	_ "local-ai/docs" // This will import the auto-generated docs
)

func main() {
	r := api.NewRouter()
	fmt.Println("Server starting on :8080")
	r.Run(":8080")
}
