package main

import (
	"bufio"
	"database/sql"
	"log"
	"os"
	"strings"

	"github.com/coopernurse/gorp"
	_ "github.com/lib/pq"
)

type Teams struct {
	Id       int64  `db:"id"`
	City     string `db:"city"`
	Nickname string `db:"nickname"`
	Stadium  string `db:"stadium"`
}

func main() {
	db := initDb()

	// Open the file containing team info and wrap it in scanner
	teamData, err := os.Open("teams.txt")
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer teamData.Close()
	scanner := bufio.NewScanner(teamData)

	// Line by line, get the information and load it into the DB
	for scanner.Scan() {
		teamLine := scanner.Text()
		splitLine := strings.Split(teamLine, ",")
		newTeam := &Teams{
			City:     splitLine[0],
			Nickname: splitLine[1],
			Stadium:  splitLine[2],
		}

		err = db.Insert(newTeam)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}
}

func initDb() *gorp.DbMap {
	db, err := sql.Open("postgres", "user=nfl database=nfl_app sslmode=disable")
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
