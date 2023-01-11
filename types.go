package main

type Config struct {
	User     string
	Database string
	Password string

	Token string
}

type Player struct {
	Id       int
	Username string
	Level    int
	Rank     string
}

type MatchSQL struct {
	MatchId string
	Players string
	Date    int
	Game    string
	Winner  string
}

type MatchListSQL struct {
	Repetitions int16
	Players     string
}

type Match struct {
	Version int
	Game    string
	Server  string
	Owned   bool

	Start int
	End   int

	Winner  map[string]interface{}
	Players []MatchPlayer
	Events  []MatchEvent
}

type MatchPreview struct {
	Id       string
	Game     string
	Duration int
}

type MatchPlayer struct {
	Id int
	//Duels
	Kills     int
	WinStreak int
	Dead      bool

	//BlockParty
	Levels int
}

type MatchEvent struct {
	Type string
	Time int

	//Duels
	Killer       int
	Target       int
	KillerHealth string
}
