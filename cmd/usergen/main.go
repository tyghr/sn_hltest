package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/bxcodec/faker/v3"
	// _ "github.com/go-sql-driver/mysql"
)

type fakeUser struct {
	UserName  string `db:"username" json:"username" faker:"username"`
	Password  string `faker:"password"`
	Name      string `db:"name" json:"name" faker:"first_name"`
	SurName   string `db:"surname" json:"surname" faker:"last_name"`
	BirthDate string `db:"birthdate" json:"birthdate" faker:"date"`
	Gender    string `db:"gender" json:"gender" faker:"oneof: M, F"`
	City      string `db:"city" json:"city" faker:"timezone"`
	Interest  string `json:"interests" faker:"word"`
}

var (
	httpTimeOut = time.Second * 10
)

func main() {
	client := &http.Client{
		Timeout: httpTimeOut,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}
	defer client.CloseIdleConnections()

	wg := sync.WaitGroup{}
	taskCh := make(chan int, 10)

	start := time.Now()

	for worker := 0; worker < 20; worker++ {
		go func() {
			for range taskCh {
				err := createUser(client)
				if err != nil {
					log.Println(err)
				}
				wg.Done()
			}
		}()
	}

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		taskCh <- i
	}
	close(taskCh)
	wg.Wait()

	fmt.Printf("time elasped: %v\n", time.Since(start).Seconds())
}

func createUser(client *http.Client) error {
	registerURL := "http://127.0.0.1/register"

	u := fakeUser{}
	err := faker.FakeData(&u)
	if err != nil {
		return fmt.Errorf("failed creating fake user: %v", err)
	}

	data := url.Values{
		"username":        {u.UserName},
		"first_name":      {u.Name},
		"second_name":     {u.SurName},
		"gender":          {u.Gender},
		"birthdate":       {u.BirthDate},
		"city":            {u.City},
		"interests":       {u.Interest},
		"password":        {u.Password},
		"repeat_password": {u.Password},
	}

	resp, err := client.PostForm(registerURL, data)
	if err != nil {
		return fmt.Errorf("failed post request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed post request (status): %v", err)
	}

	return nil
}
