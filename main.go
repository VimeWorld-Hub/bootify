package main

import (
	"database/sql"
	"flag"
	"github.com/BurntSushi/toml"
	"github.com/go-co-op/gocron"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	conf       Config
	connection *sql.DB
)

func main() {
	configFile := flag.String("config", "config.toml", "path to the config file")
	flag.Parse()

	var err error
	if _, err = toml.DecodeFile(*configFile, &conf); err != nil {
		log.Fatalln(err)
	}
	connection = createConnection(conf.User, conf.Password, conf.Database)
	go telegram()

	s := gocron.NewScheduler(time.UTC)
	_, err = s.Every("15s").Do(check)
	if err != nil {
		log.Fatal(err)
	}
	s.StartAsync()
	s.StartBlocking()
}

func check() {
	matches := getMatchLatest()
	for _, preview := range matches {
		//Отсеивание матчей, длина которых больше 1 минуты
		if preview.Duration > 60 {
			continue
		}
		//Отсеивание игр, на которых можно буститься
		if (preview.Game == "ZOMBIECLAUS") || (preview.Game == "WHITECOLD") {
			continue
		}

		//Получение полной инфы матча
		match := getMatchInfo(preview.Id)

		//Проверка, есть ли в матче игроки
		if len(match.Players) < 1 {
			continue
		}

		//Проверка, является ли матч приватным
		if match.Owned || strings.HasPrefix(match.Server, "OS") {
			continue
		}

		//Проверка, обработан ли был матч ранее
		if getMatch(connection, preview.Id) {
			continue
		}

		//Дополнительные проверки на дуэли
		continued := false

		if preview.Game == "DUELS" {
			for _, event := range match.Events {
				if event.Type == "kill" {
					killerHealth, err := strconv.ParseFloat(event.KillerHealth, 32)
					if err != nil {
						log.Print(killerHealth)
						log.Fatal(err)
					}

					if killerHealth < 17.0 {
						continued = true
						break
					}
				}
			}
		}

		winners := getMatchWinner(match)

		var players []string
		for _, player := range match.Players {
			if (checkBlockStatus(connection, player.Id, preview.Game)) && (find(winners, strconv.Itoa(player.Id))) {
				continued = true
				break
			}
			players = append(players, strconv.Itoa(player.Id))
		}

		//Если что-то нашлось. Костыль нужен, чтобы завершить не цикл перебора, а основной
		if continued {
			continue
		}

		sort.Strings(players)
		addMatch(connection, preview.Id, strings.Join(players, " "), preview.Game, strings.Join(winners, ";"))
	}
}

func createConnection(username string, password string, database string) *sql.DB {
	db, err := sql.Open("mysql", username+":"+password+"@/"+database)
	if err != nil {
		log.Fatal(err)
	}
	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db
}
