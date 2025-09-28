package main

import (
	"fmt"
	"indexing/api"

	_ "indexing/docs" // This will import the auto-generated docs
)

func main() {
	r := api.NewRouter()
	fmt.Println("Server starting on :8080")
	r.Run(":8080")
}
