package internal

const (
	CMD_DEFAULT = iota
	CMD_COMMANDS
	CMD_USERNAME
	CMD_JOIN
	CMD_MSG
	CMD_DM
	CMD_ROOM
	CMD_MEMBERS
	CMD_ROOMS
	CMD_QUIT
	CMD_HELP
)

var Commands = []string{
	"/commands",
	"/username",
	"/join",
	"/msg",
	"/dm",
	"/room",
	"/members",
	"/rooms",
	"/quit",
	"/help",
}

type Command struct {
	id     int
	client *Client
	args   []string
}
