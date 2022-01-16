package log

import (
	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

func main() {
	// Log as JSON instead of the default ASCII formatter.
	Logger.SetFormatter(&logrus.JSONFormatter{})
}
