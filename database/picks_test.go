package database

import (
	"fmt"
	"testing"
)

func TestGetWeeklyPicks(t *testing.T) {
	db := NflDb()
	picks := GetWeeklyPicks(db, 1, 2014, 1)

	for _, pick := range picks {
		fmt.Printf("%v", pick)
	}

}
