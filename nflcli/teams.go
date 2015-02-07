package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"log"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)

type Teams struct {
	Id       int64  `db:"id"`
	City     string `db:"city"`
	Nickname string `db:"nickname"`
	Stadium  string `db:"stadium"`
}

func inputTeams(c *cli.Context) {
	scanner := bufio.NewScanner(bytes.NewBufferString(teams))

	// Line by line, get the information and load it into the DB
	for scanner.Scan() {
		teamLine := scanner.Text()
		splitLine := strings.Split(teamLine, ",")
		newTeam := &Teams{
			City:     splitLine[0],
			Nickname: splitLine[1],
			Stadium:  splitLine[2],
		}

		err := db.Insert(newTeam)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("postgres", "host=/run/postgresql user=nfl dbname=nfl_app sslmode=disable")
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = db.Ping()
	if err != nil {
		log.Fatalf(err.Error())
	}

	dbmap := &gorp.DbMap{Db: db, Dialect: gorp.PostgresDialect{}}
	dbmap.AddTableWithName(Teams{}, "teams").SetKeys(true, "Id")

	return dbmap
}

var teams = `Buffalo,Bills,Ralph Wilson Stadium
Miami,Dolphins,Sun Life Stadium
New England,Patriots,Gilette Stadium
New York,Jets,MetLife Stadium
Baltimore,Ravens,M&T Bank Stadium
Cincinatti,Bengals,Paul Brown Stadium
Cleveland,Browns,First Energy Stadium
Pittsburgh,Steelers,Heinz Field
Houston,Texans,Reliant Stadium
Indianapolis,Colts,Lucas Oil Stadium
Jacksonville,Jaguars,EverBank Field
Tennessee,Titans,LP Field
Denver,Broncos,Mile High Stadium
Kansas City,Chiefs,Arrowhead Stadium
Oakland,Raiders,O.co Coliseum
San Diego,Chargers,Qualcomm Stadium
Dallas,Cowboys,AT&T Stadium
New York,Giants,MetLife Stadium
Philadelphia,Eagles,Lincoln Financial Field
Washington,Redskins,FedEx Field
Chicago,Bears,Soldier Field
Detroit,Lions,Ford Field
Green Bay,Packers,Lambeau Field
Minnesota,Vikings,TCF Bank Stadium
Atlanta,Falcons,Georiga Dome
Carolina,Panthers,Bank of America Stadium
New Orleans,Saints,Mercedes-Benz Superdome
Tampa Bay,Buccaneers,Raymond James Stadium
Arizona,Cardinals,University of Phoenix Stadium
St. Louis,Rams,Edward Jones Dome
San Francisco,49ers,Candlestick Park
Seattle,Seahawks,CenturyLink Field`
