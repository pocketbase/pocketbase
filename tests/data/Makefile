
# Dependencies: pg-loader, 

install-pgloader:
	brew install pgloader

run-docker:
	docker run -d --rm --name postgres -e POSTGRES_USER=user -e POSTGRES_PASSWORD=pass -p 5432:5432 postgres:alpine

dump-db:
# 1. Convert `data.db` to `data.sqlite-dump.sql` and `data.pg-dump.sql`
	sqlite3 data.db .dump > data.sqlite-dump.sql
	-PGPASSWORD=pass dropdb -h 127.0.0.1 -U user pb-data
	PGPASSWORD=pass createdb -h 127.0.0.1 -U user pb-data
	-pgloader --with "preserve index names" sqlite://$$(pwd)/data.db pgsql://user:pass@127.0.0.1/pb-data
	docker exec -it postgres pg_dump -U user -h 127.0.0.1 pb-data > data.pg-dump.sql
# 2. Convert `auxiliary.db` to `auxiliary.sqlite-dump.sql` and `auxiliary.pg-dump.sql`
	sqlite3 auxiliary.db .dump > auxiliary.sqlite-dump.sql
	-PGPASSWORD=pass dropdb -h 127.0.0.1 -U user pb-auxiliary
	PGPASSWORD=pass createdb -h 127.0.0.1 -U user pb-auxiliary
	-pgloader --with "preserve index names" sqlite://$$(pwd)/auxiliary.db pgsql://user:pass@127.0.0.1/pb-auxiliary
	docker exec -it postgres pg_dump -U user -h 127.0.0.1 pb-auxiliary > auxiliary.pg-dump.sql


PREFIX=pb_test_
drop-matching-dbs:
	@echo "Looking for databases starting with $(PREFIX)..."
	@PGPASSWORD=pass psql -h 127.0.0.1 -U user -d postgres -t -c "SELECT datname FROM pg_database WHERE datname LIKE '$(PREFIX)%'" | \
	while read dbname; do \
	  dbname=$$(echo $$dbname | xargs); \
	  if [ ! -z "$$dbname" ]; then \
	    echo "Dropping database: $$dbname"; \
	    PGPASSWORD=pass psql -h 127.0.0.1 -U user -d postgres -c "DROP DATABASE \"$$dbname\""; \
	  fi; \
	done
