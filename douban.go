package dbfm

import (
	"encoding/json"
	"fmt"
	color "github.com/daviddengcn/go-colortext"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	// "strconv"
	"time"
)

type FmChannel struct {
	Name_en string `json:"name_en"`
	Seq_id  int    `json:"seq_id"`
	Abbr_en string `json:"abbr_en"`
	Name    string `json:"name"`
	// Channel_id string    `json:"channel_id"`
	Channel_id interface{} `json:"channel_id"`
}

type FmChannels struct {
	Channels []FmChannel `json:"channels"`
}

func (c *FmChannels) Print() {
	channel_len := len(c.Channels)
	if channel_len > 40 {
		channel_len = 40
	}

	for i := 0; i < channel_len-1; i += 2 {
		fmt.Printf("%3d. %v\t%3d. %s\n", i, c.Channels[i].Name, i+1, c.Channels[i+1].Name)
	}
}

type Song struct {
	Album   string `json:"album"`
	Picture string `json:"picture"`
	Url     string `json:"url"`
	Title   string `json:"title"`
}

type SongResponse struct {
	R                   int    `json:"r"`
	Is_show_quick_start int    `json:"is_show_quick_start"`
	Songs               []Song `json:"song"`
}

type FmAccount struct {
	User_name string `json:"user_name"`
	User_id   string `json:"user_id"`
	Token     string `json:"token"`
	Expire    string `json:"expire"`
	Err       string `json:"err"`
	R         int    `json:"r"`
	Email     string `json:"email"`
}

func (fm *Fm) Login() *FmAccount {
	login_url := "http://www.douban.com/j/app/login"
	data := make(url.Values)
	data["email"] = []string{fm.Email}
	data["password"] = []string{fm.Password}
	data["app_name"] = []string{"radio_desktop_win"}
	data["version"] = []string{"100"}

	res, err := http.PostForm(login_url, data)
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer res.Body.Close()

	var reader io.ReadCloser
	account := &FmAccount{}

	reader = res.Body
	if reader != nil {
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			log.Fatalln(err)
		}

		err = json.Unmarshal(body, &account)
		if err != nil {
			log.Fatalln(err)
		}
	}
	return account
}

func (fm *Fm) GetPlaylist(channel string) {
	req_url, _ := url.Parse("http://douban.fm/j/mine/playlist")
	req_url.RawQuery = url.Values{
		"channel": {channel},
		"type":    {"n"},
		"pb":      {"64"},
		"from":    {"mainsite"},
	}.Encode()
	res, _ := http.Get(req_url.String())
	defer res.Body.Close()

	data, _ := ioutil.ReadAll(res.Body)

	songs := &SongResponse{}

	err := json.Unmarshal(data, &songs)
	if err != nil {
		log.Fatalln(err)
	}
	for _, s := range songs.Songs {
		fm.PlayList <- s
	}
}

// not use
func (fm *Fm) ShuffleSongs(songs []Song) {
	for i, _ := range rand.Perm(len(songs)) {
		fm.PlayList <- songs[i]
	}
}

func GetChannels() *FmChannels {
	req_url := `http://www.douban.com/j/app/radio/channels`
	res, err := http.Get(req_url)
	if err != nil {
		os.Exit(1)
	}
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		os.Exit(1)
	}
	channels := new(FmChannels)

	err = json.Unmarshal(data, &channels)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return channels
}

func (c *FmChannels) CMap() map[int]string {
	cMap := make(map[int]string)
	for i, channel := range c.Channels {
		cMap[i] = fmt.Sprintf("%v", channel.Channel_id)
	}
	return cMap
}

type Fm struct {
	Player   string
	Email    string
	Password string
	PlayId   string
	PlayList chan Song
	Playing bool
}

func NewFm() *Fm {
	fm := &Fm{}
	fm.PlayList = make(chan Song, 1000)
	return fm
}

var current string

func (fm *Fm) MainLoop() {
	fm.Play()

	for {
		fmt.Scanf("%s\n", &current)

		switch {
		case current == "p":
			fm.Play()
		case current == "s":
			fm.Stop()
		case current == "n":
			fm.Next()
		case current == "q":
			fm.Exit()
		}

		time.Sleep(time.Second * 1)
	}
}

func (fm *Fm) Play() {
	if fm.Playing == true {
		return
	}

	fm.Playing = true

	go func(){
		for {
			if fm.Playing == false{
				return
			}

			if len(fm.PlayList) <= 1 {
				fm.GetPlaylist(fm.PlayId)
			}

			current_song := <-fm.PlayList

			color.ChangeColor(color.Red, true, color.Black, false)
			fmt.Print(`♥ `)
			color.ChangeColor(color.Green, false, color.Black, false)
			fmt.Print(current_song.Title + " ")
			color.ResetColor()

			time.Sleep(time.Second * 1)

			cmd := exec.Command(fm.Player, current_song.Url)
			cmd.Start()
			cmd.Wait()
		}
	}()
}

func (fm *Fm) Next() {
	cmd := exec.Command("taskkill", "/F", "/IM", "mplayer.exe")
	cmd.Run()
}

func (fm *Fm) Stop() {
	fm.Playing = false

	cmd := exec.Command("taskkill", "/F", "/IM", "mplayer.exe")
	/*
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	*/
	cmd.Run()
}

func (fm *Fm) Exit() {
	fm.Stop()
	os.Exit(0)
}
