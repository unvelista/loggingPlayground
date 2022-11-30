package logging

import (
	"io"
	"os"

	"github.com/sirupsen/logrus"
)

func InitLogger(f *os.File) *logrus.Logger {
	logger := logrus.New()
	logger.SetOutput(io.MultiWriter(f, os.Stdout))
	logger.SetLevel(logrus.InfoLevel)
	return logger
}
