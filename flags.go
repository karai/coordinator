package main

import "flag"

// parseFlags This evaluates the flags used when the program was run
// and assigns the values of those flags according to sane defaults.
func parseFlags() {
	flag.IntVar(&karaiAPIPort, "apiport", 4200, "Port to run Karai Coordinator API on.")
	flag.IntVar(&karaiP2PPort, "p2pport", 4201, "Port to listen for P2P messages on.")
	flag.BoolVar(&isCoordinator, "coordinator", false, "Run as coordinator.")
	flag.Parse()
}
