package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.WrapRouter(cors.WrapNext(cors.Options{})) // crs

	app.Get("/songs/_search", GetSongs)
	app.Run(iris.Addr(":8080"))
}

func GetSongs(c iris.Context) {
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

	c.StatusCode(resp.StatusCode)
	c.JSON(result)
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
