.PHONY: postgres

GO_INSTALLED   := $(shell which go)
GO_VERSION     ?= 1.14.1

ifeq ($(findstring arm,$(shell uname -m)), arm)
	ARCH = armv6l
else ifeq ($(findstring x86,$(shell uname -m)), x86)
	ARCH = amd64
else 
$(error unable to determine ARCH)
endif

ifeq ($(findstring Linux,$(shell uname -s)), Linux)
	OS = linux
else ifeq ($(findstring Darwin,$(shell uname -s)), Darwin)
	OS = darwin

else
$(error unable to determine OS)
endif

install-dev: go-get-private build_tools ircd postgress

build-tools: 
ifeq ($(GO_INSTALLED),)
	$(info installing Golang v${GO_VERSION} on OS=${OS}, ARCH=${ARCH})
	curl -o /tmp/go${GO_VERSION}.${OS}-${ARCH}.tar.gz https://dl.google.com/go/go${GO_VERSION}.${OS}-${ARCH}.tar.gz
	sudo tar -C /usr/local -xzf /tmp/go${GO_VERSION}.${OS}-${ARCH}.tar.gz
	echo 'export PATH=$$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
endif # end check for go installed
	$(info installing git)
	sudo apt install git -y


go-get-private:
	go env -w GOPRIVATE=pault.ag
	git config --global url."git@github.com:".insteadOf "https://github.com/"

postgres:
	$(info installing and setting up postgresql)
	sudo apt install -y postgresql
	sudo service postgresql start
	sudo -u postgres psql -c "CREATE USER euler WITH PASSWORD 'euler'" 
	sudo -u postgres psql -c 'CREATE DATABASE euler'

####### IRC stuff #########
ifeq ($(findstring armv6,$(shell uname -m)), armv6)
IRC_ARCH = armv6
else ifeq ($(findstring armv7,$(shell uname -m)), armv7)
IRC_ARCH = armv7
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

ircd: ircd-stage-files
	$(info installing ircd)
ifneq ($(whoami), root)
	sudo mv $(shell ls -d /home/oragono/oragono-2.0.0*)/oragono /home/oragono/oragono
	sudo chown oragono:oragono /home/oragono/oragono
	sudo mv $(shell ls -d /home/oragono/oragono-2.0.0*)/languages /home/oragono/
	sudo chown -R oragono:oragono /home/oragono/languages
	sudo systemctl daemon-reload
	sudo systemctl enable oragono.service
	sudo systemctl start oragono.service
else 
	mv $(shell ls -d /home/oragono/oragono-2.0.0*)/oragono /home/oragono/oragono
	chown oragono:oragono /home/oragono/oragono
	mv $(shell ls -d /home/oragono/oragono-2.0.0*)/languages /home/oragono/
	chown -R oragono:oragono /home/oragono/languages
	systemctl daemon-reload
	systemctl enable oragono.service
	systemctl start oragono.service
endif

ircd-stage-files:
	$(info ircd config for os: $(IRC_OS) arch: $(IRC_OS))
	wget -O /tmp/oragono.tar.gz https://github.com/oragono/oragono/releases/download/v2.0.0/oragono-2.0.0-$(IRC_OS)-$(IRC_ARCH).tar.gz
ifneq ($(whoami), root)
	sudo adduser --system --group oragono
	sudo tar -xf /tmp/oragono.tar.gz -C /home/oragono
	sudo rm -rf /tmp/oragono.tar.gz
	sudo cp support_files/ircd.yaml /home/oragono/ircd.yaml
	sudo chown oragono:oragono /home/oragono/ircd.yaml
	sudo cp support_files/oragono.service /etc/systemd/system/oragono.service
else
	adduser --system --group oragono
	tar -xf /tmp/oragono.tar.gz -C /home/oragono
	rm -rf /tmp/oragono.tar.gz
	cp support_files/ircd.yaml /home/oragono/ircd.yaml
	chown oragono:oragono /home/oragono/ircd.yaml
	cp support_files/oragono.service /etc/systemd/system/oragono.service
endif

ircd-clean:
	$(info cleaning up ircd for os: $(IRC_OS) arch: $(IRC_OS))
ifneq ($(whoami), root)
	sudo systemctl stop oragono.service
	sudo systemctl disable oragono.service
	sudo deluser --remove-home oragono
	sudo rm /etc/systemd/system/oragono.service
	sudo systemctl daemon-reload
else 
	systemctl stop oragono.service
	systemctl disable oragono.service
	deluser --remove-home oragono
	rm /etc/systemd/system/oragono.service
	systemctl daemon-reload
endif

##### install code ####
install-tool: install-quaditor postgres build-tools

install-quaditor:
	go build -o ~/quaditorcli ./cmd/*.go
	cp examples/bob_alice/bob_friends_keys_query.json ~/bob_friends_keys_query.json
	cp examples/bob_alice/quads.json ~/keys_data_set.json
	~/quaditorcli load -f ~/keys_data_set.json