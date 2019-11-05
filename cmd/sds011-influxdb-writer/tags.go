package main

import (
	"strings"
)

func toTags(tags string) map[string]string {
	result := map[string]string{}
	for _, tagString := range strings.Split(tags, ",") {
		tagString = strings.TrimSpace(tagString)
		tagNameValue := strings.SplitN(tagString, "=", 2)
		if len(tagNameValue) >= 2 {
			result[tagNameValue[0]] = tagNameValue[1]
		}
	}
	return result
}