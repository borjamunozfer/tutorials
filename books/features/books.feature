Feature: Books feature
As an user
I want to be able to open, create, read, remove and write files

Scenario: Remove an existent book
Given the book "La Odisea" exists
When I delete the book
Then the book "La Odisea" does not exists

Scenario: Create an empty book
Given the book "La Odisea" does not exist
When I create the book "La Odisea"
Then the file is created empty.

Scenario: Open a book
Given the book "La Odisea" exists
When I open the book
Then the file is opened correctly

Scenario: Open a book that does not exist
Given the book "La Odisea" does not exist
When I open the book non existent
Then the file is created

Scenario: Read a book 
Given the book "La Odisea" exists
And the book is not empty
When i read the book "La Odisea" fully
Then the book content is returned

Scenario: Read a book line by line
Given the book "La Odisea" exists
And the book is not empty
When I read the book by line
Then I got one line at time

Scenario: Write to an existent book
Given the book "La Odisea" exists
When I write content 
Then the book content is updated

Scenario: Write to an unexistent book
Given the book "La Odisea" does not exists
When I write content to unexistent book
Then the book is created and updated