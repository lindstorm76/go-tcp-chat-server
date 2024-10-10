package internal

import "net"

type Room struct {
	name    string
	members map[net.Addr]*Client
}

func (room *Room) broadcast(sender *Client, message string) {
	for clientAddr, client := range room.members {
		if clientAddr != sender.conn.RemoteAddr() {
			client.msg(message)
		}
	}
}

func CreateNewRoom(name string) *Room {
	return &Room{
		name:    name,
		members: make(map[net.Addr]*Client),
	}
}
