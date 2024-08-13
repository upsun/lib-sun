package utils

import "testing"

const APP_NAME = "test"

func TestStartReporters(t *testing.T) {
	StartReporters(APP_NAME)
}

func TestStartCrashReporter(t *testing.T) {
	StartCrashReporter(APP_NAME)
}

func TestStartAnalyticsReportes(t *testing.T) {
	StartAnalyticsReportes(APP_NAME)
}
