// This file informs the "mage" command how to build and publish the
// send-email and send-raw-email commands.
//
// $ mage [build]   builds the code
// $ mage install   builds the code and install on server
//
//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/target"
)

var linux = map[string]string{"GOOS": "linux"}

var Default = Build

// Backend builds all of the back-end binaries.
func Build() {
	mg.Deps(SendEmail, SendRawEmail)
}

// SendEmail builds and installs the send-email program.
func SendEmail() error {
	if changed, err := target.Dir("dist/send-email", "send-email", "scholaemail.go"); err != nil {
		return err
	} else if !changed {
		return nil
	}
	return sh.RunWith(linux, "go", "build", "-o", "dist/send-email", "./send-email")
}

// SendRawEmail builds and installs the send-email program.
func SendRawEmail() error {
	if changed, err := target.Dir("dist/send-raw-email", "send-raw-email", "scholaemail.go"); err != nil {
		return err
	} else if !changed {
		return nil
	}
	return sh.RunWith(linux, "go", "build", "-o", "dist/send-raw-email", "./send-raw-email")
}

func Install() (err error) {
	mg.Deps(Build)
	return sh.Run("scp", "dist/send-email", "dist/send-raw-email", "schola:bin")
}
