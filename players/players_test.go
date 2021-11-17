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

func iRegisterPlayers() error {
	RegisterPlayers()
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

func InitializeScenario(ctx *godog.ScenarioContext) {

	ctx.Step(`^I got (\d+) players$`, iGotPlayers)
	ctx.Step(`^I register players$`, iRegisterPlayers)
	ctx.Step(`^there are no players$`, thereAreNoPlayers)
}
