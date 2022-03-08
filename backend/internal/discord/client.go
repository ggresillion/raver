package discord

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

const (
	api = "https://discordapp.com/api"
)

type DiscordApiError struct {
	Code    int
	Message string
}

func newDiscordApiError(code int, message string) *DiscordApiError {
	return &DiscordApiError{code, message}
}

func (e *DiscordApiError) Error() string {
	return fmt.Sprint("received error from discord api: ", e.Message)
}

type DiscordClient struct {
	token string
	mutex *sync.Mutex
}

var clients = map[string]*DiscordClient{}

func NewDiscordClient(token string) *DiscordClient {
	client := clients[token]
	if client != nil {
		return client
	}
	client = &DiscordClient{token, &sync.Mutex{}}
	clients[token] = client
	return client
}

// Make a request to the discord API
func (c *DiscordClient) request(verb string, path string) (io.Reader, error) {

	c.mutex.Lock()

	log.Println(fmt.Sprint("sending request to discord: ", verb, path))
	res, err := c.make(verb, path)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		return nil, newDiscordApiError(res.StatusCode, string(b))
	}

	c.handleRateLimit(res)

	return res.Body, nil
}

// Do the actual request
func (c *DiscordClient) make(verb string, path string) (*http.Response, error) {
	req, err := http.NewRequest(verb, fmt.Sprint(api, path), nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", fmt.Sprint("Bearer ", c.token))

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Handle rate limit from discord by waiting for X-RateLimit-Reset-After seconds if the bucket if empty
func (c *DiscordClient) handleRateLimit(res *http.Response) error {
	if res.Header.Get("X-RateLimit-Remaining") == "0" {
		secondsToWait, err := strconv.Atoi(res.Header.Get("X-RateLimit-Reset-After"))
		if err != nil {
			return err
		}

		go func() {
			defer c.mutex.Unlock()
			time.Sleep(time.Duration(secondsToWait) * time.Second)
		}()
		return nil
	}
	c.mutex.Unlock()
	return nil
}
