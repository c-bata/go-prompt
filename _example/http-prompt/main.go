package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"

	prompt "github.com/c-bata/go-prompt"
)

type RequestContext struct {
	url    *url.URL
	header http.Header
	client *http.Client
}

var ctx *RequestContext

// See https://github.com/eliangcs/http-prompt/blob/master/http_prompt/completion.py
var suggestions = []prompt.Suggest{
	// Command
	{"cd", "Change URL/path"},
	{"exit", "Exit http-prompt"},

	// HTTP Method
	{"delete", "DELETE request"},
	{"get", "GET request"},
	{"patch", "GET request"},
	{"post", "POST request"},
	{"put", "PUT request"},

	// HTTP Header
	{"Accept", "Acceptable response media type"},
	{"Accept-Charset", "Acceptable response charsets"},
	{"Accept-Encoding", "Acceptable response content codings"},
	{"Accept-Language", "Preferred natural languages in response"},
	{"ALPN", "Application-layer protocol negotiation to use"},
	{"Alt-Used", "Alternative host in use"},
	{"Authorization", "Authentication information"},
	{"Cache-Control", "Directives for caches"},
	{"Connection", "Connection options"},
	{"Content-Encoding", "Content codings"},
	{"Content-Language", "Natural languages for content"},
	{"Content-Length", "Anticipated size for payload body"},
	{"Content-Location", "Where content was obtained"},
	{"Content-MD5", "Base64-encoded MD5 sum of content"},
	{"Content-Type", "Content media type"},
	{"Cookie", "Stored cookies"},
	{"Date", "Datetime when message was originated"},
	{"Depth", "Applied only to resource or its members"},
	{"DNT", "Do not track user"},
	{"Expect", "Expected behaviors supported by server"},
	{"Forwarded", "Proxies involved"},
	{"From", "Sender email address"},
	{"Host", "Target URI"},
	{"HTTP2-Settings", "HTTP/2 connection parameters"},
	{"If", "Request condition on state tokens and ETags"},
	{"If-Match", "Request condition on target resource"},
	{"If-Modified-Since", "Request condition on modification date"},
	{"If-None-Match", "Request condition on target resource"},
	{"If-Range", "Request condition on Range"},
	{"If-Schedule-Tag-Match", "Request condition on Schedule-Tag"},
	{"If-Unmodified-Since", "Request condition on modification date"},
	{"Max-Forwards", "Max number of times forwarded by proxies"},
	{"MIME-Version", "Version of MIME protocol"},
	{"Origin", "Origin(s} issuing the request"},
	{"Pragma", "Implementation-specific directives"},
	{"Prefer", "Preferred server behaviors"},
	{"Proxy-Authorization", "Proxy authorization credentials"},
	{"Proxy-Connection", "Proxy connection options"},
	{"Range", "Request transfer of only part of data"},
	{"Referer", "Previous web page"},
	{"TE", "Transfer codings willing to accept"},
	{"Transfer-Encoding", "Transfer codings applied to payload body"},
	{"Upgrade", "Invite server to upgrade to another protocol"},
	{"User-Agent", "User agent string"},
	{"Via", "Intermediate proxies"},
	{"Warning", "Possible incorrectness with payload body"},
	{"WWW-Authenticate", "Authentication scheme"},
	{"X-Csrf-Token", "Prevent cross-site request forgery"},
	{"X-CSRFToken", "Prevent cross-site request forgery"},
	{"X-Forwarded-For", "Originating client IP address"},
	{"X-Forwarded-Host", "Original host requested by client"},
	{"X-Forwarded-Proto", "Originating protocol"},
	{"X-Http-Method-Override", "Request method override"},
	{"X-Requested-With", "Used to identify Ajax requests"},
	{"X-XSRF-TOKEN", "Prevent cross-site request forgery"},
}

func livePrefix() (string, bool) {
	if ctx.url.Path == "/" {
		return "", false
	}
	return ctx.url.String() + "> ", true
}

func executor(in string) {
	in = strings.TrimSpace(in)

	var method, body string
	blocks := strings.Split(in, " ")
	switch blocks[0] {
	case "exit":
		fmt.Println("Bye!")
		os.Exit(0)
	case "cd":
		if len(blocks) < 2 {
			ctx.url.Path = "/"
		} else {
			ctx.url.Path = path.Join(ctx.url.Path, blocks[1])
		}
		return
	case "get", "delete":
		method = strings.ToUpper(blocks[0])
	case "post", "put", "patch":
		if len(blocks) < 2 {
			fmt.Println("please set request body.")
			return
		}
		body = strings.Join(blocks[1:], " ")
		method = strings.ToUpper(blocks[0])
	}
	if method != "" {
		req, err := http.NewRequest(method, ctx.url.String(), strings.NewReader(body))
		if err != nil {
			fmt.Println("err: " + err.Error())
			return
		}
		req.Header = ctx.header
		res, err := ctx.client.Do(req)
		if err != nil {
			fmt.Println("err: " + err.Error())
			return
		}
		result, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println("err: " + err.Error())
			return
		}
		fmt.Printf("%s\n", result)
		ctx.header = http.Header{}
		return
	}

	if h := strings.Split(in, ":"); len(h) == 2 {
		// Handling HTTP Header
		ctx.header.Add(strings.TrimSpace(h[0]), strings.Trim(h[1], ` '"`))
	} else {
		fmt.Println("Sorry, I don't understand.")
	}
}

func completer(in prompt.Document) []prompt.Suggest {
	w := in.GetWordBeforeCursor()
	if w == "" {
		return []prompt.Suggest{}
	}
	return prompt.FilterHasPrefix(suggestions, w, true)
}

func main() {
	var baseURL = "http://localhost:8000/"
	if len(os.Args) == 2 {
		baseURL = os.Args[1]
		if strings.HasSuffix(baseURL, "/") {
			baseURL += "/"
		}
	}
	u, err := url.Parse(baseURL)
	if err != nil {
		log.Fatal(err)
	}
	ctx = &RequestContext{
		url:    u,
		header: http.Header{},
		client: &http.Client{},
	}

	p := prompt.New(
		executor,
		completer,
		prompt.OptionPrefix(u.String()+"> "),
		prompt.OptionLivePrefix(livePrefix),
		prompt.OptionTitle("http-prompt"),
	)
	p.Run()
}
