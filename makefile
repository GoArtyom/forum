build:
	docker build -t forum .
run:
	docker run -p 8081:8081 forum
initdb:
	sqlite3 forum.sqlite3 < migrations/init.sql