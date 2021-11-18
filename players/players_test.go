package players

import (
	"fmt"

	"github.com/cucumber/godog"
)

func iGotPlayers(expectedPlayers int) error {
	gotPlayers := CountPlayers()

	if gotPlayers != expectedPlayers {
		return fmt.Errorf("Expected %d players but found %d", expectedPlayers, gotPlayers)
	}

	return nil
}

func iRegisterPlayers(newPlayers int) error {

	RegisterPlayers(newPlayers)
	return nil
}

func thereAreNoPlayers() error {

	totalPlayers := CountPlayers()
	fmt.Printf("There are %d players", totalPlayers)
	if totalPlayers != 0 {
		return fmt.Errorf("There are %d players registered yet", totalPlayers)
	}
	return nil
}

func thereArePlayers(startPlayers int) error {
	fmt.Printf("We have %d players from start\n", startPlayers)
	if startPlayers != 0 {
		RegisterPlayers(startPlayers)
	}
	return nil
}

func iGotPlayersv2() error {
	return godog.ErrPending
}

func iRegisterPlayersv2() error {
	return godog.ErrPending
}

func thereAreStartPlayersTable(playersTable *godog.Table) error {

	heads := playersTable.Rows[0]
	fmt.Println("Cabecera tabla:")
	for _, head := range heads.Cells {
		fmt.Printf("\t %s", head.Value)
	}
	for _, row := range playersTable.Rows[1:] {
		fmt.Println()
		for _, cell := range row.Cells {
			fmt.Printf("\t %s", cell.Value)
		}
	}
	return godog.ErrPending
}

func InitializeScenario(ctx *godog.ScenarioContext) {

	// Reset our slice of Players before run every Scenario.
	ctx.BeforeScenario(func(sc *godog.Scenario) {
		players = nil
	})

	ctx.Step(`^I got (\d+) players$`, iGotPlayers)
	ctx.Step(`^I register (\d*) players$`, iRegisterPlayers)
	ctx.Step(`^there are no players$`, thereAreNoPlayers)
	ctx.Step(`^there are (\d+) players$`, thereArePlayers)
	ctx.Step(`^I got players$`, iGotPlayers)
	ctx.Step(`^I register players$`, iRegisterPlayers)
	ctx.Step(`^there are <startPlayers> players$`, thereAreStartPlayersTable)

}
