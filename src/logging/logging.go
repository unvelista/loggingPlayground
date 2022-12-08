package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger() *logrus.Entry {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{
		DisableTimestamp: true,
		PrettyPrint:      true, // pretty print breaks parsing in Fluent Bit
	})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)
	logger.SetReportCaller(true)

	logEntry := setDefaultFields(logger)

	return logEntry
}

func setDefaultFields(logger *logrus.Logger) *logrus.Entry {
	return logger.WithFields(logrus.Fields{
		"defaultField1": "foo",
	})
}
