package main

import (
	"strings"
)

func ExamplePrintGraph() {
	jsonData := `
	{
		"tx_type": 1,
		"tx_hash": "1cae10466de1b4bca6197bd38bc2cfea961ec9d6b21baf2923caad523ca7fa99",
		"tx_prev": "1ae3529687d8340213f6ffc2c5a4d2747ff35575eb12841e7bfea60660cd69c1",
		"tx_data": [
		{
			"coord_pub": "cfae4ecca8ed282ab51ae16bf755e510b1ce52d6e84263f357b0a524691e9259",
			"coord_sig": "2efc5209e8f19c7d0f50ac08107b0c5690ba173c6371d8d9a9b46b1790b4e709ce7e492671cb76203ff4b40e40bbaf2075db7741eeae8b4a641fe36bfa8b880b",
			"signed_msg": "48a49aa318fa36ee623bae1b25e23d87a21923d889f620e62106009699de17fc"
		},
		{
			"participant_matrix": [
				{
					"index": 0,
					"participants": [
						"77a4781055c9d26c3136afefa7823037e898f92105bcaaeac385719d42b50f20",
						"b7e0d0839883bf207b39e032343fa1569dbfeb6075c9fa9e5de41e025c367e90",
						"6eab8b2f2e55e0b10a527e2a9ec85acc3c89e3f9cb92fa30e9af0bc49b5181c3"
					]
				}
			]
		}
	  ]
	}`

	rdr := strings.NewReader(jsonData)

	printGraph(rdr)

	// Output:
	// here we go
	// 1cae10466de1b4bca6197bd38bc2cfea961ec9d6b21baf2923caad523ca7fa99
	// 1ae3529687d8340213f6ffc2c5a4d2747ff35575eb12841e7bfea60660cd69c1
	// 1
	// [
	// 		{
	//			"coord_pub": "cfae4ecca8ed282ab51ae16bf755e510b1ce52d6e84263f357b0a524691e9259",
	//			"coord_sig": "2efc5209e8f19c7d0f50ac08107b0c5690ba173c6371d8d9a9b46b1790b4e709ce7e492671cb76203ff4b40e40bbaf2075db7741eeae8b4a641fe36bfa8b880b",
	//			"signed_msg": "48a49aa318fa36ee623bae1b25e23d87a21923d889f620e62106009699de17fc"
	//		},
	//		{
	//			"participant_matrix": [
	//				{
	//					"index": 0,
	//					"participants": [
	//						"77a4781055c9d26c3136afefa7823037e898f92105bcaaeac385719d42b50f20",
	//						"b7e0d0839883bf207b39e032343fa1569dbfeb6075c9fa9e5de41e025c367e90",
	//						"6eab8b2f2e55e0b10a527e2a9ec85acc3c89e3f9cb92fa30e9af0bc49b5181c3"
	//					]
	//				}
	//			]
	//		}
	//	  ]
}
