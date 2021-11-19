Feature: Todos Feature
As an user
I want to be able to
create and query TODOs 

Scenario: Create posts
Given valid request body
When I send "POST" request to "posts"
Then the response status should be "201 Created"
And the response should match json:
"""
{
        "id": 101,
        "userId": 2021,
        "title": "test title",
        "body": "test body"
}        
"""

Scenario: Get all posts
When I send "GET" request to "posts"
And there are posts created
Then the response should be not empty

@testtag
Scenario: Get one post by id
When I send "GET" request to "posts/1"
Then the response should match json:
"""
{
    "userId": 1,
    "id": 1,
    "title": "sunt aut facere repellat provident occaecati excepturi optio reprehenderit",
    "body": "quia et suscipit\nsuscipit recusandae consequuntur expedita et cum\nreprehenderit molestiae ut ut quas totam\nnostrum rerum est autem sunt rem eveniet architecto"
}
"""

