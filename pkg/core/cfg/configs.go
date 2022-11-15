package cfg

import (
	"time"

	"github.com/spf13/viper"
)

// config keys
const (
	// PromURL prometheus url
	PromURL = "prometheus.url"
	// PromBasePath the base path of the prometheus URL
	PromBasePath = "prometheus.basepath"
	// CoreStartdelay lets the core wait for the given duration
	CoreStartdelay = "core.startdelay"

	// BusWsPath is the url path the bus listens ons
	BusWsPath = "bus.ws.path"

	// AlertDelay is the timespan (time.Duration) that must have passed in order to gernerate an alert
	AlertDelay = "alert.delay"
	// AlertIntervall is the intervall (time.Duration) after which an alert is resend
	AlertIntervall = "alert.intervall"
	// AlertLevel is the min level for an alert to get escalated
	AlertLevel = "alert.level"
	// AlertSubject (Mail) Subject of the alert
	AlertSubject = "alert.subject"
	// AlertMailFrom from mail adr
	AlertMailFrom = "alert.mail.from"
	// AlertMailSMTPHost smtp host
	AlertMailSMTPHost = "alert.mail.smtp.host"
	// AlertMailSMTPPort smtp port
	AlertMailSMTPPort = "alert.mail.smtp.port"
	// AlertIncidentCorrelationEvents is the number of events used for correlation
	AlertIncidentCorrelationEvents = "alert.event.correlation.events"
	// AlertIncidentCorrelationReopenTime is the duration during which a icident is reopened
	AlertIncidentCorrelationReopenTime = "alert.event.correlation.reopentime"

	// AlertDestinations defines destinations for alerting (handled in alerter)
	AlertDestinations = "alert.destinations"
	// AlertRules defines rules alert handling (handled in alerter)
	AlertRules = "alert.rules"

	// AlertVisualiserURL is the URL of the webserver serving details
	AlertVisualiserURL = "alert.visualiserURL"

	// StatusTimeout is the duration after which a status will be unknow if no event is received
	StatusTimeout = "status.timeout"
	// StatusCleanup is the duration after which a status will be removed if no event is received
	StatusCleanup = "status.cleanup"

	// PasswdChange should the password be changed
	PasswdChange = "password.change"
	// PasswdChgSz szeanrios to change the password
	PasswdChgSz = "password.szenarios"
	// PasswdRuleUpper upper case letters to be used in password ("" use default)
	PasswdRuleUpper = "password.rule.upper"
	// PasswdRuleLower lower case letters to be used in password ("" use default)
	PasswdRuleLower = "password.rule.lower"
	// PasswdRuleDigit digits to be used in password ("" use default)
	PasswdRuleDigit = "password.rule.digits"
	// PasswdRuleSymbols symbols to be used in password ("" use default)
	PasswdRuleSymbols = "password.rule.symbols"
	// PasswdRuleLength length of the password
	PasswdRuleLength = "password.rule.length"
	// PasswdRuleNumSymbols number of digits in the password
	PasswdRuleNumDigits = "password.rule.numDigits"
	// PasswdRuleNumSymbols number of symbols in the password
	PasswdRuleNumSymbols = "password.rule.numSymbols"
)

func init() {
	viper.SetDefault(PasswdChange, false)
	viper.SetDefault(PasswdChgSz, []string{})
	viper.SetDefault(PasswdRuleUpper, "")
	viper.SetDefault(PasswdRuleLower, "")
	viper.SetDefault(PasswdRuleDigit, "")
	viper.SetDefault(PasswdRuleLength, 20)
	viper.SetDefault(PasswdRuleNumDigits, 3)
	viper.SetDefault(PasswdRuleNumSymbols, 3)

	viper.SetDefault(PromURL, "http://localhost:9090/")
	viper.SetDefault(PromBasePath, "/")
	viper.SetDefault(CoreStartdelay, 100*time.Millisecond)
	viper.SetDefault(BusWsPath, "/meta/message")

	viper.SetDefault(AlertLevel, "warning")
	viper.SetDefault(AlertDelay, 15*time.Minute)
	viper.SetDefault(AlertIntervall, 14*24*time.Hour)
	viper.SetDefault(AlertSubject, "SOM Alert:")

	viper.SetDefault(AlertMailFrom, "alert@som-monitoring.net")
	viper.SetDefault(AlertMailSMTPPort, 25)
	viper.SetDefault(AlertIncidentCorrelationEvents, 6)
	viper.SetDefault(AlertIncidentCorrelationReopenTime, time.Hour)

	viper.SetDefault(StatusTimeout, 6*time.Hour)
	viper.SetDefault(StatusCleanup, 7*24*time.Hour)

}
