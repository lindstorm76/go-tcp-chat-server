# Go TCP Chat Server

Made with ❤️ by Thanapong Angkha

## Usage

**Connect a client to the server**

For `telnet` use `telnet localhost 8080`

For `netcat` use `nc localhost 8080`

| Command                     | Description                                   |
| :---                        | :----:                                     |
| /commands                   | See all available commands                    |
| /username `username`        | Set client username                           |
| /join `room_name`           | Join a room                                   |
| /msg `message`...           | Send a message to the current room            |
| /dm `username` `message`... | Send a direct message to a member in the room |
| /room                       | See what room you're currently in             |
| /members                    | See all members in the room you're in         |
| /rooms                      | See all available rooms on the server         |
| /quit                       | Quit the current room                         |
| /help `command`             | See how a command can be used                 |

## How to run

1. Build the app

```
docker build -t go-tcp-chat-server .
```

2. Run it on any port you want, here I'll use 8080

```
docker run --name go-tcp-chat-server-container -p 8080:8080 go-tcp-chat-server
```
