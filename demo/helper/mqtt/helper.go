package custmqtt

import "strings"

func ParseTopic(rawTopic string) []string {
	return strings.Split(rawTopic, "/")
}

