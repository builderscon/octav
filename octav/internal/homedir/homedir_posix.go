// +build !darwin,!windows

package homedir

import (
	"fmt"
	"os"
)

func Get() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("error: Environment variable HOME not set")
	}

	return home, nil
}
