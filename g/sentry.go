package g

import (
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitSentry() {
	// init sentry before init gin
	if err := sentry.Init(sentry.ClientOptions{
		Dsn: viper.GetString("sentry.dsn"),
	}); err != nil {
		log.Printf("Sentry initialization failed: %v", err)
	}
}
