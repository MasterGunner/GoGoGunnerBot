package girc

import (
	"regexp"
)

type listener struct {
	name     string
	query    *regexp.Regexp       //regex
	callback func(*IRC, []string) //the function to execute, no idea how this actually would work in Go. Check to see hwo I did this in Ruby?
	helptext string
}

// RegisterListeners populates the array of Listener functions ("temporary" function, as I have no way to dynamically load Listeners right now).
func RegisterListeners(i *IRC) {
	/*
	 * ADMINISTRATIVE LISTENERS
	 */

	//Join the designated channel
	i.AddListener(listener{
		name:  "JoinChannel",
		query: regexp.MustCompile("(?i) PRIVMSG (.*) :" + i.commandChar + "Join (.*)"),
		callback: func(i *IRC, info []string) {
			i.Join(info[2])
			i.channels = append(i.channels, info[2])
		},
		helptext: "Make GunnerBot join a new channel. Usage: " + i.commandChar + "Join #Channel",
	})

	//Leave the designated channel
	i.AddListener(listener{
		name:  "LeaveChannel",
		query: regexp.MustCompile("(?i) PRIVMSG (.*) :" + i.commandChar + "Leave ?$"),
		callback: func(i *IRC, info []string) {
			i.Leave(info[1])

			//Remove from internal channel list
			for j := 0; j < len(i.channels); j++ {
				if i.channels[j] == info[1] { //Find channel with matching name.
					i.channels = append(i.channels[:j], i.channels[j+1:]...) //Remove channel from list.
					j = j - 1                                                //Decrease counter to avoid skipping the next channel.
				}
			}
		},
		helptext: "Make GunnerBot leave the channel. Usage: " + i.commandChar + "Leave",
	})

	//Bot Source
	i.AddListener(listener{
		name:  "Source",
		query: regexp.MustCompile("(?i) PRIVMSG (.*) :" + i.commandChar + "Source ?$"),
		callback: func(i *IRC, info []string) {
			i.Say(info[1], "Don't judge me! https://github.com/MasterGunner/GoGoGunnerBot")
		},
	})

	//Basic echo/repeat/say function
	i.AddListener(listener{
		name:  "Echo",
		query: regexp.MustCompile("(?i) PRIVMSG (.*) :" + i.commandChar + "Echo (.*)"),
		callback: func(i *IRC, info []string) {
			i.Say(info[1], info[2])
		},
		helptext: "Echo the provided text. Usage: " + i.commandChar + "Echo DerpDerpDerp",
	})
}
