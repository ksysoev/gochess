package events

type EventGameMove struct {
	GameID      string
	Move        string
	Position    string
	PlayerWhite string
	PlayerBlack string
}

type EventGameStart struct {
	GameID      string
	PlayerBlack string
	PlayerWhite string
	Position    string
}
