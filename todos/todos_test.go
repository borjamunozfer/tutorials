package todos

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"

	"github.com/cucumber/godog"
)

type TodoFeature struct {
	todoClient *todoClient
	response   *http.Response
	body       []byte
}

func (tf *TodoFeature) iSendRequestTo(ctx context.Context, method, endpoint string) error {

	//set path to URL struct
	tf.todoClient.baseURL.Path = endpoint
	resp := &http.Response{}
	var err error
	if method == "GET" && endpoint == "posts/1" {
		resp, err = tf.todoClient.GetPost()
		if err != nil {
			return err
		}

	} else if method == "GET" {
		resp, err = tf.todoClient.GetPosts()
		if err != nil {
			return err
		}
	}

	if method == "POST" {
		// Get body from Step context shared.
		// Maps return tuple internally with a bool if the Key (postBody) exists.
		if val, ok := ctx.Value("postBody").(map[string]interface{}); ok {
			// we found our map from the Context
			resp, err = tf.todoClient.CreatePost(val)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("Could not find by postBody key from Context")
		}
	}

	//Update struct TodoFeature with response obtained.
	tf.response = resp
	return nil
}

func (tf *TodoFeature) theResponseShouldBeNotEmpty() error {

	if len(tf.body) == 0 {
		return fmt.Errorf("Response body is empty")
	}
	return nil
}

func (tf *TodoFeature) theResponseShouldMatchJson(bodyjson *godog.DocString) (err error) {
	// we need to compare our expected json response with our resp (which internally represents a Post struct)
	// our response from the previous step is stored in the struct http.Response

	/*defer guarantees that this sentence executes lastly and always before the function execution ends (in this case, our theResponseShouldMatchJson)
	  Body is internally an stream object, so we need to close the flux always and free the resource to avoid memory leaks in case we get any error
	*/
	defer tf.response.Body.Close()
	var expected, actual Post

	body, err := io.ReadAll(tf.response.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal([]byte(bodyjson.Content), &expected); err != nil {
		return
	}
	if err = json.Unmarshal(body, &actual); err != nil {
		return
	}
	//util to compare structs, maps, etc.
	if !reflect.DeepEqual(expected, actual) {
		return fmt.Errorf("got %v want %v\n", expected, actual)
	}

	return nil
}

func (tf *TodoFeature) theResponseStatusShouldBe(status string) (err error) {

	if tf.response.Status != status {
		return fmt.Errorf("Got %s http status, but expected %s", tf.response.Status, status)
	}
	return nil
}

func (tf *TodoFeature) thereArePostsCreated() (err error) {

	//get response to see if *resp is not empty.
	var posts []Post

	defer tf.response.Body.Close()

	body, err := io.ReadAll(tf.response.Body)
	if err != nil {
		return err
	}

	//save body to our struct TodoFeature just in case we need
	tf.body = body

	//BODY []bytes is a JSON  --> []Post
	if err = json.Unmarshal(body, &posts); err != nil {
		return
	}

	if len(posts) == 0 {
		return fmt.Errorf("There is no posts created")
	}

	return nil
}

func (tf *TodoFeature) validRequestBody(ctx context.Context) (context.Context, error) {

	newBodyPost := map[string]interface{}{
		"title":  "test title",
		"userId": 2021,
		"body":   "test body",
	}

	return context.WithValue(ctx, "postBody", newBodyPost), nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	// initialize structs.
	tc := &todoClient{baseURL: &url.URL{Host: "https://jsonplaceholder.typicode.com/"},
		client: &http.Client{}}
	tfeature := TodoFeature{todoClient: tc}

	ctx.Step(`^I send "([^"]*)" request to "([^"]*)"$`, tfeature.iSendRequestTo)
	ctx.Step(`^the response should be not empty$`, tfeature.theResponseShouldBeNotEmpty)
	ctx.Step(`^the response should match json:$`, tfeature.theResponseShouldMatchJson)
	ctx.Step(`^the response status should be "([^"]*)"$`, tfeature.theResponseStatusShouldBe)
	ctx.Step(`^there are posts created$`, tfeature.thereArePostsCreated)
	ctx.Step(`^valid request body$`, tfeature.validRequestBody)
}
