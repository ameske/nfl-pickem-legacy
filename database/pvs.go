package database

import (
	"log"

	"github.com/coopernurse/gorp"
)

type Pvs struct {
	Id    int64  `db:"id"`
	Type  string `db:"type"`
	Seven int    `db:"seven"`
	Five  int    `db:"five"`
	Three int    `db:"three"`
	One   int    `db:"one"`
}

func GetPvs(db *gorp.DbMap, weekId int64) Pvs {
	var pvs Pvs
	err := db.SelectOne(&pvs, "SELECT pvs.* FROM weeks JOIN pvs ON weeks.pvs_id = pvs.id WHERE weeks.id = $1", weekId)
	if err != nil {
		log.Fatalf("GetPvs: %s", err.Error())
	}

	return pvs
}
