'''
weeklyScores.py

Scrapes the NFL's website to pull the reults from the given year and week's games

@author: ameske
@date: 10/5/14
'''

from bs4 import BeautifulSoup
import urllib2
import json
import sys

schedule_url = 'http://www.nfl.com/schedules/2014/REG'
year = -1
week = -1

def home_team(tag):
    return tag.has_attr('class') and "team-name" in tag['class'] and "home" in tag['class']

def away_team(tag):
    return tag.has_attr('class') and "team-name" in tag['class'] and "away" in tag['class']

def home_team_score(tag):
    return tag.has_attr('class') and "team-score" in tag['class'] and "home" in tag['class']

def away_team_score(tag):
    return tag.has_attr('class') and "team-score" in tag['class'] and "away" in tag['class']

def get_results(year, week_no):
    #FETCH THE HTML WITH THE SCORES EMBEDDED
    page = urllib2.urlopen(schedule_url + str(week_no))
    resultsHTML = BeautifulSoup(page.read())
    
    #WE CAN UNIQUELY ID A GAME BASED ON THE HOME TEAM, GATHER THEM AND THE SCORES
    home_teams          = [ str(tag.string) for tag in resultsHTML.find_all(home_team) ]
    home_team_scores    = [ str(tag.string) for tag in resultsHTML.find_all(home_team_score) ]
    away_team_scores    = [ str(tag.string) for tag in resultsHTML.find_all(away_team_score) ]
    results = zip(home_teams, home_team_scores, away_team_scores)
    
    #JSON'IFY THE RESULTS
    games = []
    for result in results:
        d = {
                "year": year,
                "week": week_no,
                "home": result[0],
                "home_score": int(result[1]),
                "away_score": int(result[2]),
        }                
        games.append(d)
                
    return games

def main():
  year = int(sys.argv[1])
  week = int(sys.argv[2])
  results = get_results(int(year), int(week))
  fd = open('{0}-Week{1}-Results.json'.format(year, week), 'w')
  fd.write(json.dumps(results, indent=4, separators=(',', ': ')))
  fd.close()

if __name__ == "__main__":
  main()
