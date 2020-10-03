package main

// semverInfo Version string constructor
func semverInfo() string {
	majorSemver := "0"
	minorSemver := "21"
	patchSemver := "2"
	wholeString := majorSemver + "." + minorSemver + "." + patchSemver
	return wholeString
}
