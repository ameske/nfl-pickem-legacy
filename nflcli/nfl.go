package main

import (
	"os"

	"github.com/ameske/nfl-pickem/database"
	"github.com/codegangsta/cli"
	"github.com/coopernurse/gorp"
)

var (
	db     *gorp.DbMap
	config Config
)

type Config struct {
	StaticDir string `yaml:"STATIC_DIR"`
}

func init() {
	/*
		config := Config{}
		configBytes, err := ioutil.ReadFile("/opt/ameske/etc/go_nfl/nfl.yaml")
		if err != nil {
			log.Fatalf("Error reading config file: %s", err.Error())
		}
		err = yaml.Unmarshal(configBytes, &config)
		if err != nil {
			log.Fatalf("Error parsing config file: %s", err.Error())
		}
	*/
	db = database.NflDb("host=localhost port=5432")
}

func main() {
	app := cli.NewApp()
	app.Name = "nfl"
	app.Version = "0.0.1"
	app.Usage = "manage the go_nfl web application"
	app.Action = func(c *cli.Context) {
		cli.ShowAppHelp(c)
	}
	app.Commands = []cli.Command{
		{
			Name:      "add",
			ShortName: "a",
			Usage:     "manual management of database information",
			Subcommands: []cli.Command{
				{
					Name:  "user",
					Usage: "add a user to the database",
					Flags: []cli.Flag{
						cli.StringFlag{
							Name:  "first",
							Value: "",
							Usage: "first name",
						},
						cli.StringFlag{
							Name:  "last",
							Value: "",
							Usage: "last name",
						},
						cli.StringFlag{
							Name:  "email",
							Value: "",
							Usage: "email/username",
						},
						cli.StringFlag{
							Name:  "password",
							Value: "",
							Usage: "password",
						},
					},
					Action: inputUser,
				},
				{
					Name:  "picks",
					Usage: "add picks to the database",
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "year",
							Value: -1,
							Usage: "year to input picks",
						},
						cli.IntFlag{
							Name:  "week",
							Value: -1,
							Usage: "week to input picks",
						},
						cli.StringFlag{
							Name:  "user",
							Value: "",
							Usage: "user to input picks for",
						},
					},
					Action: inputPicks,
				},
				{
					Name:   "teams",
					Usage:  "add teams to the database",
					Action: inputTeams,
				},
			},
		},
		{
			Name:      "import",
			ShortName: "i",
			Usage:     "options for scraping info from nfl website",
			Subcommands: []cli.Command{
				{
					Name:   "schedule",
					Usage:  "pulls schedule information from the nfl",
					Action: schedule,
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "year, y",
							Value: -1,
							Usage: "year to pull schedule",
						},
						cli.IntFlag{
							Name:  "week, w",
							Value: -1,
							Usage: "week to pull schedule",
						},
					},
				},
				{
					Name:   "scores",
					Usage:  "pulls completed game scores from the nfl",
					Action: scores,
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "year, y",
							Value: -1,
							Usage: "year to pull scores",
						},
						cli.IntFlag{
							Name:  "week, w",
							Value: -1,
							Usage: "week to pull scores",
						},
					},
				},
			},
		},
		{
			Name:      "grade",
			ShortName: "gr",
			Usage:     "grade users picks for the specified week",
			Action:    grade,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "year, y",
					Value: -1,
					Usage: "year to grade picks",
				},
				cli.IntFlag{
					Name:  "week, w",
					Value: -1,
					Usage: "week to grade picks",
				},
			},
		},
		{
			Name:      "html",
			ShortName: "ge",
			Usage:     "generate static html pages for the website",
			Subcommands: []cli.Command{
				{
					Name:   "results",
					Usage:  "generates the results page",
					Action: results,
					Flags: []cli.Flag{
						cli.IntFlag{
							Name:  "year, y",
							Value: -1,
							Usage: "year to generate results",
						},
						cli.IntFlag{
							Name:  "week, w",
							Value: -1,
							Usage: "week to generate results",
						},
					},
				},
			},
		},
		{
			Name:      "testdata",
			ShortName: "test",
			Usage:     "generate test data for development or unit tests",
			Action:    generateTestData,
			Flags: []cli.Flag{
				cli.IntFlag{
					Name:  "weeks, w",
					Value: 1,
					Usage: "number of weeks of data to create",
				},
				cli.IntFlag{
					Name:  "users, u",
					Value: 1,
					Usage: "number of test users to create",
				},
				cli.BoolFlag{
					Name:  "fakepicks",
					Usage: "create fake picks for the test users",
				},
				cli.BoolFlag{
					Name:  "fakescores",
					Usage: "create fake scores for the games",
				},
			},
		},
	}

	app.Run(os.Args)
}
