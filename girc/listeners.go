package girc

import (
	"regexp"
)

type listener struct {
	name     string
	query    *regexp.Regexp       //regex
	callback func(*IRC, []string) //the function to execute, no idea how this actually would work in Go. Check to see hwo I did this in Ruby?
	helptext string
}
