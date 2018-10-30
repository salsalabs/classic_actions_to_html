package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"regexp"
	"strings"
	"sync"
	"text/template"
	"time"

	"github.com/salsalabs/godig"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
)

const html = "html"

const actionTemplate = `
<!DOCTYPE html>
<html>
  <head>
    <title>{{.Title}}</title>
  </head>
  <body>
    <div>
      <h1>{{.Title}}</title>
    </div>
    <div>
      {{.Description}}
    </div>
  </body>
</html>
`

//action is a targeted/blind/MCTA/petition action object from Salsa's database.
//Created 30-Oct-2018 09:04:08 by schema-maker (github.com/salsalabs/godig/cmd/schema-maker/main.go)
type action struct {
	ActionKey     string `json:"action_KEY"`
	DateCreated   string `json:"Date_Created"`
	LastModified  string `json:"Last_Modified"`
	ReferenceName string `json:"Reference_Name"`
	Title         string `json:"Title"`
	Description   string `json:"Description"`
}

//exists returns true if the specified file exists.
func exists(f string) bool {
	_, err := os.Stat(f)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatalf("%v %v\n", err, f)
	}
	return true
}

//proc accepts actions from the input queue and handles them.
func proc(in chan action) {
	for {
		b, ok := <-in
		if !ok {
			break
		}
		err := handle(b)
		if err != nil {
			log.Printf("proc: key %v, %v\n", b.ActionKey, err)
		}
	}
}

//filename parses a action and returns a filename with the specified
//extension.
func filename(b action, ext string) string {
	const form = "Mon Jan 02 2006 15:04:05 GMT-0700 (MST)"
	x := b.DateCreated
	s := ""
	d := "Unknown"
	if len(x) == 0 {
		x = b.LastModified
	}
	if len(x) != 0 {
		t, _ := time.Parse(form, x)
		d = t.Format("2006-01-02")
		s = strings.Replace(b.Title, "/", " ", -1)
		if len(s) == 0 {
			s = strings.Replace(b.ReferenceName, "/", " ", -1)
		}
		r, err := regexp.Compile("<.+?>")
		if err != nil {
			panic(err)
		}
		s = r.ReplaceAllString(s, "")
	}
	if len(s) == 0 {
		s = "Unknown"
	}
	s = strings.TrimSpace(s)
	return fmt.Sprintf("%v - %v - %v.%v", d, b.ActionKey, s, ext)
}

//handle accepts a action and writes the description as HTML.
func handle(b action) error {
	fn := filename(b, "html")
	fn = path.Join(html, fn)
	if exists(fn) {
		log.Printf("%s: HTML already exists\n", b.ActionKey)
		return nil
	}
	f, err := os.Create(fn)
	if err != nil {
		return err
	}

	var t = template.Must(template.New("t").Parse(actionTemplate))
	err = t.Execute(f, &b)
	log.Printf("wrote %s\n", fn)
	return nil
}

//push reads the action table and pushes action onto a queue.
func push(api *godig.API, summary bool, in chan action) error {
	t := api.NewTable("action")
	offset := int32(0)
	c := 500
	for c != 0 {
		var a []action
		err := t.Many(offset, c, "", &a)
		if err != nil {
			return err
		}
		log.Printf("Read %v records from offset %v\n", len(a), offset)
		c = len(a)
		offset += int32(c)
		for _, b := range a {
			if summary {
				fmt.Println(filename(b, "html"))
			} else {
				in <- b
			}
		}
	}
	close(in)
	return nil
}

//scrub handles the cases where resource URLs are on domains that Salsa no
//longer supports.
func scrub(x string) string {
	s := strings.Replace(x, "org2.democracyinaction.org", "org2.salsalabs.com", -1)
	s = strings.Replace(s, "salsa.democracyinaction.org", "org.salsalabs.com", -1)
	s = strings.Replace(s, "hq.demaction.org", "org.salsalabs.com", -1)
	s = strings.Replace(s, "cid:", "https:", -1)
	return s
}

//main accepts inputs form the user and processes actions into HTML files.
func main() {
	var (
		app     = kingpin.New("classic_actions_to_html", "A command-line app to read actions, correct DIA URLs and write contents as HTML.")
		login   = app.Flag("login", "YAML file with login credentials").Required().String()
		summary = app.Flag("summary", "Show action dates, keys and titles.  Does not write HTML").Default("false").Bool()
	)
	app.Parse(os.Args[1:])
	api, err := (godig.YAMLAuth(*login))
	if err != nil {
		log.Fatalf("%v\n", err)
	}
	if !exists(html) {
		err := os.Mkdir(html, os.ModePerm)
		if err != nil && !os.IsExist(err) {
			log.Fatalf("%v, %v\n", err, html)
		}
	}

	var wg sync.WaitGroup
	in := make(chan action, 100)

	go func(in chan action, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		proc(in)
	}(in, &wg)

	go (func(api *godig.API, summary bool, in chan action, wg *sync.WaitGroup) {
		wg.Add(1)
		defer wg.Done()
		err = push(api, summary, in)
		if err != nil {
			log.Fatalf("%v\n", err)
		}
	})(api, *summary, in, &wg)

	//Settle time.
	time.Sleep(10000)
	wg.Wait()
}
