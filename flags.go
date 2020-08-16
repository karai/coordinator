package main

import "flag"

// parseFlags This evaluates the flags used when the program was run
// and assigns the values of those flags according to sane defaults.
func flags() {
	flag.StringVar(&matrixToken, "matrixToken", "", "Matrix homeserver token")
	flag.StringVar(&graphDir, "graphDir", "./graph", "Path where graph objects should be saved")
	flag.StringVar(&batchDir, "batchDir", "./graph/batch", "Path where batched transactions should be saved")
	flag.StringVar(&matrixURL, "matrixURL", "", "Matrix homeserver URL")
	flag.StringVar(&matrixRoomID, "matrixRoomID", "", "Room ID for matrix publishd events")
	flag.IntVar(&chunkSize, "chunkSize", 100, "Number of transactions per batch on disk.")
	flag.IntVar(&karaiAPIPort, "apiport", 4200, "Port to run Karai Coordinator API on.")
	flag.BoolVar(&consumeData, "consume", false, "Consume data from sources.")
	flag.BoolVar(&isCoordinator, "coordinator", false, "Run as coordinator.")
	flag.BoolVar(&wantsClean, "clean", false, "Clear all peer certs and graph objects")
	flag.BoolVar(&wantsFiles, "write", true, "Write each graph object to disk.")
	flag.BoolVar(&wantsMatrix, "matrix", false, "Enable Matrix functions. Requires -matrixToken, -matrixURL, and -matrixRoomID")
	flag.Parse()
}
