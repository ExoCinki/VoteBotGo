package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/gookit/ini/v2/dotnev"
	"golang.org/x/net/publicsuffix"
)

func main() {

	vote()

	duration := 3*time.Hour + 1*time.Minute // PROD
	//duration := 1 * time.Minute

	c := time.Tick(duration)
	for _ = range c {
		vote()
	}
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
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List}) // Create the cookie jar.

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
	hasher := md5.New()                                 // Create the hasher.
	hasher.Write([]byte("kikugalanet" + password))      // Hash the password.
	hex.EncodeToString(hasher.Sum(nil))                 // Encode the hash.
	passwordUser := hex.EncodeToString(hasher.Sum(nil)) // Hash the password.
	fmt.Println("password : ", passwordUser)
	data := url.Values{"login": {username}, "password": {string(passwordUser[:])}, "remember": {"0"}} // Create the data.

	//5e151efef74f1eee52fee448ca7d3158

	res, err = client.PostForm("https://level-flyff.fr/ajax/connexion.php", data) // Post the cookie.
	if err != nil {
		fmt.Println(err)
		return
	}

	if res.StatusCode == 200 {
		fmt.Println("Vous êtes connecté avec", dotnev.Get("USERNAME"), "!")
	}

	fmt.Println("Wating 5 seconds...")
	time.Sleep(5 * time.Second) // Wait 5 seconds.

	fmt.Println("Accès to Vote Button !!")
	fmt.Println("Wating 5 seconds...")
	time.Sleep(5 * time.Second) // Wait 5 seconds.

	dataVote := url.Values{"id": {"1"}}                                                 // Create the data.
	res, err = client.PostForm("https://level-flyff.fr/ajax/recompenses.php", dataVote) // Get the page index.
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Vote effectué !")
	fmt.Println("Heures du vote : " + time.Now().Format("02-01-2006 15:04:05"))
}
