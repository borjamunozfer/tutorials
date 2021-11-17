package players

import "fmt"

type Player struct {
	Name string
}

//var player1 = Player{Name: "Player1"}
var players = []Player{}

//strconv
func RegisterPlayers() error {

	for i := 0; i < 5; i++ {
		player := Player{Name: fmt.Sprintf("Player %d", i)}
		players = append(players, player)
	}
	return nil
}

func GetPlayers() []Player {
	return players
}

func CountPlayers() int {
	return len(players)
}
