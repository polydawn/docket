package commands

import (
	. "fmt"
)

type BuildCmdOpts struct {
	DockerH     string `short:"H"                    description:"Where to connect to docker daemon."`
	Source      string `short:"s" long:"source"      description:"Container source.      (default: graph)"`
	Destination string `short:"d" long:"destination" description:"Container destination. (default: graph)"`
	NoOp        bool   `long:"noop" description:"Set the container command to /bin/true."`
	Epoch       bool   `long:"epoch" description:"Force all file modtimes to epoch."`
}

const DefaultBuildTarget = "build"

//Transforms a container
func (opts *BuildCmdOpts) Execute(args []string) error {
	//Load settings
	hroot := LoadHroot(args, DefaultBuildTarget, opts.Source, opts.Destination)

	//We're building; launch upstream image
	hroot.launchImage = hroot.image.Upstream
	Println("Building from", hroot.image.Upstream, "to", hroot.image.Name)

	//If desired, set the command to /bin/true and do not modify destination image name
	//We'd love to not launch the container at all, but docker's export is completely broken.
	// 'docker export ubuntu' --> 'Error: No such container: ubuntu' --> :(
	if opts.NoOp {
		hroot.settings.Command = []string{ "/bin/true" }
	}

	//Prepare source & destination
	hroot.PrepareInput()
	hroot.PrepareOutput()

	//Start or connect to a docker daemon
	hroot.StartDocker(opts.DockerH)
	hroot.PrepareCache()
	hroot.Launch()

	//Perform any destination operations required
	hroot.ExportBuild(opts.Epoch)

	hroot.Cleanup()
	return nil
}
