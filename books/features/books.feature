Feature: Books feature
As an user
I want to be able to open, create, read, remove and write files

Scenario: Create empty book
Given the book "La Odisea" does not exist
When I create the book "La Odisea"
Then the file is created empty.

Scenario: Open book
Given the book "La Odisea" exists
When I open the book
Then the file is opened correctly

Scenario: Read book
Given the book "La Odisea" exists
When i read the book fully
Then the book content is returned

Scenario: Read book line by line
Given the book "La Odisea" exists
When I read the book by line
Then I got one line at time

Scenario: Write to existent book
Given the book "La Odisea" exists
When I write content 
Then the book content is updated

Scenario: Write to unexistent book
Given the book "La Odisea" does not exists
When I write content to unexistent book
Then the book is created and updated