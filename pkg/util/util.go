package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path"

	jsoniter "github.com/json-iterator/go"
	http "github.com/valyala/fasthttp"
)

func getTokenPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	// create .telegraphcl directory if it doesn't exist
	if _, err := os.Stat(path.Join(usr.HomeDir, ".telegraphcl")); os.IsNotExist(err) {
		os.Mkdir(path.Join(usr.HomeDir, "/.telegraphcl"), 0755)
	}
	return path.Join(usr.HomeDir, ".telegraphcl", "telegraph.token")
}

var tokenPath string = getTokenPath()

type Response struct {
	Ok     bool            `json:"ok"`
	Error  string          `json:"error,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
}

func MakeRequest(path string, payload interface{}) ([]byte, error) {
	parser := jsoniter.ConfigFastest

	src, err := parser.Marshal(payload)
	if err != nil {
		return nil, errors.New("Failed to marshal payload")
	}

	u := http.AcquireURI()
	defer http.ReleaseURI(u)
	u.SetScheme("https")
	u.SetHost("api.telegra.ph")
	u.SetPath(path)

	req := http.AcquireRequest()
	defer http.ReleaseRequest(req)
	req.SetRequestURIBytes(u.FullURI())
	req.Header.SetMethod(http.MethodPost)
	req.Header.SetUserAgent("toby3d/telegraph")
	req.Header.SetContentType("application/json")
	req.SetBody(src)

	resp := http.AcquireResponse()
	defer http.ReleaseResponse(resp)

	if err := http.Do(req, resp); err != nil {
		return nil, err
	}

	r := new(Response)
	if err := parser.Unmarshal(resp.Body(), r); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !r.Ok {
		return nil, fmt.Errorf("API error: %s", r.Error)
	}

	if r.Result == nil {
		return nil, errors.New("API response result is nil")
	}

	return r.Result, nil
}

func StoreAccessToken(token string) (int, error) {
	err := ioutil.WriteFile(tokenPath, []byte(token), 0644)
	if err != nil {
		fmt.Println(err)
		return 1, errors.New("Couldn't save token :(")
	}
	return 0, nil
}

func FetchAccessToken() (string, error) {
	data, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
