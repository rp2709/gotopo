package viewsheds_generation

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/rp2709/gotopo/gotopo"
)

const viewshedSizeBits = 128
const minimumAngleFromHorizon = math.Pi / 12
const minimumAngleFromZenith = math.Pi / 12

type command int

const (
	printDirections command = iota
	visualize
	createWorkAreas
	generate
)

var commandFromName = map[string]command{
	"print-directions":  printDirections,
	"visualize":         visualize,
	"generate":          generate,
	"create-work-areas": createWorkAreas,
}

var commandDescription = map[command]string{
	printDirections: "Computes all the direction vectors used to generate viewsheds and prints them as a list of vectors to stdout in this format: size(x,y,z),...,(x,y,z),",
	visualize:       "Generate one viewshed at specified coordinates and produces an image",
	createWorkAreas: "Divides the extent of all samples in manageable work areas.",
	generate:        "Generate a hemispherical viewshed of arbitrary precision every sample point using a digital surface model",
}

type commandFunction func(gotopo.DatabaseConfig, viewshedsConfigT, []string) error

var commandEntrypoint = map[command]commandFunction{
	printDirections: doPrintDirections,
	visualize:       doVisualizeViewsheds,
	createWorkAreas: doCreateWorkAreas,
	generate:        doGenerateViewsheds,
}

func Run(databaseConfig gotopo.DatabaseConfig, args []string) error {
	// viewshed-generation specific flags
	var viewshedConfig viewshedsConfigT
	flags := flag.NewFlagSet("viewsheds", flag.ExitOnError)
	flags.Usage = usage(flags)
	flags.StringVar(&viewshedConfig.VrtFilePath, "vrt", "", "Path to .vrt file indexing the surface model's tiles")
	flags.Float64Var(&viewshedConfig.MaxDistance, "max-distance", 40000, "maximum distance to explore from each sample point")

	err := flags.Parse(args)
	if err != nil {
		return err
	}

	// read subcommand
	if flags.NArg() < 1 {
		flags.Usage()
		os.Exit(1)
	}
	cmd, ok := commandFromName[flags.Arg(0)]
	if !ok {
		flags.Usage()
		os.Exit(2)
	}

	commandArgs := flags.Args()[1:]
	return commandEntrypoint[cmd](databaseConfig, viewshedConfig, commandArgs)
}

func usage(viewshedsFlags *flag.FlagSet) func() {
	return func() {
		fmt.Fprintln(os.Stderr, "Usage: gotopo [options] viewsheds [options] <subcommand>")
		fmt.Fprintln(os.Stderr, "Available commands:")
		for name, cmd := range commandFromName {
			fmt.Fprintf(os.Stderr, "  - %s : %s\n", name, commandDescription[cmd])
		}
		fmt.Fprintln(os.Stderr, "Options:")
		viewshedsFlags.PrintDefaults()
	}
}

type viewshedsConfigT struct {
	VrtFilePath string
	MaxDistance float64
}
