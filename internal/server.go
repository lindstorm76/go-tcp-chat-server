package internal

import (
	"fmt"
	"strings"
)

type Server struct {
	rooms map[string]*Room

	CommandChan chan Command
}

func NewServer() *Server {
	return &Server{
		rooms:       make(map[string]*Room),
		CommandChan: make(chan Command),
	}
}

func (server *Server) Run() {
	for command := range server.CommandChan {
		client := command.client

		switch command.id {
		case CMD_COMMANDS:
			server.getAvailableCommands(client)
		case CMD_USERNAME:
			server.setUsername(client, command.args)
		case CMD_JOIN:
			server.joinRoom(client, command.args)
		case CMD_MSG:
			server.sendMessage(client, command.args)
		case CMD_DM:
			server.sendDirectMessage(client, command.args)
		case CMD_ROOM:
			server.getCurrentRoom(client)
		case CMD_MEMBERS:
			server.getMembersInARoom(client)
		case CMD_ROOMS:
			server.getAvailableRooms(client)
		case CMD_QUIT:
			server.quit(client)
		case CMD_HELP:
			server.help(client, command.args)
		default:
			client.err(fmt.Errorf("unknown command, to see available commands use /commands"))
		}
	}
}

func (server *Server) getAvailableCommands(client *Client) {
	client.msg(fmt.Sprintf("available commands:\n%s", strings.Join(Commands, "\n")))
	client.msg("see how can a command be used with /help <command> e.g. /help /dm")
}

func (server *Server) setUsername(client *Client, args []string) {
	if len(args) == 0 {
		client.err(fmt.Errorf("your username can't be empty"))

		return
	}

	username := args[0]

	client.username = username
	client.msg(fmt.Sprintf("you've changed your username to %s", username))
}

func (server *Server) joinRoom(client *Client, args []string) {
	roomName := args[0]

	room, ok := server.rooms[roomName]

	if !ok {
		room = CreateNewRoom(roomName)
		server.rooms[room.name] = room
	}

	room.members[client.conn.RemoteAddr()] = client

	if client.currentRoom != nil {
		server.quitCurrentRoom(client)
	}

	client.currentRoom = room

	client.msg(fmt.Sprintf("welcome to %s!", room.name))
	room.broadcast(client, fmt.Sprintf("%s has joined the room!", client.username))
}

func (server *Server) sendMessage(client *Client, args []string) {
	message := strings.Join(args, " ")

	if len(message) == 0 {
		return
	}

	client.currentRoom.broadcast(client, fmt.Sprintf("%s: %s", client.username, message))
}

func (server *Server) sendDirectMessage(sender *Client, args []string) {
	if sender.currentRoom == nil {
		sender.err(fmt.Errorf("you must join a room first"))

		return
	}

	if len(args) == 0 {
		sender.err(fmt.Errorf("/dm needs a recipient and a message"))

		return
	}

	username := args[0]
	message := strings.Join(args[1:], " ")

	if len(message) == 0 {
		return
	}

	for _, client := range sender.currentRoom.members {
		if username == client.username {
			client.msg(fmt.Sprintf("(dm) %s: %s", sender.username, message))

			return
		}
	}

	sender.err(fmt.Errorf("username %s not found in %s", username, sender.currentRoom.name))
}

func (sever *Server) getCurrentRoom(client *Client) {
	if client.currentRoom == nil {
		client.err(fmt.Errorf("you are not in a room"))

		return
	}

	client.msg(fmt.Sprintf("you are currently in %s", client.currentRoom.name))
}

func (server *Server) getMembersInARoom(client *Client) {
	if client.currentRoom == nil {
		client.err(fmt.Errorf("you must be in a room to see its members"))

		return
	}

	var members []string

	for _, member := range client.currentRoom.members {
		if client.username == member.username {
			members = append(members, fmt.Sprintf("(you) %s", member.username))

			continue
		}

		members = append(members, member.username)
	}

	client.msg(fmt.Sprintf("members in %s: %s", client.currentRoom.name, strings.Join(members, ", ")))
}

func (server *Server) getAvailableRooms(client *Client) {
	if len(server.rooms) == 0 {
		client.msg("no rooms available, feel free to create a new one with /join <room_name>")

		return
	}

	var roomNames []string

	for roomName := range server.rooms {
		roomNames = append(roomNames, roomName)
	}

	client.msg(fmt.Sprintf("available rooms are: %s", strings.Join(roomNames, ", ")))
}

func (server *Server) quit(client *Client) {
	if client.currentRoom == nil {
		client.err(fmt.Errorf("you are not in a room"))

		return
	}

	client.msg("sad to see you go :(")

	server.quitCurrentRoom(client)
}

func (server *Server) help(client *Client, args []string) {
	if len(args) == 0 {
		client.err(fmt.Errorf("a command name is needed for /help"))

		return
	}

	command := args[0]

	switch command {
	case "/commands":
		client.msg("</commands> is used to see all available commands")
	case "/username":
		client.msg("</username> flynn is used to set your username to 'flynn'")
	case "/join":
		client.msg("</join> #general is used to join '#general' or create it if not exist")
	case "/msg":
		client.msg("</msg> hello friend is used to send 'hello friend' to all members in the room you're in")
	case "/dm":
		client.msg("</dm> sam hello there! is used to send a dm 'hello there!' to sam")
	case "/room":
		client.msg("</room> is used to see what room you're currently in")
	case "/members":
		client.msg("</members> is used to see all members in the room you're in")
	case "/rooms":
		client.msg("</rooms> is used to see all rooms on the server")
	case "/quit":
		client.msg("</quit> is used to quit the room you're in")
	default:
		client.msg("unknown command")
	}
}

func (server *Server) quitCurrentRoom(client *Client) {
	delete(client.currentRoom.members, client.conn.RemoteAddr())

	client.currentRoom.broadcast(client, fmt.Sprintf("%s has left the room!", client.username))
	client.currentRoom = nil
}
