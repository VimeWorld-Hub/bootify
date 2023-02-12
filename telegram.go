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

	bot, err = tgbotapi.NewBotAPI(conf.TelegramToken)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.From.ID != conf.Admin {
				return
			}
			delim := strings.Split(update.Message.Text, " ")

			switch delim[0][1:] {
			case "matches":
				sendMessage(matches(update))
			case "get":
				sendMessage(get(update, delim))
			case "add":
				sendMessage(add(update, delim))
			case "clean":
				sendMessage(clean(update, delim))
			default:
				sendMessage(generateMessage(update, "🥵 <strong>Неизвестная команда</strong>"))
			}
		}
	}
}

// Команда для удаления матчей с определёнными игрокам
func clean(update tgbotapi.Update, delim []string) tgbotapi.MessageConfig {
	if len(delim) < 2 {
		return generateMessage(update, "<strong>📚 /clean [id1] [id2] [id3]</strong>")
	}

	removeMatches(db, strings.Join(delim[1:], " "))
	return generateMessage(update, "<strong>🧹 Матчи, в которых участвовали эти игроки, удалены</strong>")
}

// Команда для добавления в ЧС (игнор-лист)
func add(update tgbotapi.Update, delim []string) tgbotapi.MessageConfig {
	if len(delim) < 3 {
		return generateMessage(update, "<strong>📚 /add [игра] [никнейм]</strong>")
	}

	game := delim[1]
	players := getPlayersFromNames(delim[2:])

	for _, player := range players {
		addToBlockList(db, player.Id, game, int(update.Message.From.ID))
	}

	return generateMessage(update, "<strong>🖤 Игроки добавлены в ЧС</strong>")
}

// Получение матчей с определёнными игроками
func get(update tgbotapi.Update, delim []string) tgbotapi.MessageConfig {
	if len(delim) < 2 {
		return generateMessage(update, "<strong>📚 /get [id1] [id2] [id3]</strong>")
	}

	matches := getMatches(db, strings.Join(delim[1:], " "))
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
			winnersStats[playersAPI[wn].Username] = append(winnersStats[playersAPI[wn].Username], fmt.Sprintf("%s {%s} %s https://vimetop.ru/matches#%s", match.Game, convertDateFull(match.Date), winnersCurrent, match.MatchId))

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
			fmt.Sprintf("%s {%s} %s https://vimetop.ru/matches#%s", match.Game, convertDateFull(match.Date), winnersCurrent, match.MatchId))
	}

	for username, player := range winnersStats {
		length := len(player)

		if length < 9 {
			continue
		}

		gameList := player
		if length > 25 {
			gameList = gameList[0:25]
		}

		sendMessage(generateMessage(update, generateReport(
			username,
			strings.Join(remove(playersStats, username), " "),
			convertDateMini(matches[0].Date),
			strings.Join(gameList, "\n"),
		)))
	}

	return generateMessage(update, generateReport(
		strings.Join(winners, " "),
		strings.Join(players, " "),
		convertDateMini(matches[0].Date),
		strings.Join(games, "\n"),
	))
}

// Получение всех подозрительных матчей
func matches(update tgbotapi.Update) tgbotapi.MessageConfig {
	matches := getSuspectMatches(db, 9)
	if matches == nil {
		return generateMessage(update, "<strong>😊 Нет подозрительных матчей</strong>")
	}

	var list []string
	for _, match := range matches {
		list = append(list, fmt.Sprintf("{%v} <code>%s</code>", match.Repetitions, match.Players))
	}

	return generateMessage(update, "<strong>🤡 Возможные бустеры:</strong>\n"+strings.Join(list, "\n"))
}

func generateReport(violators string, helpers string, date string, games string) string {
	return fmt.Sprintf("%s\n%s\n%s\n%s\n%s",
		violators,
		helpers,
		conf.VimeWorldUsername,
		date,
		games,
	)
}

// TODO: объединить функции
func sendMessage(messageConfig tgbotapi.MessageConfig) {
	_, err := bot.Send(messageConfig)
	if err != nil {
		log.Fatal(err)
	}
}

func generateMessage(update tgbotapi.Update, text string) tgbotapi.MessageConfig {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	return msg
}
