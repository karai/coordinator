package main

// semverInfo Version string constructor
func semverInfo() string {
	majorSemver := "0"
	minorSemver := "22"
	patchSemver := "0"
	wholeString := majorSemver + "." + minorSemver + "." + patchSemver
	return wholeString
}
