package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const ENDPOINT = "https://api.vimeworld.com/"

func req(URL string) []byte {
	resp, err := http.Get(URL)
	if err != nil {
		log.Fatal(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	return body
}

func getMatchInfo(ID string) Match {
	body := req(fmt.Sprintf(ENDPOINT+"match/%v", ID))
	matchInfo := Match{}

	err := json.Unmarshal(body, &matchInfo)
	if err != nil {
		log.Print(string(body[:]))
		log.Fatal(err)
	}

	return matchInfo
}

func getMatchWinner(match Match) []string {
	var players []string
	winner := match.Winner["player"]
	if winner != nil {
		players = append(players, strconv.Itoa(int(winner.(float64))))
	}

	winner = match.Winner["players"]
	if winner != nil {
		winners := winner.([]interface{})
		for i := 0; i < len(winners); i++ {
			players = append(players, strconv.Itoa(int(winners[i].(float64))))
		}
	}

	winner = match.Winner["team"]
	if winner != nil {
		//log.Print(winner.(string))
	}

	return players
}

func getMatchLatest() []MatchPreview {
	body := req(ENDPOINT + "match/latest")
	var matchList []MatchPreview

	err := json.Unmarshal(body, &matchList)
	if err != nil {
		log.Fatal(err)
	}

	return matchList
}

func getPlayersFromNames(players []string) map[int]Player {
	var list = make(map[int]Player)

	body := req(ENDPOINT + "user/name/" + strings.Join(players, ","))
	var playersList []Player
	err := json.Unmarshal(body, &playersList)

	if err != nil {
		log.Fatal(err)
	}

	for _, player := range playersList {
		list[player.Id] = player
	}

	return list
}

func getPlayers(players string) map[int]Player {
	var list = make(map[int]Player)

	body := req(ENDPOINT + "user/" + strings.Join(strings.Split(players, " "), ","))
	var playersList []Player
	err := json.Unmarshal(body, &playersList)

	if err != nil {
		log.Fatal(err)
	}

	for _, player := range playersList {
		list[player.Id] = player
	}

	return list
}
