build:
	docker build -t forum .
run:
	docker run -p 8080:8080 forum
initdb:
	sqlite3 forum.sqlite3 < migrations/init.sql