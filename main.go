package main

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/gookit/ini/v2/dotnev"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

func main() {
	vote()

	duration := 3*time.Hour + 1*time.Minute

	loadingbar(duration)
	c := time.Tick(duration)
	for _ = range c {
		loadingbar(duration)
		vote()
	}
}

func loadingbar(duration time.Duration) {

	tickduration := time.Minute

	p := mpb.New(mpb.WithWidth(64))
	total := duration / tickduration
	name := "Timer :"

	bar := p.New(int64(total),
		mpb.BarStyle().Lbound("╢").Filler("▌").Tip("▌").Padding("░").Rbound("╟"),
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
		),
		mpb.AppendDecorators(decor.Percentage()),
	)

	for i := 0; i < int(total); i++ {
		time.Sleep(tickduration)
		bar.Increment()
	}

	p.Wait()
}

func vote() {

	err := dotnev.Load("./", ".env")

	if err != nil {
		fmt.Println(err)
		return
	}

	username := dotnev.Get("USERNAME") // Username.
	password := dotnev.Get("PASSWORD") // Password.

	// Create the cookie jar.
	jar, err := cookiejar.New(nil) // Create the cookie jar.

	if err != nil { // Check if the cookie jar is created.
		fmt.Println(err)
		return
	}

	client := &http.Client{Jar: jar} // Create the client.

	res, err := client.Get("https://level-flyff.fr/") // Get the page index.
	if err != nil {
		fmt.Println(err)
		return
	}

	passwordUser := md5.Sum([]byte("kikugalanet" + password))                                         // Hash the password.
	data := url.Values{"login": {username}, "password": {string(passwordUser[:])}, "remember": {"0"}} // Create the data.

	res, err = client.PostForm("https://level-flyff.fr/ajax/connexion.php", data) // Post the cookie.
	if err != nil {
		fmt.Println(err)
		return
	}

	if res.StatusCode == 200 {
		fmt.Println("Vous êtes connecté avec", dotnev.Get("USERNAME"), "!")
	}

	dataVote := url.Values{"id": {"1"}}                                                 // Create the data.
	res, err = client.PostForm("https://level-flyff.fr/ajax/recompenses.php", dataVote) // Get the page index.
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Vote effectué !")
}
