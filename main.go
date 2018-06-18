package main

import (
	"io/ioutil"
	"os"
	"strings"

	"github.com/jmccann/drone-github-status/plugin"

	"github.com/Sirupsen/logrus"
	"github.com/urfave/cli"
)

var revision string // build number set at compile-time

func main() {
	app := cli.NewApp()
	app.Name = "github status plugin"
	app.Usage = "github status plugin"
	app.Action = run
	app.Version = revision
	app.Flags = []cli.Flag{

		//
		// plugin args
		//

		cli.StringFlag{
			Name:   "api-key",
			Usage:  "api key to access github api",
			EnvVar: "PLUGIN_API_KEY,GITHUB_RELEASE_API_KEY,GITHUB_TOKEN",
		},
		cli.StringFlag{
			Name:   "base-url",
			Value:  "https://api.github.com/",
			Usage:  "api url, needs to be changed for ghe",
			EnvVar: "PLUGIN_BASE_URL,GITHUB_BASE_URL",
		},
		cli.StringSliceFlag{
			Name:   "context",
			Usage:  "status context(s) to create/update",
			EnvVar: "PLUGIN_CONTEXT,PLUGIN_CONTEXTS",
		},
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug logging",
			EnvVar: "PLUGIN_DEBUG",
		},
		cli.StringFlag{
			Name:   "description",
			Usage:  "status description",
			EnvVar: "PLUGIN_DESCRIPTION",
		},
		cli.StringFlag{
			Name:   "file",
			Usage:  "status read from file",
			EnvVar: "PLUGIN_FILE",
		},
		cli.StringFlag{
			Name:   "state",
			Usage:  "status state",
			EnvVar: "PLUGIN_STATE",
		},
		cli.StringFlag{
			Name:   "target-url",
			Usage:  "url to have status link to",
			EnvVar: "PLUGIN_TARGET_URL,DRONE_BUILD_LINK",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "basic auth password",
			EnvVar: "PLUGIN_PASSWORD,GITHUB_PASSWORD,DRONE_NETRC_PASSWORD",
		},
		cli.StringFlag{
			Name:   "username",
			Usage:  "basic auth username",
			EnvVar: "PLUGIN_USERNAME,GITHUB_GITHUB_USERNAME,DRONE_NETRC_USERNAME",
		},

		//
		// drone env
		//

		cli.StringFlag{
			Name:   "build-status",
			Usage:  "drone build status, can be used to derive state from",
			EnvVar: "DRONE_BUILD_STATUS",
		},
		cli.StringFlag{
			Name:   "commit-sha",
			Usage:  "commit-sha to assign status to",
			EnvVar: "DRONE_COMMIT_SHA",
		},
		cli.StringFlag{
			Name:   "repo-name",
			Usage:  "repository name",
			EnvVar: "DRONE_REPO_NAME",
		},
		cli.StringFlag{
			Name:   "repo-owner",
			Usage:  "repository owner",
			EnvVar: "DRONE_REPO_OWNER",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	if c.Bool("debug") {
		logrus.SetLevel(logrus.DebugLevel)
	}

	logrus.WithFields(logrus.Fields{
		"Revision": revision,
	}).Info("Drone Github Status Plugin Version")

	// Read state from file
	state := c.String("state")
	logrus.WithField("state", state).Debug("user defined state")
	if state == "" && c.String("file") != "" {
		if _, err := os.Stat(c.String("file")); err == nil {
			dat, err := ioutil.ReadFile(c.String("file"))

			if err != nil {
				logrus.Errorf("failed to read state from file %s", c.String("file"))
				state = "error"
			} else {
				state = strings.TrimSuffix(string(dat), "\n")
				logrus.WithField("state", state).Debug("state provided from file")
			}
		}
	}

	p, err := plugin.NewFromCLI(c)
	if err != nil {
		return err
	}
	p.State = state

	return p.Exec()
}
