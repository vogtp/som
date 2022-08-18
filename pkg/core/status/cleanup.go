package status

// Cleanup removes groups with Level below OK (i.e. Unkown)
func Cleanup(s Status) {
	grp, ok := s.(*statusGroup)
	if !ok {
		return
	}
	children := make([]Grouper, 0, len(grp.children))
	for _, c := range grp.Szenarios() {
		if c.Level() < OK {
			continue
		}
		clenupSz(c)
		children = append(children, c)
	}
	grp.children = children
}

func clenupSz(grpr SzenarioGroup) {
	grp, ok := grpr.(*szGroup)
	if !ok {
		return
	}
	children := make([]Grouper, 0, len(grp.children))
	for _, c := range grp.Regions() {
		if c.Level() < OK {
			continue
		}
		clenupRg(c)
		children = append(children, c)
	}
	grp.children = children
}

func clenupRg(grpr RegionGroup) {
	grp, ok := grpr.(*regGroup)
	if !ok {
		return
	}
	children := make([]Grouper, 0, len(grp.children))
	for _, c := range grp.Users() {
		if c.Level() < OK {
			continue
		}
		children = append(children, c)
	}
	grp.children = children
}
