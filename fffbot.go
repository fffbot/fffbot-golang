package main

import (
	"strings"
	"errors"
	"log"
	"os"

	"net/http"
	"io/ioutil"

	"github.com/turnage/graw"
	"github.com/turnage/graw/reddit"

	"github.com/fffbot/fffbot/html2md"
)

func fetchPage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func parsePageToReply(url string, body []byte) (string, error) {
	str := string(body)
	headerIndex := strings.Index(str, "<h2")

	if headerIndex == -1 {
		return "", errors.New("header not found")
	}

	afterHeader := string(([]rune(str))[headerIndex:])

	footerIndex := strings.Index(afterHeader, "\"footer\"")

	if footerIndex == -1 {
		return "", errors.New("footer not found")
	}

	fromHeaderToFooter := string(([]rune(afterHeader))[:footerIndex])

	footerDivIndex := strings.LastIndex(fromHeaderToFooter, "<div")

	if footerDivIndex == -1 {
		return "", errors.New("footer div not found")
	}

	fromHeaderToFooterDiv := string(([]rune(fromHeaderToFooter))[:footerDivIndex])

	md := strings.Replace(strings.Replace(html2md.Convert(fromHeaderToFooterDiv), "(/blog/)", "(https://www.factorio.com/blog/)", 1), "    Hello", "Hello", 1) +
		`

*****
^(Fetched from: ` + url + `)

^(Beep boop I'm a bot; reply or message me to share bugs & suggestions)
`
	return md, nil
}

type reminderBot struct {
	bot reddit.Bot
}

func (r *reminderBot) Post(p *reddit.Post) error {
	if strings.Contains(p.URL, "factorio.com/blog/post/fff") {
		url := p.URL

		Info.Println("FFF post detected; Title:", p.Title, "; URL:", url)
		page, err := fetchPage(url)

		if err != nil {
			Error.Println("Error fetching from", url, ":", err)
			return nil
		}

		reply, err := parsePageToReply(url, page)

		if err != nil {
			Error.Println("Error parsing:", err)
			return nil
		}

		Info.Println("Reply formulated, posting")
		e := r.bot.Reply(p.Name, reply)
		if e != nil {
			Error.Println("Reply failed:", err, "; full reply:\n", reply)
		}
		return e
	}

	return nil
}

var (
	Info  *log.Logger
	Error *log.Logger
)

func main() {
	Info = log.New(os.Stdout, "[INF] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERR] ", log.Ldate|log.Ltime|log.Lshortfile)

	Info.Println("Initializing fffbot/1.0")

	if bot, err := botFromEnv(); err != nil {
		Error.Println("Failed to create bot handle:", err)
	} else {
		cfg := graw.Config{Subreddits: []string{"bottesting", "factorio"}}
		handler := &reminderBot{bot: bot}
		if _, wait, err := graw.Run(handler, bot, cfg); err != nil {
			Error.Println("Failed to start graw run:", err)
		} else {
			Info.Println("GRAW initialized, listening for posts")
			Error.Println("graw run failed:", wait())
		}
	}
}

func botFromEnv() (reddit.Bot, error) {
	agent := "GRAW:fffbot:1.0 (by /u/fffbot)"

	clientId := os.Getenv("GRAW_CLIENT_ID")
	clientSecret := os.Getenv("GRAW_CLIENT_SECRET")
	username := os.Getenv("GRAW_USERNAME")
	password := os.Getenv("GRAW_PASSWORD")

	app := reddit.App{ID: clientId,
		Secret: clientSecret,
		Username: username,
		Password: password}

	return reddit.NewBot(reddit.BotConfig{
		Agent: agent,
		App:   app,
		Rate:  0,
	})
}
