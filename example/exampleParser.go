package main

import (
	"fmt"
	"github.com/Wuvist/mustache"
	"github.com/airdispatch/dpl"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type Quick struct {
	Actions    []dpl.Action
	Content    template.HTML
	ActionName string
}

func (q Quick) CheckPath(p string) bool {
	return p == q.ActionName
}

func main() {
	xmlFile, err := os.Open("example/notes.dpl")
	if err != nil {
		fmt.Println("Couldn't open file.")
		return
	}
	defer xmlFile.Close()

	b, _ := ioutil.ReadAll(xmlFile)

	o, err := dpl.ParseDPL(b)
	if err != nil {
		fmt.Println(err)
		return
	}

	document := `
	<html>
		<head>
			<title>AirDispatch Plugin Test</title>
			<link href="//netdna.bootstrapcdn.com/bootstrap/3.0.2/css/bootstrap.min.css" rel="stylesheet">
			<meta name="viewport" content="initial-scale=1">
		</head>
		<body>
			<div class="container">
				<h1>AirDispatch Plugin Tester</h1>
				<div class="row">
					<div class="col-sm-2">
						<p>Action List</p>
						<div class="list-group">
							{{ $out := . }}
							{{ range .Actions }}
							<a class="list-group-item {{ if $out.CheckPath .Name }}active{{ end }}" href="/{{ .Name }}">{{ .Name }}</a>
							{{ end }}
						</div>
					</div>
					<div class="col-sm-10">
						<p>Content</p>
						{{ .Content }}
					</div>
				</div>
			</div>
		</body>
	</html>
	`
	t, _ := template.New("d").Parse(document)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := strings.Split(r.URL.Path, "/")
		identifier := -1
		var action string
		if len(path) == 2 && path[1] == "" {
			action = o.DefaultAction
		} else if len(path) == 2 {
			action = path[1]
		} else if len(path) == 3 {
			action = path[1]
			var err error
			identifier, err = strconv.Atoi(path[2])

			if err != nil {
				t.Execute(w, "Couldn't parse id value.")
				return
			}
		} else {
			t.Execute(w, "Couldn't determine action")
		}
		fmt.Println("Getting Request for", action, "with id", identifier)
		a := o.Actions[action]

		data := mustache.Render(a.HTML, getFakeParameters(o, identifier))

		t.Execute(w, Quick{o.Action, template.HTML(data), action})
	})
	// fmt.Println(o)
	fmt.Println("Ready to accept Connections.")
	log.Fatal(http.ListenAndServe(":2048", nil))
}

func getFakeParameters(p *dpl.Plugin, id int) map[string]interface{} {
	if id != -1 {
		return map[string]interface{}{
			"message": getMessageForId(p, id),
			"actions": getActions(p, ""),
		}
	} else {
		return map[string]interface{}{
			"messages": getMessages(p),
			"actions":  getActions(p, ""),
		}
	}
}

func getMessages(p *dpl.Plugin) []map[string]interface{} {
	return []map[string]interface{}{
		getMessageForId(p, 15),
	}
}

func getMessageForId(p *dpl.Plugin, id int) map[string]interface{} {
	return map[string]interface{}{
		"title":    "Hello, there.",
		"body":     "adsfasdfreg eropijg adf;j",
		"category": "asdf reqt ",
		"id":       id,
		"actions":  getActions(p, strconv.Itoa(id)),
	}
}

func getActions(p *dpl.Plugin, id string) map[string]interface{} {
	output := make(map[string]interface{})
	for _, v := range p.Actions {
		vmap := make(map[string]interface{})
		url := "/" + v.Name
		if id != "" {
			url += "/" + id
		}
		vmap["url"] = url
		vmap["link"] = template.HTML("<a href='" + url + "'>" + v.Name + "</a>")
		output[v.Name] = vmap
	}
	return output
}
