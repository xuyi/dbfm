package main

import (
	"fmt"
	"github.com/xuyi/dbfm"
)

func main() {
	fm := dbfm.NewFm()
	fm.Email = ""
	fm.Password = ""
	fm.Player = `C:\Program Files (x86)\mplayer\mplayer.exe`

	channels := dbfm.GetChannels()
	channels.Print()
	cMap := channels.CMap()

	var playlist_index int
	fmt.Print("Channel:(0) ")
	fmt.Scanf("%d\n", &playlist_index)

	// account := fm.Login()
	// songs := account.GetPlaylist(cMap[playlist_index])
	songs := dbfm.GetPlaylist(cMap[playlist_index])
	go fm.ShuffleSongs(songs)
	fm.MainLoop()
}
