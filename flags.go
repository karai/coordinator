package main

import "flag"

// parseFlags This evaluates the flags used when the program was run
// and assigns the values of those flags according to sane defaults.
func flags() {
	flag.StringVar(&matrixToken, "matrixToken", "", "Matrix homeserver token")
	flag.StringVar(&matrixURL, "matrixURL", "", "Matrix homeserver URL")
	flag.StringVar(&matrixRoomID, "matrixRoomID", "", "Room ID for matrix publishd events")
	flag.IntVar(&karaiAPIPort, "apiport", 4200, "Port to run Karai Coordinator API on.")
	flag.BoolVar(&wantsClean, "clean", false, "Clear all peer certs")
	flag.BoolVar(&wantsMatrix, "matrix", false, "Enable Matrix functions. Requires -matrixToken, -matrixURL, and -matrixRoomID")
	flag.Parse()
}
