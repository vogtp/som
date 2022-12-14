package cdp

import (
	"fmt"
	"strings"

	"github.com/vogtp/som/pkg/monitor/szenario"
)

// Dump prints the body and its size to log
func (cdp *Engine) Dump() szenario.CheckFunc {
	return func(body *string) error {
		cdp.hcl.Infof("Got Body:\n%s\nSize: %v", *body, len(*body))
		return nil
	}
}

// Contains looks for a string in the body
func (cdp *Engine) Contains(s string) szenario.CheckFunc {
	return func(body *string) error {
		if !strings.Contains(*body, s) {
			return fmt.Errorf("%s not found in body", s)
		}
		cdp.hcl.Infof("Found %s in body", s)
		return nil
	}
}

// NotContains looks for a string in the body and errs if found
func (cdp *Engine) NotContains(s string) szenario.CheckFunc {
	return func(body *string) error {
		if strings.Contains(*body, s) {
			return fmt.Errorf("%s is shown", s)
		}
		cdp.hcl.Infof("Not found %s in body", s)
		return nil
	}
}

// Bigger checks if the size of the body (in bytes) in bigger than i
func (cdp *Engine) Bigger(i int) szenario.CheckFunc {
	return func(body *string) error {
		if len(*body) < i {
			return fmt.Errorf("body is %v should be bigger %v", len(*body), i)
		}
		cdp.hcl.Infof("Body size OK: %v > %v", len(*body), i)
		return nil
	}
}

// Strings gets the body as plaintext
func (cdp *Engine) Strings(html *string) szenario.CheckFunc {
	return func(body *string) error {
		*html = *body
		return nil
	}
}
