// Package version provides version information for the go-file-organizer CLI tool.
package version

import "fmt"

// Version represents the current version of the application
const Version = "v1.2.1"

// AppName is the name of the application
const AppName = "go-file-organizer"

// GetVersionInfo returns formatted version information
func GetVersionInfo() string {
	return fmt.Sprintf("%s version %s", AppName, Version)
}
