// Package girc contains the core IRC connectivity functionality
package girc

import (
	"bufio"
	"fmt"
	"net"
	"regexp"
	"strings"

	"github.com/MasterGunner/GoGoGunnerBot/utilities"
)

// IRCInterface does the...interface stuff for the IRC Connection? I still really don't quite understand what Interfaces are for.
type IRCInterface interface {
	StartupConfig(string, int, []string, string)
	Connect()
	Join(string)
	Leave(string)
	AddListener(listener)
	RemoveListener(string)
	Handle(string)
	Send(string)
	Say(string, string)
}

// IRC struct holds the variables for the IRC connection.
type IRC struct {
	connection *net.Conn
	//connection *io.ReadWriter
	server    string
	port      int
	channels  []string
	nick      string
	listeners []listener //Probably have to create a new Type/Struct for listeners, and then create an array of those here.
}

// StartupConfig configures the socket and initial recievers.
func (i *IRC) StartupConfig(server string, port int, channels []string, nick string) {
	//"constructor" operations
	i.server = server
	i.port = port
	i.channels = channels
	i.nick = nick

	//Register listener for Ping/Pong.

	//Register Other Listeners
	i.AddListener(listener{
		name:  "echo",
		query: regexp.MustCompile(" PRIVMSG (.*) :}Echo (.*)"),
		callback: func(i *IRC, info []string) {
			i.Say(info[1], info[2])
		},
		helptext: "Echo Message",
	})
}

// Connect sets up the connection to the IRC server.
func (i *IRC) Connect() {
	utilities.Log("Establishing Connection...")

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", i.server, i.port))
	if err != nil {
		utilities.Log(fmt.Sprintf("Could not connect to server: %v", err))
		return
	}
	i.connection = &conn

	i.Send("NICK " + i.nick)
	i.Send("USER " + i.nick + " 8 *:Golang IRC bot")

	//Set socket listeners for connect and disconnect. <-Does GO even have these?
	////Wait for connection to be established, and then send connection details.
	////If channels are specified, join them after 3 seconds.
	//On disconnect...

	//Loop to listen to Socket; call Handler function in new thread on new input
	quit := false
	for quit != true {
		//message, err := i.connection.Read()
		message, err := bufio.NewReader(*i.connection).ReadString('\n')
		if err != nil {
			utilities.Log(fmt.Sprintf("%v", err))
		}

		go i.Handle(string(message))
	}
}

// Join an IRC Channel.
func (i *IRC) Join(channel string) {
	//IRC.Send('JOIN ' + channel) ?
	utilities.Log("JOINED CHANNEL - " + channel)
}

// Leave an IRC Channel
func (i *IRC) Leave(channel string) {
	//IRC.Send('JOIN ' + channel) ?
	utilities.Log("LEFT CHANNEL - " + channel)
}

// AddListener sets new listeners for incoming messages.
func (i *IRC) AddListener(l listener) {
	i.listeners = append(i.listeners, l)
}

// RemoveListener remove listeners for incoming messages.
func (i *IRC) RemoveListener(name string) {

}

// Handle decides which listeners to pass incoming messages on to.
func (i *IRC) Handle(message string) {
	utilities.Log(fmt.Sprintf("RECV - %s", message))
	if strings.Index(message, "PING ") == 0 {
		i.Send(fmt.Sprintf("PONG %s", message[5:]))
	}

	for j := 0; j < len(i.listeners); j++ {
		l := i.listeners[j]
		info := l.query.FindStringSubmatch(message)
		if len(info) > 1 {
			l.callback(i, info)
		}
	}
}

//Send raw information.
func (i *IRC) Send(message string) {
	fmt.Fprint(*i.connection, fmt.Sprintf("%s\r\n", message))
	utilities.Log(message)
}

//Say something to a channel
func (i *IRC) Say(channel string, message string) {
	if strings.Index(message, "ACTION") == 0 {
		message = fmt.Sprintf("\x01%s\x01", message)
	}
	i.Send(fmt.Sprintf("PRIVMSG %s :%s\r\n", channel, message))
	utilities.Log(fmt.Sprintf("SAID(%s) - %s", channel, message))
}
