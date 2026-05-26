
TARGET = scrob
TARGETDIR = bin

SRCS = $(wildcard *.go) $(wildcard */*.go)

.PHONY: clean

$(TARGET): $(SRCS)
	go build -o $(TARGETDIR)/$(TARGET) .

clean:
	rm -r $(TARGETDIR)
