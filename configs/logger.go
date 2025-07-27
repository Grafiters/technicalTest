package configs

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

type LoggerFormat struct {
	infoColor    *color.Color
	warningColor *color.Color
	errorColor   *color.Color
}

func NewLogger() *LoggerFormat {
	return &LoggerFormat{
		infoColor:    color.New(color.FgGreen),
		warningColor: color.New(color.FgYellow),
		errorColor:   color.New(color.FgRed),
	}
}

func (l *LoggerFormat) Info(a interface{}) {
	l.log(l.infoColor, "INFO: ", a)
}

func (l *LoggerFormat) Warning(a interface{}) {
	l.log(l.warningColor, "WARNING: ", a)
}

func (l *LoggerFormat) Error(message string, err interface{}) {
	l.log(l.errorColor, "ERROR", fmt.Sprintf("%s: %v", message, err))
}

func (l *LoggerFormat) log(c *color.Color, level string, data interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	var msg string

	switch v := data.(type) {
	case string:
		msg = v
	case []interface{}:
		msg = fmt.Sprintf("%v", v)
	default:
		msg = fmt.Sprintf("%+v", v)
	}

	c.Printf("%s [%s] %s\n", timestamp, level, msg)
}
