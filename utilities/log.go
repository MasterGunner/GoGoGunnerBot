// Package utilities contains utility functions for GunnerBot
package utilities

import (
	"fmt"
	"time"
)

// Log is a general logging function (to be expanded)
func Log(msg string) {
	date := time.Now()
	//date := time.Now().String()[0:20]
	output := fmt.Sprintf("%v : %s", date, msg)
	fmt.Println(output)
}
