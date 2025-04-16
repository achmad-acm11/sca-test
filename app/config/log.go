package config

import "github.com/sirupsen/logrus"

func ConfigLog() *logrus.Logger {

	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetLevel(logrus.DebugLevel)

	return logger
}
