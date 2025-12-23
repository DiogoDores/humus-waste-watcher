.PHONY: setup-dev

setup-dev:
	flyctl ssh sftp get /app/data/poop_tracker.db ./data/poop_tracker.db
	flyctl machine stop 7849945cee4248
	go run main/bot.go

get-db:
	flyctl ssh sftp get /app/data/poop_tracker.db ./data/poop_tracker.db

run:
	go run main/bot.go

stop:
	flyctl machine stop 7849945cee4248

start:
	flyctl machine start 7849945cee4248

deploy:
	flyctl deploy