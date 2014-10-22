package main

import (
	"fmt"
	"github.com/xuyi/dbfm"
)

func main() {
	channels := dbfm.GetChannels()
	channels.Print()
	cMap := channels.CMap()

	var playlist_index int
	fmt.Print("Channel:(0) ")
	fmt.Scanf("%d\n", &playlist_index)

	fm := dbfm.NewFm()
	fm.Email = ""
	fm.Password = ""
	fm.Player = `C:\Program Files (x86)\mplayer\mplayer.exe`
	fm.PlayId = cMap[playlist_index]

	fm.MainLoop()
}
