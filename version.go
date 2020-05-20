package main

// Version string
func semverInfo() string {
	var majorSemver, minorSemver, patchSemver, wholeString string
	majorSemver = "0"
	minorSemver = "5"
	patchSemver = "3"
	wholeString = majorSemver + "." + minorSemver + "." + patchSemver
	return wholeString
}
