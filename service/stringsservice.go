package service

import (
	"golang.org/x/exp/slices"
	"regexp"
	"strings"
)

var PipaDefinitions = []string{"пипа", "пипы", "пипе", "пипам", "пипк", "пипо", "пипу", "пипи"}

var HehDefinitions = []string{"haha", "hehe", "хаха", "хах", "хехе", "хех"}

var nonAlphanumericRegex = regexp.MustCompile(`[^\p{L}\p{N} ]+`)

func clearString(str string) string {
	return nonAlphanumericRegex.ReplaceAllString(str, "")
}

func ContainsPipa(text string) bool {
	return ContainsString(text, PipaDefinitions)
}

func ContainsHeh(text string) bool {
	return ContainsString(text, HehDefinitions)
}

func ContainsString(text string, slice []string) bool {
	for _, val := range slice {
		if strings.Contains(strings.ToLower(strings.ReplaceAll(text, " ", "")), val) {
			return true
		}
	}
	return false
}

func ContainsYes(text string) bool {
	if slices.Contains(strings.Split(strings.ToLower(clearString(text)), " "), "да") {
		return true
	}
	return false
}
