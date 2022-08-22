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
	// AlertMailTo receipient of main alert
	AlertMailTo = "alert.mail.to"
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

	// AlertTeamsWebhook webhook of teams messages
	AlertTeamsWebhook = "alert.teams.webhook"
)

func init() {

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

}
