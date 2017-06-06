// Package girc contains the core IRC connectivity functionality
package girc

import (
	"bufio"
	"fmt"
	"net"
	"strings"

	"github.com/MasterGunner/GoGoGunnerBot/utilities"
)

// IRCInterface does the...interface stuff for the IRC Connection? I still really don't quite understand what Interfaces are for.
type IRCInterface interface {
	NewClient(string, int, []string, string, string)
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
	server      string
	port        int
	channels    []string
	nick        string
	commandChar string
	listeners   []listener //Probably have to create a new Type/Struct for listeners, and then create an array of those here.
}

// NewClient creates and returns an IRC Client instance.
func NewClient(server string, port int, channels []string, nick string, commandChar string) *IRC {
	return &IRC{
		server:      server,
		port:        port,
		channels:    channels,
		nick:        nick,
		commandChar: commandChar,
	}
}

// Connect sets up the connection to the IRC server.
func (i *IRC) Connect() {
	//Establish basic connection
	utilities.Log("Establishing Connection...")

	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", i.server, i.port))
	if err != nil {
		utilities.Log(fmt.Sprintf("Could not connect to server: %v", err))
		return
	}
	i.connection = &conn

	i.Send("NICK " + i.nick)
	i.Send("USER " + i.nick + " 8 *:Golang IRC bot")

	//Register Listeners (should probably do this before establishing connection)
	RegisterListeners(i)

	//Connect to provided channels
	for j := 0; j < len(i.channels); j++ {
		i.Join(i.channels[j])
	}

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
	i.Send("JOIN " + channel)
	utilities.Log("JOINED CHANNEL - " + channel)
}

// Leave an IRC Channel
func (i *IRC) Leave(channel string) {
	i.Send("PART " + channel)
	utilities.Log("LEFT CHANNEL - " + channel)
}

// AddListener sets new listeners for incoming messages.
func (i *IRC) AddListener(l listener) {
	i.listeners = append(i.listeners, l)
}

// RemoveListener remove listeners for incoming messages.
func (i *IRC) RemoveListener(name string) {
	//a = append(a[:i], a[i+1:]...)
	for j := 0; j < len(i.listeners); j++ {
		if i.listeners[j].name == name { //Find listener with matching name.
			i.listeners = append(i.listeners[:j], i.listeners[j+1:]...) //Remove listener from list.
			j = j - 1                                                   //Decrease counter to avoid skipping the next listener.
		}
	}
}

// Handle decides which listeners to pass incoming messages on to.
func (i *IRC) Handle(message string) {
	//Log Recieved Message
	utilities.Log(fmt.Sprintf("RECV - %s", message))

	//Repond to Ping/Pong queries from server.
	if strings.Index(message, "PING ") == 0 {
		i.Send(fmt.Sprintf("PONG %s", message[5:]))
	}

	//Loop through listener arrays and trigger any functions as appropriate.
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
