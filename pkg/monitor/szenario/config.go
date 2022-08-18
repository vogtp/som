package szenario

const (
	// UserTypeAll name of the usertype that contains all szenarios
	UserTypeAll = "all"
)

var (
	// NoConfig means that the sznario config has not been loaded and set
	NoConfig *Config = New()
)

// Config holds all Szenario informations
type Config struct {
	userTypes map[string]*UserType
	allSz     *UserType
}

// New creates a new config
func New() *Config {
	c := &Config{
		userTypes: make(map[string]*UserType),
	}
	c.allSz = MustUserType(c.CreateUsertType(UserTypeAll, "usertype containing all szenarios"))
	return c
}

// SzenarioCount returns the number of szenarios
func (c Config) SzenarioCount() int {
	return len(c.allSz.Szenarios)
}
