Feature: Author writer
As an Author
I want to be able to write Books

Scenario: Author create a new book
Given the book "La Odisea" does not exist
When I create the book "La Odisea"
Then a new empty book is created

Scenario: Author writes book
Given the book "La Odisea" exists
When I open the book "La Odisea"
And I write lines
Then the book "La Odisea" is updated

Scenario: Author delete line
Given the book "La Odisea" exists
When I open the book "La Odisea"
And I delete "last" line
Then the book line is removed

