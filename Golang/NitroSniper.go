package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

var (
	Token string
)

func Claim(Code string) int {
	if _, err := strconv.Atoi(Code); err != nil {
		Client := &http.Client{}
		req, _ := http.NewRequest("POST", "https://ptb.discordapp.com/api/v6/entitlements/gift-codes/"+Code+"/redeem", nil)
		req.Header.Set("Authorization", Token)
		Res, _ := Client.Do(req)
		return Res.StatusCode
	}
	return 0
}

func Message(s *discordgo.Session, Message *discordgo.MessageCreate) {
	if Message.Author != nil && !Message.Author.Bot { // would use a return but I for some reason like this better
		regex := regexp.MustCompile("[A-Za-z0-9]{24}|[A-Za-z0-9]{16}").FindStringSubmatch(Message.Content)
		if len(regex) > 0 {
			if Claim(regex[0]) == 200 {
				fmt.Println("claimed: https://discord.gift/" + regex[0])
			}
		}
	}
}

func Edit(s *discordgo.Session, Message *discordgo.MessageUpdate) {
	if Message.Author != nil && !Message.Author.Bot {
		regex := regexp.MustCompile("[A-Za-z0-9]{24}|[A-Za-z0-9]{16}").FindStringSubmatch(Message.Content)
		if len(regex) > 0 {
			if Claim(regex[0]) == 200 {
				fmt.Println("claimed: https://discord.gift/" + regex[0])
			}
		}
	}
}

func main() {
	fmt.Print("Enter your Token: ")
	fmt.Scan(&Token)
	dg, _ := discordgo.New(Token)
	dg.AddHandler(Message)
	dg.AddHandler(Edit)
	dg.Open()
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
}
