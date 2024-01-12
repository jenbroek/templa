# templa version
VERSION = 0.2.0

# paths
PREFIX = /usr/local
MANPREFIX = $(PREFIX)/share/man

# flags
LDFLAGS = '-X main.VERSION=$(VERSION)'
