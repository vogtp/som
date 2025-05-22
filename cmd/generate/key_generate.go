package main

import (
	"errors"
	"fmt"
	"os"
	"text/template"

	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
)

var genKeyCtl = &cobra.Command{
	Use:   "key <package> <file>",
	Short: "Generate key file for encryption",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		pkg := args[0]
		file := args[1]
		_, err := os.Stat(file)
		if err == nil {
			// nothing to do
			return nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("Error reading stat of file %q: %w", file, err)
		}
		return genKeyFile(pkg, file)
	},
}

func genKeyFile(pkg string, file string) error {
	fmt.Printf("Generating keyfile %q for package %q\n", file, pkg)
	tplName := "keyTemplate"
	tpl, err := template.New(tplName).Parse(keyTemplate)
	if err != nil {
		return fmt.Errorf("cannot parse template: %w", err)
	}
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	key, err := genKey()
	if err != nil {
		return fmt.Errorf("cannot generate key: %w", err)
	}
	return tpl.Execute(f, map[string]interface{}{
		"pkg": pkg,
		"key": fmt.Sprintf("%#v", []byte(key)),
	})
}

func genKey() (string, error) {
	pwGen, err := password.NewGenerator(&password.GeneratorInput{})
	if err != nil {
		return "", fmt.Errorf("cannot create generator: %w", err)
	}
	return pwGen.Generate(40, 10, 10, false, true)
}

const keyTemplate = `package {{ .pkg }}

import "github.com/vogtp/som/pkg/core"

// len: 40
var myKey = {{ .key }}

func init() {
	core.Keystore.Add(myKey)
}
`
