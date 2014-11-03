package database

func UsersMap(users []Users) map[int64]Users {
	um := make(map[int64]Users)
	for _, u := range users {
		u := u
		um[u.Id] = u
	}

	return um
}

func GamesMap(games []Games) map[int64]Games {
	gm := make(map[int64]Games)
	for _, g := range games {
		g := g
		gm[g.Id] = g
	}

	return gm
}
