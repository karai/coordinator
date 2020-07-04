package main

// semverInfo Version string constructor
func semverInfo() string {
	var majorSemver, minorSemver, patchSemver, wholeString string
	majorSemver = "0"
	minorSemver = "16"
	patchSemver = "2"
	wholeString = majorSemver + "." + minorSemver + "." + patchSemver
	return wholeString
}
