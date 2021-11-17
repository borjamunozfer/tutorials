Feature: Players

Scenario: There are players
Given there are no players
When I register players
Then I got 6 players
