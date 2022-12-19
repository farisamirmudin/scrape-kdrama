package lib

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func Search() string {
	// Create a new bufio.Reader to read from standard input
	reader := bufio.NewReader(os.Stdin)

	c := color.New(color.FgGreen)

	// Apply the color to a string and get the colored string
	fmt.Print(c.Sprint("Search Drama: "))

	// Read the input until the newline character
	input, err := reader.ReadString('\n')
	if err != nil {
		panic("Error reading input")
	}
	input = strings.TrimSuffix(input, "\n")
	input = strings.ReplaceAll(input, ` `, `-`)

	return input
}
