go_nfl
======

NOTE: This is my first attempt at making the NFL pickem' pool. It's a completely server side application. I have abandoned maintaining this version of the application in favor of splitting the app into a RESTful API and a completely independent front end. Since the architecture is completely distinct, and the git history is not relevant any more to the new version, I have decided to split the respositories rather than maintain a legacy branch. www.github.com/ameske/nfl-pickem will refer to the new version of the application.

A customized NFL Pick-Em Pool manager written in Go.

There weren't any pick-em managers that gave the ability to set up customized scoring. This is an automated
pick-em-pool application that uses our scoring system. 

The system works as follows:

1. Each week, pick the winners of all of the games.
2. Next, assign point values to the games. You have one 7 point game, two 5 point games, and five 3
point games based on the number of byes in a week. The remainder of the games
are worth only one point.

There is currently support for:

- Login/Logout and session management
- Admin picks management page (set/correct/update individual picks for all users for a given week)
- Automatic importing of game scores after each "wave" of games completes
- Automatic grading of picks after each "wave" of games completes
- Dynamically generated results page for each week
- Dynamically generated standings page for each week
- Dynamically generated picks form based on the current NFL week and whether games have already started or not
- E-mail notifications to users when picks are submitted/modified
- CLI app for manually managing score imports and grading
