package conf

import (
	"path/filepath"
	"testing"
	"github.com/coocood/assrt"
)

func parser() *TomlConfigParser {
	return &TomlConfigParser{}
}

func TestTomlParser(t *testing.T) {
	//Testing library, current directory, and some config strings
	assert := assrt.NewAssert(t)
	cwd, _ := filepath.Abs(".")
	nwd, _ := filepath.Abs("..")
	f1, f2, f3, f4, f5, settings := "", "", "", "", "", "[settings]\n"


	//
	//	Default config
	//
	conf := parser().GetConfig()
	assert.Equal(DefaultConfiguration, *conf)


	//
	//	Basic settings
	//
	f1 = `
		# Custom DNS servers
		dns = [ "8.8.8.8" ]

		# What container folder to start in
		folder = "/hroot"

		# Interactive mode
		attach = true

		# Delete the container after running
		purge = true
	`
	conf = parser().
		AddConfig(settings + f1, ".").
		GetConfig()
	expect := DefaultConfiguration //Use default fields with a few exceptions
	expect.Settings.DNS = []string{ "8.8.8.8" }
	expect.Settings.Folder = "/hroot"
	expect.Settings.Attach = true
	expect.Settings.Purge = true
	assert.Equal(expect, *conf)


	//
	//	Mount localizing
	//
	f2 = `
		# Folder mounts
		#	(host folder, container folder, 'ro' or 'rw' permissions)
		mounts = [
			[ ".../", "/boxen",    "ro"],  # The top folder
			[ "./",   "/hroot",   "rw"],  # The current folder
		]
	`
	conf = parser().
		AddConfig(settings + f1 + f2, "..").
		GetConfig()
	expect.Settings.Mounts = [][]string{
		[]string{ nwd, "/boxen",  "ro" },
		[]string{ cwd, "/hroot", "rw" },
	}
	assert.Equal(expect, *conf)


	//
	//	Settings override
	//
	f3 = `
		folder = "/home"
	`
	conf = parser().
		AddConfig(settings + f1 + f2, "..").
		AddConfig(settings + f3,      "." ).
		GetConfig()
	expect.Settings.Folder = "/home"
	assert.Equal(expect, *conf)


	//
	// Image names
	//
	f4 = `
	[image]
		name     = "example.com/ubuntu/12.04"
		upstream = "index.docker.IO/ubuntu/12.04"
	`
	conf = parser().
			AddConfig(settings + f1 + f2, "..").
			AddConfig(settings + f3 + f4, "." ).
			GetConfig()
	expect.Image.Name     = "example.com/ubuntu/12.04"
	expect.Image.Upstream = "index.docker.IO/ubuntu/12.04"
	assert.Equal(expect, *conf)


	//
	//	Target settings override
	//

	f5 = `
	# This is where you specify run targets.
	# Targets let you take different actions with the same container.
	[target.bash]
		command = [ "/bin/bash" ]
		dns = [ "8.8.4.4" ]
	`
	conf = parser().
		AddConfig(settings + f1 + f2,      "..").
		AddConfig(settings + f3 + f4 + f5, "." ).
		GetConfig()
	expect.Settings.DNS = append(expect.Settings.DNS, "8.8.4.4")
	expect.Settings.Command = []string{ "/bin/bash" }
	assert.Equal(1, len(conf.Targets))
	assert.Equal(expect.Settings, conf.Targets["bash"])
}
