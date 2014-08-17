'''
weeklyScores.py

Scrapes the NFL's website to pull the reults from the given year and week's games

@author: ameske
@date: 8/17/14
'''
from bs4 import BeautifulSoup
import urllib2

schedule_url = 'http://www.nfl.com/schedules/2014/REG'

def get_results(year, week_no):
    games = []

    #FETCH THE HTML WITH THE SCORES EMBEDDED
    page = urllib2.urlopen(schedule_url + str(week_no))
    resultsHTML = BeautifulSoup(page.read())
    
    #WE CAN UNIQUELY ID A GAME BASED ON THE HOME TEAM, GATHER THEM AND THE SCORES
    home_teams          = [ str(tag.string) for tag in resultsHTML.find_all(class_="team-name home ") ]
    home_team_scores    = [ str(tag.string) for tag in resultsHTML.find_all(class_="team-score home ") ]
    away_team_scores    = [ str(tag.string) for tag in resultsHTML.find_all(class_="team-score away ") ]
    
    results = zip(home_teams, home_team_scores, away_team_scores)
    
    for result in results:
        d = {
                "year": year,
                "week": week_no,
                "home": result[0],
                "home_score": result[1],
                "away_score": result[0],
        }                
        games.append(d)
                
    return games

def main():
  results = get_results(int(year), int(week)))
  fd = open(year + '-Week' + week + '_Results.json')
  fd.write(json.dumps(results))
  fd.close()

if __name__ == "__main__":
  main()
