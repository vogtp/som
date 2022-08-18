package status

import (
	"fmt"

	"github.com/vogtp/go-hcl"
)

// GrpReg holds all by json instanceble groups
var GrpReg *grpReg = &grpReg{
	grps: map[string]Grouper{},
}

type grpReg struct {
	grps map[string]Grouper
}

// Add a group (best in init())
func (r *grpReg) Add(g Grouper) {
	hcl.Debugf("Adding group to registry: %T\n", g)
	r.grps[fmt.Sprintf("%T", g)] = g
}

func (r *grpReg) new(n string, key string) (Grouper, error) {
	g, ok := r.grps[n]
	if !ok {
		return nil, fmt.Errorf("no such group %s in reqistry", n)
	}
	return g.New(key), nil
}

func init() {
	GrpReg.Add(&Group{})
	GrpReg.Add(&evtGroup{})
	GrpReg.Add(&szGroup{})
	GrpReg.Add(&regGroup{})
	GrpReg.Add(&usrGroup{})
}
