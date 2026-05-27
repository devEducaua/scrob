PREFIX ?= /usr/local
BINDIR ?= $(PREFIX)/bin

TARGET = scrob
TARGETDIR = bin

SRCS = $(wildcard *.go) $(wildcard */*.go)

.PHONY: clean install uninstall

$(TARGET): $(SRCS)
	go build -o $(TARGETDIR)/$(TARGET) .

install:
	cp $(TARGETDIR)/$(TARGET) $(BINDIR)/

uninstall:
	rm $(BINDIR)/$(TARGET)

clean:
	rm -r $(TARGETDIR)
