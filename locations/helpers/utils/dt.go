package utils

import (
	dtclient "labs/service-mesh/locations/helpers/clients"
	"strings"
)

func TimezoneToReadable(tz dtclient.Timezone) dtclient.Timezone {
	newTz := dtclient.Timezone{}

	if len(tz.Area) > 0 {
		newTz.Area = strings.ReplaceAll(tz.Area, "_", " ")
	}
	if len(tz.Location) > 0 {
		newTz.Location = strings.ReplaceAll(tz.Location, "_", " ")
	}
	if len(tz.Region) > 0 {
		newTz.Region = strings.ReplaceAll(tz.Region, "_", " ")
	}

	return newTz
}
