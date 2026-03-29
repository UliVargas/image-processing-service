package main

import (
	"fmt"
	"os"

	"image-processing-service/internal/modules/file"
	"image-processing-service/internal/modules/session"
	"image-processing-service/internal/modules/user"

	"ariga.io/atlas-provider-gorm/gormschema"
)

func main() {
	stmts, err := gormschema.New("postgres").Load(
		&user.User{},
		&session.Session{},
		&file.File{},
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load gorm schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(stmts)
}
