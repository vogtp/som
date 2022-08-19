package term

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
)

// Read a value from the teminal or the args
func Read(name string, args []string, idx int) string {
	if idx < len(args) && len(args[idx]) > 0 {
		return args[idx]
	}
	r := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", name)
	val, err := r.ReadString('\n')
	if err != nil {
		fmt.Printf("Error reading: %v", err)
		return ""
	}
	return strings.TrimSpace(val)
}

// Password read a password from the teminal or the args
func Password(name string, args []string, idx int) string {
	if idx < len(args) && len(args[idx]) > 0 {
		return args[idx]
	}
	fmt.Printf("%s: ", name)
	pw, err := gopass.GetPasswdMasked()
	if err != nil {
		fmt.Printf("Error reading: %v", err)
		return ""
	}
	return strings.TrimSpace(string(pw))
}
