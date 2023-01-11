package main

import (
	"database/sql"
	"log"
)

func addMatch(connection *sql.DB, id string, players string, game string, winner string) {
	_, err := connection.Exec("INSERT INTO `matches`(`match_id`, `players`, `date`, `game`, `winner`) VALUES (?, ?, ?, ?, ?)", id, players, getCurrentDate(), game, winner)
	if err != nil {
		log.Fatal(err)
	}
}

func getMatch(connection *sql.DB, id string) bool {
	var receivedMatch string
	err := connection.QueryRow("SELECT `game` FROM `matches` WHERE `match_id` = ?", id).Scan(&receivedMatch)

	return err == nil
}

func removeMatches(players string) {
	connection.QueryRow("DELETE FROM `matches` WHERE `players` = ?", players)
}

func getMatches(players string) []MatchSQL {
	var list []MatchSQL
	rows, err := connection.Query("SELECT * FROM `matches` WHERE `players` = ? LIMIT 30", players)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var (
			id      string
			players string
			date    int
			game    string
			winner  string
		)

		if err := rows.Scan(&id, &players, &date, &game, &winner); err != nil {
			log.Fatal(err)
		}
		list = append(list, MatchSQL{id, players, date, game, winner})
	}

	return list
}

func getSuspectMatches(repetitions int) []MatchListSQL {
	var list []MatchListSQL
	rows, err := connection.Query("SELECT COUNT(*) as repetitions, players FROM `matches` GROUP BY players HAVING repetitions > ?", repetitions)
	if err != nil {
		log.Fatal(err)
	}

	for rows.Next() {
		var (
			count   int16
			players string
		)
		if err := rows.Scan(&count, &players); err != nil {
			log.Fatal(err)
		}
		list = append(list, MatchListSQL{count, players})
	}

	return list
}

func addToBlockList(connection *sql.DB, id int, game string, admin int) {
	_, err := connection.Exec("INSERT INTO `blocked`(`user_id`, `game`, `admin`, `date`) VALUES (?, ?, ?, ?)", id, game, admin, getCurrentDate())
	if err != nil {
		log.Fatal(err)
	}
}

func checkBlockStatus(connection *sql.DB, id int, game string) bool {
	var receivedGame string
	err := connection.QueryRow("SELECT `game` FROM `blocked` WHERE `user_id` = ? and `game` = ?", id, game).Scan(&receivedGame)

	return err == nil
}
