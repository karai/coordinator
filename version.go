package main

// semverInfo Version string constructor
func semverInfo() string {
	majorSemver := "0"
	minorSemver := "19"
	patchSemver := "19"
	wholeString := majorSemver + "." + minorSemver + "." + patchSemver
	return wholeString
}
