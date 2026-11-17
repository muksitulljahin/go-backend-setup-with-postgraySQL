package main

import (
	"fmt"
	"io"
	"os"

	"ariga.io/atlas-provider-gorm/gormschema"
	imagecollection "github.com/go-backend-setup-with-postgray/api/v1/imageCollection"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(&imagecollection.ImageCollection{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}
	io.WriteString(os.Stdout, stmts)
}
