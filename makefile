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

####### IRC stuff #########
ifeq ($(findstring arm6,$(shell uname -m)), arm6)
IRC_ARCH = arm6
else ifeq ($(findstring arm7,$(shell uname -m)), arm7)
IRC_ARCH = arm7
else ifeq ($(findstring x86,$(shell uname -m)), x86)
IRC_ARCH = x64
else
$(error unable to support arch for irc)
endif

ifeq ($(findstring Linux,$(shell uname -s)), Linux)
IRC_OS = linux
else ifeq ($(findstring Darwin,$(shell uname -s)), Darwin)
IRC_OS = darwin
else
$(error unable to determine OS)
endif

ircd:
	$(info os: $(IRC_OS) arch: $(IRC_OS))
	wget -O /tmp/oragono.tar.gz https://github.com/oragono/oragono/releases/download/v2.0.0/oragono-2.0.0-$(IRC_OS)-$(IRC_ARCH).tar.gz
ifneq ($(whoami), root)
	sudo adduser --system --group oragono
	sudo tar -xf /tmp/oragono.tar.gz -C /home/oragono
	sudo cp support_files/ircd.yaml /home/oragono/ircd.yaml
	sudo cp $(shell ls -d /home/oragono/oragono*)/oragono /home/oragono/oragono
	sudo cp support_files/oragono.service /etc/systemd/system/oragono.service
	sudo systemctl daemon-reload
	sudo systemctl enable ircd.service
	sudo systemctl start ircd.service
else 
	adduser --system --group oragono
	tar -xf /tmp/oragono.tar.gz -C /home/oragono
endif
ifneq ($(whoami), root)
	
else
	cp support_files/ircd.yaml $(shell ls -d /opt/oragono*)/ircd.yaml
	cp support_files/ircd.service /etc/systemd/system/ircd.service
	systemctl daemon-reload
	systemctl enable ircd.service
	systemctl start ircd.service
endif

