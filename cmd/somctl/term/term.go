package term

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/howeyc/gopass"
)

// ReadOrArgs a value from the teminal or the args
func ReadOrArgs(name string, args []string, idx int, defaultValue string) string {
	if args != nil && idx > -1 &&
		idx < len(args) && len(args[idx]) > 0 {
		return args[idx]
	}
	if len(defaultValue) > 0 {
		name = fmt.Sprintf("%s (%s)", name, defaultValue)
	}
	val := Read(name)

	if len(defaultValue) > 0 && len(val) < 1 {
		val = defaultValue
	}
	return val
}

// Read a value from the teminal
func Read(msg string) string {
	r := bufio.NewReader(os.Stdin)
	fmt.Printf("%s: ", msg)
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
