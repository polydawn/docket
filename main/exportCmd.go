package main

import (
	// "github.com/jessevdk/go-flags"
	"polydawn.net/docket/confl"
	"polydawn.net/docket/crocker"
)

type exportCmdOpts struct{}

func (opts *exportCmdOpts) Execute(args []string) error {
	//Get the target
	target := ""
	if len(args) == 1 {
		target = args[0]
	} else {
		target = "build"
	}

	return WithDocker(func(dock *crocker.Dock, settings *confl.ConfigLoad) error {
		return Export(dock, settings, target)
	})
}

const ExportPath = "./image.tar" //Where to export docker images

//Exports the result of a target into docker.
func Export(dock *crocker.Dock, settings *confl.ConfigLoad, target string) error {
	//Get configuration
	config := settings.GetConfig(target)
	saveAs := settings.GetConfig(confl.DefaultTarget).Image

	//Run the build
	container := Launch(dock, config)
	container.Wait()

	//Create a tar
	container.ExportToFilename(ExportPath)

	//Import the built docker
	// Todo: add --noImport option to goflags
	dock.ImportFromFilenameTagstring(ExportPath, saveAs)

	//Remove if desired
	if config.Purge {
		container.Purge()
	}

	return nil
}