package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/shirou/gopsutil/v3/host"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Constants
// You can add as many constants as you want here and set them to anything
func getConstants() map[string]string {
	var s strings.Builder
	hostInfo, _ := host.Info()

	// uptime
	duration := time.Duration(hostInfo.Uptime) * time.Second

	days := fmt.Sprint(int(duration.Hours()) / 24)
	hours := fmt.Sprint(int(duration.Hours()) % 24)
	minutes := fmt.Sprint(int(duration.Minutes()) % 60)
	seconds := fmt.Sprint(int(duration.Seconds()) % 60)

	// you can manually adjust how everything is formatted here
	s.WriteString(days + " Days, ")
	s.WriteString(hours + " Hours, ")
	s.WriteString(minutes + " Minutes, ")
	s.WriteString(seconds + " Seconds")
	uptime := s.String()


	// local time
	time := time.Now().Format("15:04:05");

	// hostname
	hostname := hostInfo.Hostname

	// username
	username := os.Getenv("USER")
	if username == "" {
		username = os.Getenv("USERNAME") // windows fallback
	}

	// OS Platform
	platform := cases.Title(language.English, cases.Compact).String(hostInfo.OS)

	// OS name
	version := hostInfo.KernelVersion

	// all constants
	return map[string]string{
		"{uptime}":   uptime,
		"{time}":     time,
		"{hostname}": hostname,
		"{username}": username,
		"{platform}": platform,
		"{version}": version,
	}
}

func ProcessConstants(v any, constants map[string]string) {
	val := reflect.ValueOf(v)

	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)

		if !field.CanSet() {
			continue
		}

		switch field.Kind() {
		case reflect.String:
			original := field.String()
			processed := replaceConstants(original, constants)
			field.SetString(processed)

		case reflect.Struct:
			ProcessConstants(field.Addr().Interface(), constants)

		case reflect.Ptr:
			if !field.IsNil() && field.Elem().Kind() == reflect.Struct {
				ProcessConstants(field.Interface(), constants)
			}

		case reflect.Slice:
			if field.Type().Elem().Kind() == reflect.Struct {
				for j := 0; j < field.Len(); j++ {
					ProcessConstants(field.Index(j).Addr().Interface(), constants)
				}
			}
		}
	}
}

func replaceConstants(input string, constants map[string]string) string {
	for placeholder, value := range constants {
		input = strings.ReplaceAll(input, placeholder, value)
	}
	return input
}
