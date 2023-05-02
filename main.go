package main

import (
	"fmt"
	traffic "traffic/httpreq"

	"github.com/enescakir/emoji"
)

func main() {
	traffic.GetUpdates()
	fmt.Printf(string(emoji.HighSpeedTrain), string(emoji.BulletTrain))
}
