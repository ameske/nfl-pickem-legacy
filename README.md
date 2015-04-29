go_nfl
======

A customized NFL Pick-Em Pool manager written in Go.

There weren't any pick-em managers that gave the ability to set up customized scoring. This is an automated
pick-em-pool application that uses our scoring system. 

The system works as follows:

1. Each week, pick the winners of all of the games.
2. Next, assign point values to the games. You have one 7 point game, two 5 point games, and anywhere from
three to five 3 point games based on the number of byes in a week. The remainder of the games
are worth only one point.

There is currently support for:

- Importing schedules for a season
- Importing the results on a weekly basis
- Generating HTML to display a week's picks/results
- CLI app for managing the database
- Login/Logout and session management
- E-mail notifications to users when picks are submitted/modified
- Server side verification of picks (client side is in the works)
