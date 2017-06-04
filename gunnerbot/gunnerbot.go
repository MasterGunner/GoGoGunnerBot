package main

import (
	"github.com/MasterGunner/GoGoGunnerBot/girc"
)

func main() {
	var ircInterface girc.IRCInterface
	irc := girc.IRC{}

	ircInterface = &irc

	ircInterface.StartupConfig("irc.dbcommunity.org", 6667, []string{"desertbus"}, "GoGoGunnerBot")
	ircInterface.Connect()
}
