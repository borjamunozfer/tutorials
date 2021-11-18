Feature: Players


@wip
Scenario: There are players
Given there are no players
When I register 5 players
Then I got 5 players


@static
Scenario: Register 5 new players
Given there are 5 players
When I register 5 players
Then I got 10 players

Scenario: Register new players from given
Given there are <startPlayers> players
| startPlayers | newPlayers | endPlayers |
|  0           | 1          |  1         |
|  2           | 2          |  4         |
|  3           | 3          |  6         |
When I register players
Then I got players

@dynamic
Scenario Outline: Register new players
Given there are <startPlayers> players
When I register <newPlayers> players
Then I got <endPlayers> players
Examples:
    | startPlayers | newPlayers | endPlayers |
    | 0            |  5         | 5          |
    | 2            |  2         | 4          |
