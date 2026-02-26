/*
The lint_commit tool is invoked by lefthook to validate commit messages
against the Conventional Commits specification.  It reads the temporary file
created by git and exits with an error if the message does not conform.
*/
package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	// Patrón Regex para Conventional Commits
	commitPattern = `^(feat|fix|docs|style|refactor|perf|test|chore|ci|build)(\([a-z0-9-]+\))?: .+$`
	// Límite recomendado de caracteres para la primera línea
	maxLineLength = 72
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go run scripts/lint_commit.go <archivo_mensaje>")
		os.Exit(1)
	}

	commitMsgFile := os.Args[1]
	content, err := os.ReadFile(commitMsgFile)
	if err != nil {
		fmt.Printf("Error al leer el mensaje: %v\n", err)
		os.Exit(1)
	}

	// Solo validamos la primera línea del commit
	lines := strings.Split(string(content), "\n")
	firstLine := strings.TrimSpace(lines[0])

	// 1. Validar longitud
	if len(firstLine) > maxLineLength {
		fmt.Printf("❌ ERROR: El mensaje es demasiado largo (%d/%d caracteres).\n", len(firstLine), maxLineLength)
		os.Exit(1)
	}

	// 2. Validar formato con Regex
	re := regexp.MustCompile(commitPattern)
	if !re.MatchString(firstLine) {
		fmt.Println("❌ ERROR: El formato no sigue el estándar 'tipo(alcance): descripción'.")
		fmt.Println("Tipos válidos: feat, fix, docs, style, refactor, perf, test, chore, ci, build")
		os.Exit(1)
	}

	fmt.Println("✅ Mensaje de commit impecable.")
}
