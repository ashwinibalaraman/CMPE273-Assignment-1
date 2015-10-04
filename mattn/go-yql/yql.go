package yql

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	oauth "github.com/akrennmair/goauth"
)

const endpoint = "http://query.yahooapis.com/v1/public/yql"

var (
	yqlOauth *oauth.OAuthConsumer
)

func init() {
	sql.Register("yql", &YQLDriver{})
}

type YQLDriver struct{}

func (d *YQLDriver) Open(dsn string) (driver.Conn, error) {
	if len(dsn) > 1 {
		parts := strings.Split(dsn, "|")
		if len(parts) == 2 {
			return &YQLConn{http.DefaultClient, parts[0], parts[1]}, nil
		}

	}
	return &YQLConn{c: http.DefaultClient}, nil
}

type YQLConn struct {
	c      *http.Client
	key    string
	secret string
}

func (c *YQLConn) Close() error {
	c.c = nil
	return nil
}

func (c *YQLConn) Begin() (driver.Tx, error) {
	return nil, errors.New("Begin not supported")
}

func (c *YQLConn) Prepare(query string) (driver.Stmt, error) {
	return &YQLStmt{c, query}, nil
}

type YQLStmt struct {
	c *YQLConn
	q string
}

func (s *YQLStmt) Close() error {
	return nil
}

func (s *YQLStmt) NumInput() int {
	// TODO: strict check
	return strings.Count(s.q, "?")
}

func (s *YQLStmt) bind(args []driver.Value) error {
	b := s.q
	for _, v := range args {
		// TODO: strict check
		b = strings.Replace(b, "?", fmt.Sprintf("%q", v), 1)
	}
	s.q = b
	return nil
}
func (s *YQLStmt) QueryRow(args []driver.Value) (sql.Row, error) {
	fmt.Println("-------IN YQL.GO QUERYROW--------")
	var row sql.Row
	return row, nil
}
func (s *YQLStmt) Query(args []driver.Value) (driver.Rows, error) {
	if err := s.bind(args); err != nil {
		return nil, err
	}

	var res *http.Response
	var err error
	if len(s.c.key) > 1 {
		// secure
		yqlOauth := &oauth.OAuthConsumer{
			Service:         "yql",
			RequestTokenURL: "https://api.login.yahoo.com/oauth/v2/get_request_token",
			AccessTokenURL:  "https://api.login.yahoo.com/oauth/v2/get_token",
			CallBackURL:     "oob",
			ConsumerKey:     s.c.key,
			ConsumerSecret:  s.c.secret,
			Timeout:         5e9,
		}
		p := oauth.Params{}
		p.Add(&oauth.Pair{Key: "format", Value: "json"})
		p.Add(&oauth.Pair{Key: "q", Value: s.q})
		res, err = yqlOauth.Get(endpoint, p, nil)
		if err != nil {
			return nil, err
		}
	} else {
		url := fmt.Sprintf("%s?q=%s&env=store://datatables.org/alltableswithkeys&format=json", endpoint, url.QueryEscape(s.q))
		fmt.Println("url is: ", url)
		res, err = http.Get(url)
	}

	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	var data interface{}
	err = json.NewDecoder(res.Body).Decode(&data)
	if err != nil {
		fmt.Println("error is:", err)
		return nil, errors.New("Invalid Json: %v")
	}
	fmt.Println("data is", data)
	if data == nil {
		fmt.Println("data nil")
		return nil, errors.New("Unsupported result")
	}
	var ok bool
	data = data.(map[string]interface{})["query"]
	if data == nil {
		fmt.Println("query nil")
		return nil, errors.New("Unsupported result")
	}
	data = data.(map[string]interface{})["results"]
	if data == nil {
		fmt.Println("results nil")
		return nil, errors.New("Unsupported result")
	}
	results, ok := data.(map[string]interface{})
	if !ok {
		fmt.Println("data not ok")
		return nil, errors.New("Unsupported result")
	}
	fmt.Println("results-------->", results)
	for _, v := range results {
		if vv, ok := v.([]interface{}); ok {
			fmt.Println("vv ok------->", v)
			return &YQLRows{s, 0, vv}, nil
		}
	}
	fmt.Println("unknown error")
	return nil, errors.New("May be a single row response")
	//return nil, nil
}

type YQLResult struct {
	s *YQLStmt
}

func (s *YQLStmt) Exec(args []driver.Value) (driver.Result, error) {
	return nil, errors.New("Exec does not supported")
}

type YQLRows struct {
	s *YQLStmt
	n int
	d []interface{}
}

type YQLRow struct {
	s *YQLStmt
	n int
	d interface{}
}

func (rc *YQLRows) Close() error {
	return nil
}

func (rc *YQLRows) Columns() []string {
	return []string{"results"}
}

func (rc *YQLRows) Next(dest []driver.Value) error {
	if rc.n == len(rc.d) {
		return io.EOF
	}
	if s, ok := rc.d[rc.n].(string); ok {
		dest[0] = s
	} else {
		dest[0] = rc.d[rc.n]
	}
	rc.n++
	return nil
}
