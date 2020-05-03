.ONESHELL:

.PHONY: postgres

go_get_private:
	go env -w GOPRIVATE=pault.ag
	git config --global url."git@github.com:".insteadOf "https://github.com/"

postgres:
	sudo apt install -y postgresql
	sudo service postgresql start
	sudo -u postgres psql <<POSTGRES_SCRIPT 
	CREATE USER euler WITH PASSWORD 'euler';
	CREATE DATABASE euler;
	POSTGRES_SCRIPT