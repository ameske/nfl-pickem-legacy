'''
A script to get the results of the previous week's games and score the
user picks of that week. Also creates the pick objects for the upcoming
week. To be run on Tuesday morning.

@author: ameske
@date: 1/11/13
'''
from bs4 import BeautifulSoup
import urllib2
from app import models

year = models.Year.query.filter_by(year=2013).first()
week_types = {1: 'A', 2: 'A', 3:'A', 4:'B', 5:'C', 6:'B', 7:'B', 8:'C', 9:'C', 10:'C', 11:'B', 12:'C', 13:'A', 14:'A', 15:'A', 16:'A', 17:'A' }
schedule_url = 'http://www.nfl.com/schedules/2013/REG'

def get_results(week_no):
    #FETCH THE HTML WITH THE SCORES EMBEDDED
    page = urllib2.urlopen(schedule_url + str(week_no))
    resultsHTML = BeautifulSoup(page.read())
    
    #WE CAN UNIQUELY ID A GAME BASED ON THE HOME TEAM, GATHER THEM AND THE SCORES
    home_teams          = [ str(tag.string) for tag in resultsHTML.find_all(class_="team-name home ") ]
    home_team_scores    = [ str(tag.string) for tag in resultsHTML.find_all(class_="team-score home ") ]
    away_team_scores    = [ str(tag.string) for tag in resultsHTML.find_all(class_="team-score away ") ]
    
    results = zip(home_teams, home_team_scores, away_team_scores)
    
    week = models.Week.query.filter_by(week=weeok_no).first()
    #UPDATE THE GAME SCHEDULE WITH THE APPROPRIATE SCORES
    for result in results:
        home_team = models.Team.query.filter_by(nickname=result[0]).first()
        game = models.Schedule.query.filter_by(week=week,home_team=home_team).first()
        game.home_team_score = int(result[1])
        game.away_team_score = int(result[2])
   
    #COMMIT THE CHANGES TO THE GAME SCHEDULE
    db.session.commit()

    #NOW, SCORE THE USER PICKS, HANDING OVER THE WEEK OBJECT WE JUST RETREIVED
    score_games(week)

    #FINALLY, CREATE THE PICK OBJECTS FOR NEXT WEEK
    create_picks(week_no + 1)


def score_games(week):
    #GATHER ALL OF THE PICKS FOR THE CURRENT WEEK
    weekly_picks = [ pick for pick in [ game.picks for game in week.games ] ]
    
    #FOR EACH PICK, CHECK WHO THE WINNING TEAM IS AND SCORE IT
    for pick in weekly_picks:
        pick.awardedPoints = pick.points if pick.selection == pick.game.winner else 0

    users = models.User.all()
    for current_user in users:
      users_picks = models.User.user_picks_by_week(current_user, week)
      week_statistic = models.Statistic(user=user, year=year, week=week)
      for pick in users_picks:
        if pick.awardedPoints == 7:
          week_statistic.seven = week_statistic.seven + 1
        elif pick.awardedPoints == 5:
          week_statistic.five = week_statistic.five + 1
        elif pick.awardedPoints == 3:
          week_statistic.three = week_statistic.three + 1
        elif pick.awardedPoints == 1:
          week_statistic.one = week_statistic.one + 1

    #COMMIT THE CHANGES TO THE PICKS AND STATISTICS
    db.session.commit()


def create_picks(week_no):
    #GATHER SOME INFORMATION WE WILL NEED
    print week_no
    week = models.Week.query.filter_by(week=week_no).first()
    print week.id

    games = models.Schedule.query.filter_by(week=week).all()
    print "I found %d games." % len(games)

    users = models.User.all()
    print "I found %d users." % len(users)

    #FOR EACH USER, CREATE A PICK FOR EACH GAME IN THE UPCOMING WEEK
    for current_user in users:
        for current_game in games:
          pick = models.Pick(user=current_user, game=current_game)

    #COMMIT THE NEW PICKS
    models.db.session.commit()


def main():
  get_results()
  score_games()
  create_picks()

if __name__ == "__main__":
  main()
