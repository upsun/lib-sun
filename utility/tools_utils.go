package utils

import (
	"log"
	"os"
	"time"

	"github.com/getsentry/sentry-go"
	version "github.com/upsun/lib-upsun"
)

func StartReporters(appName string) {
	StartCrashReporter(appName)
	StartAnalyticsReportes(appName)
}

func StartCrashReporter(appName string) {
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://810e5677ad7fb44e87da6d0d2f2d9212@o38989.ingest.us.sentry.io/4507532824936448",
		TracesSampleRate: 1.0,
		Environment:      os.Getenv("APP_ENV"),
		Release:          "adv-initial@" + version.VERSION,
		EnableTracing:    true,
		Debug:            false,
	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("app.name", appName)
	})

	defer sentry.Flush(2 * time.Second)
	//sentry.CaptureMessage("It works 3!")
}

func StartAnalyticsReportes(appName string) {

}
