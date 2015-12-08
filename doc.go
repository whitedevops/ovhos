/*
Package ovhos is a light SDK for the OVH Object Storage service.

Installation

In the terminal:

	$ go get github.com/whitedevops/ovhos

Usage

Example:

	package main

	import (
		"fmt"
		"os"

		"github.com/whitedevops/ovhos"
	)

	// Make a new OVH Object Storage client
	var storage = &ovhos.Client{
		Region:    "XXXX",
		Container: "X",
		TenantID:  "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
		Username:  "XXXXXXXXXXXX",
		Password:  "XXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX",
	}

	func main() {
		// Check connection
		if err := storage.Ping(); err != nil {
			panic(err)
		}

		// Upload an object
		f, err := os.Open("file.txt")
		if err != nil {
			panic(err)
		}
		if err := storage.Upload("file.txt", f); err != nil {
			panic(err)
		}
		f.Close()

		// Get the URL of an object
		fmt.Println(storage.URL("file.txt"))

		// List all objects
		l, err := storage.List()
		if err != nil {
			panic(err)
		}
		fmt.Println(l)

		// Delete an object
		if err := storage.Delete("file.txt"); err != nil {
			panic(err)
		}
	}
*/
package ovhos
