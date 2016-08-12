package main

import (
	"fmt"
	"regexp"
	"strings"
)

type Router struct {
	routes   []string
	handlers []func()
	names    [][]string
}

var patternRegexp = regexp.MustCompile(`\{([^\}:]+)(?::([^\}]+))?\}`)

func (r *Router) addRoute(pattern string, fn func()) {
	parameters := patternRegexp.FindAllStringSubmatch(pattern, 10)

	names := []string{}

	for _, parameter := range parameters {
		name := parameter[1]
		names = append(names, name)
	}

	r.routes = append(r.routes, patternRegexp.ReplaceAllStringFunc(pattern, func(s string) string {
		foos := strings.SplitN(s, ":", 2)
		if len(foos) < 2 {
			return `([^/]+)`
		} else {
			return "(" + foos[1][0:len(foos[1])-1] + ")"
		}
	}))

	for i, _ := range names {
		if i == 0 {
			r.handlers = append(r.handlers, fn)
		} else {
			r.handlers = append(r.handlers, nil)
		}
	}

	r.names = append(r.names, names)
}

func (r *Router) getRegexp() *regexp.Regexp {
	return regexp.MustCompile(`\A(?:` + strings.Join(r.routes, "|") + `)\z`)
}

func (r *Router) route(path string) {
	matches := r.getRegexp().FindAllStringSubmatch(path, 2)[0][1:]
	i := 0
	for _, match := range matches {
		if len(match) != 0 {
			break
		} else {
			i++
		}
	}
	r.handlers[i]()
}

func main() {
	r := &Router{}
	r.addRoute(`/user/{name}/{id:\d+}`, func() {
		fmt.Println("route 1")
	})
	r.addRoute(`/user/{id:\d+}`, func() {
		fmt.Println("route 2")
	})
	r.addRoute(`/user/{name}`, func() {
		fmt.Println("route 3")
	})
	r.route("/user/a23")
}
