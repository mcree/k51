#!/bin/bash
set -e
go get -u github.com/FiloSottile/gvt
gvt fetch github.com/spf13/cobra
gvt fetch github.com/go-kit/kit/log/levels
gvt fetch github.com/jehiah/go-strftime
gvt fetch github.com/eclipse/paho.mqtt.golang
gvt fetch github.com/Sirupsen/logrus
gvt fetch github.com/x-cray/logrus-prefixed-formatter
