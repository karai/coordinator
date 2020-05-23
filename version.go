package main

// semverInfo Version string constructor
func semverInfo() string {
	var majorSemver, minorSemver, patchSemver, wholeString string
	majorSemver = "0"
	minorSemver = "6"
	patchSemver = "0"
	wholeString = majorSemver + "." + minorSemver + "." + patchSemver
	return wholeString
}
