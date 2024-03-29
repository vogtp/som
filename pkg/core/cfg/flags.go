package cfg

import (
	"time"

	"github.com/spf13/pflag"
)

// flag keys
const (
	// BrowserShow show the browser window
	BrowserShow = "browser.show"
	// BrowserNoClose Do not close the browser window in the end. Implies show, timeout 10m and no repeat
	BrowserNoClose = "browser.noclose"
	// CheckTimeout Check timeout (time.Duration)
	CheckTimeout = "check.timeout"
	// CheckRepeat Check intervall (time.Duration)
	CheckRepeat = "check.repeat"
	// CheckStepDelay a delay between steps (time.Duration)
	CheckStepDelay = "check.step.delay"
	// CheckUser the user the check is run with (probably temp workaround) -- not parsed as flag globally
	CheckUser = "check.user"
	// CheckRegion region the checks runs in
	CheckRegion = "check.region"
	// DataDir Folder to save output like screenshots in
	DataDir = "data.dir"
	// LogLevel error warn info debug
	LogLevel = "log.level"
	// LogSource should we log the source
	LogSource = "log.source"
	// LogJSON log in json
	LogJSON = "log.json"
	// WebURLBasePath the base path of the URL
	WebURLBasePath = "web.urlpath"
	// WebPort the port the webserver runs on
	WebPort = "web.port"

	// BusLogLevel sets the level the bus logs (default: off)
	BusLogLevel = "bus.log.level"
	// BusEndpoints are the endpoints the bus connects to
	BusEndpoints = "bus.endpoint"

	// AlertEnabled is used to disable alerting globally
	AlertEnabled = "alert.enabled"

	// CfgFile
	CfgFile = "config.file"

	// CfgSave triggers periodic config saves
	CfgSave = "config.save"

	// PasswdChangeInitalDelay initial delay of password change (ONLY FOR DEBUGGING)
	PasswdChangeInitalDelay = "password.delay"
)

func init() {
	pflag.Bool(BrowserShow, false, "Show the browser window")
	pflag.Bool(BrowserNoClose, false, "Do not close the browser window in the end. Implies show, timeout 10m  and no repeat")
	pflag.Duration(CheckTimeout, 60*time.Second, "Check timeout")
	pflag.Duration(CheckRepeat, 0, "Check intervall (e.g. 5m)")
	pflag.Duration(CheckStepDelay, 0, "Delay between steps (e.g. 100ms)")
	pflag.String(CheckRegion, "default", "The region the check runs in")
	pflag.String(DataDir, "data", "Folder to save output like screenshots in")
	pflag.String(LogLevel, "warn", "Set the loglevel: error warn info debug trace off")
	pflag.String(WebURLBasePath, "", "the base path of the URL")
	pflag.Int(WebPort, 0, "Port the webserver runs on")
	pflag.String(BusLogLevel, "off", "Log level of the bus")
	pflag.Bool(LogSource, true, "Log the source line")
	pflag.Bool(LogJSON, false, "Log in json")
	pflag.Bool(CfgSave, false, "Should the configs be written to file periodically")
	pflag.Bool(AlertEnabled, true, "Disable alerting")
	pflag.StringSlice(BusEndpoints, nil, "List of external endpoints (e.g. localhost:8080/meta/message) use multiple times to add multiple endpoints")
	pflag.String(CfgFile, "som.yml", "File with the config to load")
	pflag.Duration(PasswdChangeInitalDelay, -1, "initial delay of password change (ONLY FOR DEBUGGING)")
}
