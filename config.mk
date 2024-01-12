# templa version
VERSION = 0.3.0

# paths
PREFIX = /usr/local
MANPREFIX = $(PREFIX)/share/man

# flags
LDFLAGS = '-w -s -X main.VERSION=$(VERSION)'
