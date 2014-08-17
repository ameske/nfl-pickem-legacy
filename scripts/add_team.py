'''
install.py

Sets up a brand new instance of the nfl_picks web application.

1. Creates the database.
2. Loads the teams into the databse.

You will need to run the season scripts as you see fit.

@author: Kyle Ames
@date: May 16, 2014
'''

from app import db, models

team_data = open('Scripts/teams.txt', 'r')

# Each line is of the form <CITY,NICKNAME,STADIUM>
for line in team_data:
  line = line.strip()
  split_line = line.split(',')
  new_team = models.Team(city=split_line[0], nickname=split_line[1], stadium=split_line[2])
  db.session.add(new_team)

db.session.commit()
team_data.close()
