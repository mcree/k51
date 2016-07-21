package backend

import (
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetFormatter(new(prefixed.TextFormatter))
}
