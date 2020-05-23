package main

import "flag"

// parseFlags This evaluates the flags used when the program was run
// and assigns the values of those flags according to sane defaults.
func parseFlags() {
	flag.IntVar(&karaiPort, "port", 4200, "Port to run Karai Coordinator on.")
	flag.BoolVar(&isCoordinator, "coordinator", false, "Run as coordinator.")
	// flag.StringVar(&karaiPort, "karaiPort", "4200", "Port to run Karai")
	flag.Parse()
}
