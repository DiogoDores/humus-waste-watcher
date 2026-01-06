.PHONY: setup-dev

setup-dev:
	flyctl ssh sftp get /data/poop_tracker.db ./data/poop_tracker.db
	flyctl machine stop 7849945cee4248
	go run main/bot.go

get-db:
	flyctl ssh sftp get /data/poop_tracker.db ./data/poop_tracker.db

put-db:
	flyctl ssh sftp put ./data/poop_tracker.db /data/poop_tracker.db

delete-db:
	flyctl ssh console -a humus-waste-watcher
	rm /data/poop_tracker.db

run:
	go run main/bot.go

stop:
	flyctl machine stop 7849945cee4248

start:
	flyctl machine start 7849945cee4248

deploy:
	flyctl deploy

test:
	go test ./repository/... -v

test-coverage:
	go test ./repository/... -cover

test-benchmark:
	go test ./repository/... -bench=. -benchmem

generate-all-wrapped:
	go run main/wrapped/generate_all_wrapped.go