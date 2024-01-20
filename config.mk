# templa version
VERSION = 0.6.1

# paths
PREFIX = /usr/local
MANPREFIX = $(PREFIX)/share/man

# flags
LDFLAGS = '-w -s -X main.VERSION=$(VERSION)'
