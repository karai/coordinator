package main

type clientHeader struct {
	ClientHeaderAppName    string `json:"client_header_app_name"`
	ClientHeaderAppVersion string `json:"client_header_app_version"`
	ClientHeaderPeerID     string `json:"client_header_peer_id"`
	ClientProtocolVersion  string `json:"client_protocol_version"`
}

func connectChannel(ktx string) bool {
	// connect
	// listen for welcome
	// Initial Connection Sends N1:PK to Coord
	// if we are returning, validate signed N1:C
	// Upon successful connection, submit joinTx
	// if joinTx published, return true on connectChannel() for success
	return true
}

// The N1 also needs to know how to construct the join message, so I should add that in parallel to clientops.go
// https://discordapp.com/channels/388915017187328002/453726546868305962/719243359440339023
