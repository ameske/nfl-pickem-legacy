package database

import "database/sql"

func teamRecord(city string, nickname string) (wins int, losses int, err error) {
	// TODO: Handle the case where nobody has any games gracefully in the SQL
	s := `SELECT tmp.id, home_wins+away_wins AS wins, home_losses+away_losses AS losses
	  FROM (SELECT home_id AS id, COUNT(*) AS home_wins FROM games WHERE (home_score > away_score) GROUP BY home_id) tmp
	  JOIN (SELECT away_id AS id, COUNT(*) AS away_wins FROM games WHERE (away_score > home_score) GROUP BY away_id) tmp2 ON tmp.id = tmp2.id
	  JOIN (SELECT home_id AS id, COUNT(*) AS home_losses FROM games WHERE (home_score < away_score) GROUP BY home_id) tmp3 ON tmp2.id = tmp3.id
	  JOIN (SELECT away_id AS id, COUNT(*) AS away_losses FROM games WHERE (away_score < home_score) GROUP BY away_id) tmp4 ON tmp3.id = tmp4.id
	  JOIN teams ON tmp.id = teams.id
	  WHERE teams.city = ?1 AND teams.nickname = ?2`

	err = db.QueryRow(s, city, nickname).Scan(&wins, &losses)
	if err == sql.ErrNoRows {
		return 0, 0, nil
	} else if err != nil {
		return -1, -1, err
	}

	return
}
