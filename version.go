package main

// semverInfo Version string constructor
func semverInfo() string {
	var majorSemver, minorSemver, patchSemver, wholeString string
	majorSemver = "0"
	minorSemver = "10"
	patchSemver = "1"
	wholeString = majorSemver + "." + minorSemver + "." + patchSemver
	return wholeString
}
