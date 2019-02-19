package cmd

import (
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli"

	"gitlab.inn4science.com/gophers/forge/configs"
	"gitlab.inn4science.com/gophers/forge/generate"
)

const (
	debugFlag      = "debug"
	devFlag        = "dev"
	nomemcopyFlag  = "nomemcopy"
	nocompressFlag = "nocompress"
	nometadataFlag = "nometadata"
	tagsFlag       = "tags"
	pkgFlag        = "pkg"
	outputFlag     = "o"
	modeFlag       = "mode"
	modetimeFlag   = "modetime"
	ignoreFlag     = "ignore"
	inputFlag      = "i"
)

var BindataCmd = cli.Command{
	Name:  "bindata",
	Usage: "forge bindata <options>",
	Flags: []cli.Flag{
		cli.BoolFlag{
			Name: debugFlag,
			Usage: "Do not embed the assets, but provide the embedding API. " +
				"Contents will still be loaded from disk.",
		},
		cli.BoolFlag{
			Name: devFlag,
			Usage: "Similar to debug, but does not emit absolute paths. " +
				"Expects a rootDir variable to already exist in the generated code's package.",
		},
		cli.BoolFlag{
			Name: nomemcopyFlag,
			Usage: "Use a .rodata hack to get rid of unnecessary memcopies. " +
				"Refer to the documentation to see what implications this carries.",
		},
		cli.BoolFlag{
			Name:  nocompressFlag,
			Usage: "Assets will *not* be GZIP compressed when this flag is specified.",
		},
		cli.BoolFlag{
			Name:  nometadataFlag,
			Usage: "Assets will not preserve size, mode, and modtime info.",
		},
		cli.StringFlag{
			Name:  tagsFlag,
			Usage: "Optional set of build tags to include.",
		},
		cli.StringFlag{
			Name:  prefixFlag,
			Usage: "Optional path prefix to strip off asset names.",
		},
		cli.StringFlag{
			Name:  pkgFlag,
			Usage: "Package name to use in the generated code.",
		},
		cli.StringFlag{
			Name:  outputFlag,
			Usage: "Optional name of the output file to be generated.",
		},
		cli.UintFlag{
			Name:  modeFlag,
			Usage: "Optional file mode override for all files.",
		},
		cli.Int64Flag{
			Name:  modetimeFlag,
			Usage: "Optional modification unix timestamp override for all files.",
		},
		cli.StringSliceFlag{
			Name:  ignoreFlag,
			Usage: "Regex pattern to ignore",
		},
		cli.StringSliceFlag{
			Name:  inputFlag,
			Usage: "List of input directories/files",
		},
	},
	Action: bindataAction,
}

func bindataAction(c *cli.Context) error {
	cfg := bindataConfig(c)
	if err := cfg.Validate(); err != nil {
		return cli.NewExitError("[ERROR] "+err.Error(), 1)
	}
	if err := generate.Bindata(cfg); err != nil {
		return cli.NewExitError("[ERROR] "+err.Error(), 1)
	}
	return nil
}

func bindataConfig(c *cli.Context) *configs.BindataConfig {
	cfg := &configs.BindataConfig{
		Debug:      c.Bool(debugFlag),
		Dev:        c.Bool(devFlag),
		NoMemCopy:  c.Bool(nomemcopyFlag),
		NoCompress: c.Bool(nocompressFlag),
		NoMetadata: c.Bool(nometadataFlag),
		Tags:       c.String(tagsFlag),
		Prefix:     c.String(prefixFlag),
		Package:    c.String(pkgFlag),
		Output:     c.String(outputFlag),
		Mode:       c.Uint(modeFlag),
		ModTime:    c.Int64(modetimeFlag),
	}
	if cfg.Output == "" {
		cfg.Output = "./bindata.go"
	}
	if cfg.Package == "" {
		cfg.Package = "main"
	}
	ignore := c.StringSlice(ignoreFlag)
	cfg.Ignore = make([]*regexp.Regexp, 0)
	for _, pattern := range ignore {
		cfg.Ignore = append(cfg.Ignore, regexp.MustCompile(pattern))
	}
	input := c.StringSlice(inputFlag)
	cfg.Input = make([]configs.BindataInputConfig, 0)
	for i := range input {
		cfg.Input = append(cfg.Input, parseInput(input[i]))
	}
	return cfg
}

// parseRecursive determines whether the given path has a recursive indicator and
// returns a new path with the recursive indicator chopped off if it does.
//      /path/to/foo/...    -> (/path/to/foo, true)
//      /path/to/bar        -> (/path/to/bar, false)
func parseInput(path string) configs.BindataInputConfig {
	if strings.HasSuffix(path, "/...") {
		return configs.BindataInputConfig{
			Path:      filepath.Clean(path[:len(path)-4]),
			Recursive: true,
		}
	}
	return configs.BindataInputConfig{
		Path:      filepath.Clean(path),
		Recursive: false,
	}
}
