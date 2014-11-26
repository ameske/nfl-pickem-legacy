go_nfl
======

A customized NFL Pick-Em Pool manager written in Go.

There weren't any pick-em managers that allowed one to set up customized scoring like the pool
that I run does. This is an automated pick-em-pool application that uses our scoring system.
The system works as follows. Each week, pick the winners of all of the games. Next, assign
point values to the games. You have one 7 point game, two 5 point games, and anywhere from
three to five 3 point games based on the number of byes in a week. The remainder of the games
are worth only one point.

This is a work in progress in my spare time, so I cannot make any guarantees to when I will
complete it. You are more than welcome to use any part of the application that you like. I will try to keep an accurate report on what I'm working on and what I have completed, updating every few weeks with the high-level changes and future plans.

# Project Status Updates

Current Development:

- Pick form processing

November 25, 2014:

- Results pages for each week
- Standings pages for each week
- Changed the home page to always show the current standings
- A script to generate the static results page for a given week
- 
November 7, 2014: 

- Login and Logout
- Ability for users to change their password
- A variety of scripts to administer the application. (adding users, manually inputting picks, grading picks)
- A "make your picks" page (note: back-end processing of this form is not yet complete)
- Redirection to desired endpoints following login when accessing a protected endpoint without authorization

October 22, 2014:
- Auth and protected endpoint skeleton code
- Scripts to scrape game information from the NFL's website (python)
- Scripts to scrape game results from the NFL's website (python)
- Scripts to manage/setup a new installation of the project (go)
  - SQL DDL
  - Add new users
  - Import games from JSON generated by my game schedule scraper
  - Import results from JSON generated by my game result scraper
  - Create picks for a season
- Table definitions for use with [gorp](https://github.com/coopernurse/gorp)
