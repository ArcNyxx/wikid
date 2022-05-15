# wikid - guessing game
# Copyright (C) 2022 ArcNyxx
# see LICENCE file for licensing information

.POSIX:

include config.mk

SRC = command.go logic.go wikid.go

all: wikid

wikid: $(SRC) config.mk
	$(GO) build -o $@ ./

clean:
	rm -f wikid wikid-$(VERSION).tar.gz

dist: clean
	mkdir -p wikid
	cp README LICENCE Makefile config.mk wikid.1 $(SRC) wikid
	tar -cf - wikid | gzip -c > wikid.tar.gz
	rm -rf wikid

install: all
	mkdir -p $(PREFIX)/bin $(MANPREFIX)/man1
	cp -f wikid $(PREFIX)/bin
	chmod 755 $(PREFIX)/bin/wikid
	sed 's/VERSION/$(VERSION)/g' < wikid.1 > $(MANPREFIX)/man1/wikid.1
	chmod 644 $(MANPREFIX)/man1/wikid.1

uninstall:
	rm -f $(PREFIX)/bin/wikid $(MANPREFIX)/man1/wikid.1

.PHONY: all clean dist install uninstall
