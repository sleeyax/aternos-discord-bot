package message

import aternos "github.com/sleeyax/aternos-api"

var colorMap map[aternos.ServerStatus]int = map[aternos.ServerStatus]int{
	aternos.Online:    0x04ed23,
	aternos.Offline:   0xd12222,
	aternos.Stopping:  0xd19122,
	aternos.Starting:  0xeeff01,
	aternos.Loading:   0x01eaff,
	aternos.Preparing: 0x01eaff,
	aternos.Saving:    0xc9ba99,
}
