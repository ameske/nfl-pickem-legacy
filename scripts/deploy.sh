#!/bin/bash

go run insertTeams.go
go run insertGames.go

go run addUser.go --first Will --last Taylor --email william.taylor.09@cnu.edu --password will
go run addUser.go --first Logen --last Franklin --email logenfranklin@gmail.com --password logen
go run addUser.go --first Leslie --last Jones --email lcjones757@gmail.com --password lesli
go run addUser.go --first Rob --last Specketer --email speckerc@gmail.com --password rob
go run addUser.go --first Chris --last Ames --email amestuxedos@gmail.com --password chris
go run addUser.go --first Kathy --last Ames --email kames@vwc.edu --password Kathy

go run generatePicks -year 2014

go run updateGameScores -week 1
go run updateGameScores -week 2
go run updateGameScores -week 3
go run updateGameScores -week 4
go run updateGameScores -week 5
go run updateGameScores -week 6
go run updateGameScores -week 7
go run updateGameScores -week 8
