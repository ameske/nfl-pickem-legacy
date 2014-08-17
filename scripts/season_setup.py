'''
season_setup.py

Parses the DOM of www.nfl.com/schedules to load the database with the
matchups for the current season.

Work required each season:
  1. Update 'schedule_url' to the current root URL for the NFL's schedules.
  2. Update the 'week_types' to account for variability in the NFL schedule.
  3. You probably should also make sure the DOM hasn't changed. The NFL loves to do that.

@author: Kyle Ames
@date: May 16, 2014
'''
from bs4 import BeautifulSoup
import urllib2
import datetime
import re

week_types = {1: 'A', 2: 'A', 3:'A', 4:'B', 5:'C', 6:'B', 7:'B', 8:'C', 9:'C', 10:'C', 11:'B', 12:'C', 13:'A', 14:'A', 15:'A', 16:'A', 17:'A' }
schedule_url = 'http://www.nfl.com/schedules/2013/REG'

def is_away_team(css_class):
  return css_class == "team-name away lost" or css_class == "team-name away "

def is_home_team(css_class):
  return css_class == "team-name home lost" or css_class == "team-name home "

def load_schedule():
    for week_no in range(17, 18):
        
        # OPEN UP THE HTML FOR EACH WEEK
#        print 'Current Week: ' + str(week_no) 
#        page = urllib2.urlopen(schedule_url + str(week_no))
#        scheduleHTML = BeautifulSoup(page.read())

        page = open('week1.html')
        scheduleHTML = BeautifulSoup(page.read())
        
        # PARALLEL LISTS OF AWAY, HOME, TIME
        away_teams  = [ str(tag.string) for tag in scheduleHTML.find_all(class_="team-name away ") ]
        print len(away_teams)
        home_teams  = [ str(tag.string) for tag in scheduleHTML.find_all(class_="team-name home ") ]
        print len(home_teams)
        times       = [ str(tag.string) for tag in scheduleHTML.find_all(class_="time") ]
        dates       = [ str(tag.find_previous(class_="schedules-list-date").contents[1].strings.next()).replace(',', '') for tag in scheduleHTML.find_all(class_="list-matchup-row-team") ]

        #  MERGE times and dates AND THEN BRING BACK INTO A DATETIME LIST
        dates_times = zip(dates, times)
        dates       = [ dt[0] + " 2013 " + dt[1] + " PM EST" for dt in dates_times]
        datetimes = [ datetime.datetime.strptime(date_string, '%A %B %d %Y %I:%M %p %Z') for date_string in dates ]
        
        # CREATE MATCHUP TUPLES AND THEN SEND THEM TO THE LOADING FUNCTION
        matchups = zip(away_teams, home_teams, datetimes)      
        for matchup in matchups:
          print matchup            

def main():
    load_schedule()
    
if __name__=='__main__':
    main()
    
