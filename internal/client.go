package internal

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type Client struct {
	username    string
	currentRoom *Room
	conn        net.Conn
	commandChan chan<- Command
}

func (client *Client) err(err error) {
	client.conn.Write([]byte(fmt.Sprintf("error: %s\n", err.Error())))
}

func (client *Client) msg(message string) {
	client.conn.Write([]byte(fmt.Sprintf("> %s\n", message)))
}

func readClientInput(client *Client) {
	for {
		message, err := bufio.NewReader(client.conn).ReadString('\n')

		if err != nil {
			return
		}

		message = strings.Trim(message, "\r\n")
		args := strings.Split(message, " ")
		commandArg := strings.TrimSpace(args[0])
		args = args[1:]

		command := Command{
			id:     CMD_DEFAULT,
			client: client,
			args:   args,
		}

		switch commandArg {
		case "/commands":
			command.id = CMD_COMMANDS
		case "/username":
			command.id = CMD_USERNAME
		case "/join":
			command.id = CMD_JOIN
		case "/msg":
			command.id = CMD_MSG
		case "/dm":
			command.id = CMD_DM
		case "/room":
			command.id = CMD_ROOM
		case "/members":
			command.id = CMD_MEMBERS
		case "/rooms":
			command.id = CMD_ROOMS
		case "/quit":
			command.id = CMD_QUIT
		case "/help":
			command.id = CMD_HELP
		}

		client.commandChan <- command
	}
}

func NewClient(conn net.Conn, commandChan chan Command) {
	client := &Client{
		username:    "anonymous",
		currentRoom: nil,
		conn:        conn,
		commandChan: commandChan,
	}

	client.msg("welcome to tcp chat server!")

	readClientInput(client)
}
