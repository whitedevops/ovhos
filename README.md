<p align="center"><img src="https://cloud.githubusercontent.com/assets/9503891/11653313/a2d910ba-9d9d-11e5-9cf9-cd306ddda976.png" alt="OVHOS" title="OVHOS"><br><br></p>

OVHOS is a light SDK for the [OVH Object Storage](https://www.ovh.com/fr/cloud/storage/object-storage) service.

## Installation

```Shell
$ go get github.com/whitedevops/ovhos
```

## Usage [![GoDoc](https://godoc.org/github.com/whitedevops/ovhos?status.svg)](https://godoc.org/github.com/whitedevops/ovhos)

```Go
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
```
