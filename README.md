# humus-waste-watcher
A Telegram bot designed for fun and stats! Track daily, monthly, and yearly "activity," generate personalized overviews with a "Wrapped" feature, and celebrate the champions with a monthly podium. Built for laughs, competition, and quirky data tracking with friends.

# How to run
This bot is deployed using Fly.io and any concurrent local run will stop both the local and deployed runs. However, if the deployed bot is not running, you can run the bot by following these steps:
- Install [Go](https://go.dev/)
- From the root directory, run `go run main/bot.go`
<br/><br/>

# Releases

## 1.0.0
This release consists of the bot deployment to Fly.io with the following features:
- Check current month's and global waste counter
- Check monthly overview of the waste counter
- Check monthly podium
- Check yearly podium
- Announcement of the monthly podium at the end of each month
- Announcement of the yearly podium at the end of each year

## 1.1.0
This release adds a few more features and reduces some costs on Fly.io. Features added:
- Bot reacts to messages
- Check the monthly leaderboard
- Check the bottom three from the leaderboard