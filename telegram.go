package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strconv"
	"strings"
)

var bot *tgbotapi.BotAPI

func telegram() {
	var err error

	bot, err = tgbotapi.NewBotAPI(conf.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil { // If we got a message
			log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)
			if update.Message.From.ID != 5192986817 {
				return
			}
			delim := strings.Split(update.Message.Text, " ")

			switch delim[0][1:] {
			case "matches":
				bot.Send(matches(update))
			case "get":
				bot.Send(get(update, delim))
			case "add":
				bot.Send(add(update, delim))
			case "clean":
				bot.Send(clean(update, delim))
			default:
				bot.Send(generateMessage(update, "🥵 <strong>Неизвестная команда</strong>"))
			}
		}
	}
}

func clean(update tgbotapi.Update, delim []string) tgbotapi.MessageConfig {
	if len(delim) < 2 {
		return generateMessage(update, "<strong>📚 Недостаточно аргументов</strong>")
	}

	removeMatches(strings.Join(delim[1:], " "))
	return generateMessage(update, "+")
}

func add(update tgbotapi.Update, delim []string) tgbotapi.MessageConfig {
	if len(delim) < 3 {
		return generateMessage(update, "<strong>📚 Недостаточно аргументов</strong>")
	}

	game := delim[1]
	players := getPlayersFromNames(delim[2:])

	for _, player := range players {
		addToBlockList(connection, player.Id, game, int(update.Message.From.ID))
	}

	return generateMessage(update, "<strong>😘 Игроки добавлены в ЧС</strong>")
}

func get(update tgbotapi.Update, delim []string) tgbotapi.MessageConfig {
	if len(delim) < 2 {
		return generateMessage(update, "<strong>📚 Недостаточно аргументов</strong>")
	}

	matches := getMatches(strings.Join(delim[1:], " "))
	if len(matches) < 1 {
		return generateMessage(update, "<strong>📚 Матчей нет</strong>")
	}

	playersAPI := getPlayers(matches[0].Players)
	if playersAPI == nil {
		return generateMessage(update, "<strong>🪄 Игроки свалили</strong>")
	}

	var winners []string
	var players []string
	var games []string

	var winnersStats = make(map[string][]string)
	var playersStats []string

	for _, match := range matches {
		var winnersCurrent []string
		for _, winner := range strings.Split(match.Winner, ";") {
			if winner == "" {
				continue
			}
			wn, err := strconv.Atoi(winner)
			if err != nil {
				log.Fatal(winner + " no converted to int")
			}
			winnersCurrent = append(winnersCurrent, playersAPI[wn].Username)
			winnersStats[playersAPI[wn].Username] = append(winnersStats[playersAPI[wn].Username], fmt.Sprintf("%s {%s} %s https://vimetop.ru/matches#%s", match.Game, convertDate(match.Date), winnersCurrent, match.MatchId))

			if len(winners) < 1 {
				winners = append(winners, playersAPI[wn].Username)
			}
		}

		if len(players) < 1 {
			for _, player := range playersAPI {
				playersStats = append(playersStats, player.Username)
				if find(winners, player.Username) == false {
					players = append(players, player.Username)
				}
			}
		}
		games = append(games,
			fmt.Sprintf("%s {%s} %s https://vimetop.ru/matches#%s", match.Game, convertDate(match.Date), winnersCurrent, match.MatchId))
	}

	for username, player := range winnersStats {
		length := len(player)

		if length < 9 {
			continue
		}

		gameList := player
		if length > 20 {
			gameList = gameList[0:20]
		}
		//TODO: пофиксить тут дату
		bot.Send(generateMessage(
			update,
			fmt.Sprintf("%s\n%s\nCharkosOff\n%s\n%s",
				username,
				strings.Join(remove(playersStats, username), " "),
				convertDate(matches[0].Date),
				strings.Join(gameList, "\n"))))
	}

	//TODO: пофиксить тут дату
	return generateMessage(update, fmt.Sprintf("%s\n%s\nCharkosOff\n%s\n%s",
		strings.Join(winners, " "),
		strings.Join(players, " "),
		convertDate(matches[0].Date),
		strings.Join(games, "\n")))
}

func matches(update tgbotapi.Update) tgbotapi.MessageConfig {
	matches := getSuspectMatches(9)
	if matches == nil {
		return generateMessage(update, "<strong>🕵️‍♀️ Нет подозрительных матчей</strong>")
	}

	var list []string
	for _, match := range matches {
		list = append(list, fmt.Sprintf("{%v} <code>%s</code>", match.Repetitions, match.Players))
	}

	return generateMessage(update, "<strong>💩 Возможные бустеры:</strong>\n"+strings.Join(list, "\n"))
}

func generateMessage(update tgbotapi.Update, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	return msg
}
