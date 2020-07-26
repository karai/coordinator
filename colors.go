package main

import "runtime"

var nc = "\033[0m"

var brightblack = "\033[1;30m"
var brightred = "\033[1;31m"
var brightgreen = "\033[1;32m"
var brightyellow = "\033[1;33m"
var brightpurple = "\033[1;34m"
var brightmagenta = "\033[1;35m"
var brightcyan = "\033[1;36m"
var brightwhite = "\033[1;37m"

var black = "\033[0;30m"
var red = "\033[0;31m"
var green = "\033[0;32m"
var yellow = "\033[0;33m"
var purple = "\033[0;34m"
var magenta = "\033[0;35m"
var cyan = "\033[0;36m"
var white = "\033[0;37m"

func osCheck() {
	if runtime.GOOS == "windows" {
		nc = ""

		brightblack = ""
		brightred = ""
		brightgreen = ""
		brightyellow = ""
		brightpurple = ""
		brightmagenta = ""
		brightcyan = ""
		brightwhite = ""

		black = ""
		red = ""
		green = ""
		yellow = ""
		purple = ""
		magenta = ""
		cyan = ""
		white = ""
	}
}
