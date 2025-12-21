.PHONY: setup-dev

setup-dev:
	powershell -Command "if (Test-Path './data/poop_tracker.db') { Remove-Item './data/poop_tracker.db' }"
	flyctl ssh sftp get /app/data/poop_tracker.db ./data/poop_tracker.db
	flyctl machine stop 683d527a22d298 
	go run main/bot.go