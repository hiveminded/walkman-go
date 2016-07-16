package main

import (
	"encoding/json"
	"fmt"
	"github.com/Sirupsen/logrus"
	"github.com/kataras/iris"
	"io/ioutil"
	"net/http"
	"strings"
	"github.com/iris-contrib/middleware/cors"
)

func main() {
	iris.Use(cors.Default()) // crs
	iris.Get("/songs/_search", GetSongs)
	iris.Listen(":8080")
}

func GetSongs(c *iris.Context) {
	bs, err := json.Marshal(c.URLParams())
	assertNil(err)

	reader := strings.NewReader(string(bs))
	request, err := http.NewRequest("GET", "http://localhost:9200/songs/_search?pretty=false", reader)
	assertNil(err)
	client := &http.Client{}
	resp, err := client.Do(request)
	assertNil(err)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	assertNil(err)

	var result map[string]interface{}
	assertNil(json.Unmarshal(body, &result))

	c.JSON(resp.StatusCode, result)
}

// assertNil panic if the test is not nil
func assertNil(test interface{}, params ...interface{}) {
	if test != nil {
		f := logrus.Fields{}
		for i := range params {
			f[fmt.Sprintf("param%d", i)] = params[i]
		}

		if e, ok := test.(error); ok {
			logrus.WithFields(f).Panic(e)
		}
		logrus.WithFields(f).Panic("must be nil, but its not")
	}
}
