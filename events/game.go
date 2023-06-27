package events

// EventGameMove is an event that is published when a player makes a move.
type EventGameMove struct {
	GameID      string
	Move        string
	Position    string
	PlayerWhite string
	PlayerBlack string
}

// EventGameStart is an event that is published when a game is started.
type EventGameStart struct {
	GameID      string
	PlayerBlack string
	PlayerWhite string
	Position    string
}
