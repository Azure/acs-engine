package httpd_test

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gogo/protobuf/proto"
	"github.com/golang/snappy"
	"github.com/influxdata/influxdb/internal"
	"github.com/influxdata/influxdb/models"
	"github.com/influxdata/influxdb/prometheus/remote"
	"github.com/influxdata/influxdb/query"
	"github.com/influxdata/influxdb/services/httpd"
	"github.com/influxdata/influxdb/services/meta"
	"github.com/influxdata/influxql"
)

// Ensure the handler returns results from a query (including nil results).
func TestHandler_Query(t *testing.T) {
	h := NewHandler(false)
	h.StatementExecutor.ExecuteStatementFn = func(stmt influxql.Statement, ctx *query.ExecutionContext) error {
		if stmt.String() != `SELECT * FROM bar` {
			t.Fatalf("unexpected query: %s", stmt.String())
		} else if ctx.Database != `foo` {
			t.Fatalf("unexpected db: %s", ctx.Database)
		}
		ctx.Results <- &query.Result{StatementID: 1, Series: models.Rows([]*models.Row{{Name: "series0"}})}
		ctx.Results <- &query.Result{StatementID: 2, Series: models.Rows([]*models.Row{{Name: "series1"}})}
		return nil
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewJSONRequest("GET", "/query?db=foo&q=SELECT+*+FROM+bar", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	} else if body := strings.TrimSpace(w.Body.String()); body != `{"results":[{"statement_id":1,"series":[{"name":"series0"}]},{"statement_id":2,"series":[{"name":"series1"}]}]}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

// Ensure the handler returns results from a query passed as a file.
func TestHandler_Query_File(t *testing.T) {
	h := NewHandler(false)
	h.StatementExecutor.ExecuteStatementFn = func(stmt influxql.Statement, ctx *query.ExecutionContext) error {
		if stmt.String() != `SELECT * FROM bar` {
			t.Fatalf("unexpected query: %s", stmt.String())
		} else if ctx.Database != `foo` {
			t.Fatalf("unexpected db: %s", ctx.Database)
		}
		ctx.Results <- &query.Result{StatementID: 1, Series: models.Rows([]*models.Row{{Name: "series0"}})}
		ctx.Results <- &query.Result{StatementID: 2, Series: models.Rows([]*models.Row{{Name: "series1"}})}
		return nil
	}

	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	part, err := writer.CreateFormFile("q", "")
	if err != nil {
		t.Fatal(err)
	}
	io.WriteString(part, "SELECT * FROM bar")

	if err := writer.Close(); err != nil {
		t.Fatal(err)
	}

	r := MustNewJSONRequest("POST", "/query?db=foo", &body)
	r.Header.Set("Content-Type", writer.FormDataContentType())

	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	} else if body := strings.TrimSpace(w.Body.String()); body != `{"results":[{"statement_id":1,"series":[{"name":"series0"}]},{"statement_id":2,"series":[{"name":"series1"}]}]}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

// Test query with user authentication.
func TestHandler_Query_Auth(t *testing.T) {
	// Create the handler to be tested.
	h := NewHandler(true)

	// Set mock meta client functions for the handler to use.
	h.MetaClient.AdminUserExistsFn = func() bool { return true }

	h.MetaClient.UserFn = func(username string) (meta.User, error) {
		if username != "user1" {
			return nil, meta.ErrUserNotFound
		}
		return &meta.UserInfo{
			Name:  "user1",
			Hash:  "abcd",
			Admin: true,
		}, nil
	}

	h.MetaClient.AuthenticateFn = func(u, p string) (meta.User, error) {
		if u != "user1" {
			return nil, fmt.Errorf("unexpected user: exp: user1, got: %s", u)
		} else if p != "abcd" {
			return nil, fmt.Errorf("unexpected password: exp: abcd, got: %s", p)
		}
		return h.MetaClient.User(u)
	}

	// Set mock query authorizer for handler to use.
	h.QueryAuthorizer.AuthorizeQueryFn = func(u meta.User, query *influxql.Query, database string) error {
		return nil
	}

	// Set mock statement executor for handler to use.
	h.StatementExecutor.ExecuteStatementFn = func(stmt influxql.Statement, ctx *query.ExecutionContext) error {
		if stmt.String() != `SELECT * FROM bar` {
			t.Fatalf("unexpected query: %s", stmt.String())
		} else if ctx.Database != `foo` {
			t.Fatalf("unexpected db: %s", ctx.Database)
		}
		ctx.Results <- &query.Result{StatementID: 1, Series: models.Rows([]*models.Row{{Name: "series0"}})}
		ctx.Results <- &query.Result{StatementID: 2, Series: models.Rows([]*models.Row{{Name: "series1"}})}
		return nil
	}

	// Test the handler with valid user and password in the URL parameters.
	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewJSONRequest("GET", "/query?u=user1&p=abcd&db=foo&q=SELECT+*+FROM+bar", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d: %s", w.Code, w.Body.String())
	} else if body := strings.TrimSpace(w.Body.String()); body != `{"results":[{"statement_id":1,"series":[{"name":"series0"}]},{"statement_id":2,"series":[{"name":"series1"}]}]}` {
		t.Fatalf("unexpected body: %s", body)
	}

	// Test the handler with valid user and password using basic auth.
	w = httptest.NewRecorder()
	r := MustNewJSONRequest("GET", "/query?db=foo&q=SELECT+*+FROM+bar", nil)
	r.SetBasicAuth("user1", "abcd")
	h.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d: %s", w.Code, w.Body.String())
	} else if body := strings.TrimSpace(w.Body.String()); body != `{"results":[{"statement_id":1,"series":[{"name":"series0"}]},{"statement_id":2,"series":[{"name":"series1"}]}]}` {
		t.Fatalf("unexpected body: %s", body)
	}

	// Test the handler with valid JWT bearer token.
	req := MustNewJSONRequest("GET", "/query?db=foo&q=SELECT+*+FROM+bar", nil)
	// Create a signed JWT token string and add it to the request header.
	_, signedToken := MustJWTToken("user1", h.Config.SharedSecret, false)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signedToken))

	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d: %s", w.Code, w.Body.String())
	} else if body := strings.TrimSpace(w.Body.String()); body != `{"results":[{"statement_id":1,"series":[{"name":"series0"}]},{"statement_id":2,"series":[{"name":"series1"}]}]}` {
		t.Fatalf("unexpected body: %s", body)
	}

	// Test the handler with JWT token signed with invalid key.
	req = MustNewJSONRequest("GET", "/query?db=foo&q=SELECT+*+FROM+bar", nil)
	// Create a signed JWT token string and add it to the request header.
	_, signedToken = MustJWTToken("user1", "invalid key", false)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signedToken))

	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d: %s", w.Code, w.Body.String())
	} else if body := strings.TrimSpace(w.Body.String()); body != `{"error":"signature is invalid"}` {
		t.Fatalf("unexpected body: %s", body)
	}

	// Test handler with valid JWT token carrying non-existant user.
	_, signedToken = MustJWTToken("bad_user", h.Config.SharedSecret, false)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signedToken))

	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d: %s", w.Code, w.Body.String())
	} else if body := strings.TrimSpace(w.Body.String()); body != `{"error":"user not found"}` {
		t.Fatalf("unexpected body: %s", body)
	}

	// Test handler with expired JWT token.
	_, signedToken = MustJWTToken("user1", h.Config.SharedSecret, true)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signedToken))

	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d: %s", w.Code, w.Body.String())
	} else if !strings.Contains(w.Body.String(), `{"error":"Token is expired`) {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}

	// Test handler with JWT token that has no expiration set.
	token, _ := MustJWTToken("user1", h.Config.SharedSecret, false)
	delete(token.Claims.(jwt.MapClaims), "exp")
	signedToken, err := token.SignedString([]byte(h.Config.SharedSecret))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", signedToken))
	w = httptest.NewRecorder()
	h.ServeHTTP(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("unexpected status: %d: %s", w.Code, w.Body.String())
	} else if body := strings.TrimSpace(w.Body.String()); body != `{"error":"token expiration required"}` {
		t.Fatalf("unexpected body: %s", body)
	}

	// Test the handler with valid user and password in the url and invalid in
	// basic auth (prioritize url).
	w = httptest.NewRecorder()
	r = MustNewJSONRequest("GET", "/query?u=user1&p=abcd&db=foo&q=SELECT+*+FROM+bar", nil)
	r.SetBasicAuth("user1", "efgh")
	h.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d: %s", w.Code, w.Body.String())
	} else if body := strings.TrimSpace(w.Body.String()); body != `{"results":[{"statement_id":1,"series":[{"name":"series0"}]},{"statement_id":2,"series":[{"name":"series1"}]}]}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

// Ensure the handler returns results from a query (including nil results).
func TestHandler_QueryRegex(t *testing.T) {
	h := NewHandler(false)
	h.StatementExecutor.ExecuteStatementFn = func(stmt influxql.Statement, ctx *query.ExecutionContext) error {
		if stmt.String() != `SELECT * FROM test WHERE url =~ /http\:\/\/www.akamai\.com/` {
			t.Fatalf("unexpected query: %s", stmt.String())
		} else if ctx.Database != `test` {
			t.Fatalf("unexpected db: %s", ctx.Database)
		}
		ctx.Results <- nil
		return nil
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewRequest("GET", "/query?db=test&q=SELECT%20%2A%20FROM%20test%20WHERE%20url%20%3D~%20%2Fhttp%5C%3A%5C%2F%5C%2Fwww.akamai%5C.com%2F", nil))
}

// Ensure the handler merges results from the same statement.
func TestHandler_Query_MergeResults(t *testing.T) {
	h := NewHandler(false)
	h.StatementExecutor.ExecuteStatementFn = func(stmt influxql.Statement, ctx *query.ExecutionContext) error {
		ctx.Results <- &query.Result{StatementID: 1, Series: models.Rows([]*models.Row{{Name: "series0"}})}
		ctx.Results <- &query.Result{StatementID: 1, Series: models.Rows([]*models.Row{{Name: "series1"}})}
		return nil
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewJSONRequest("GET", "/query?db=foo&q=SELECT+*+FROM+bar", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	} else if body := strings.TrimSpace(w.Body.String()); body != `{"results":[{"statement_id":1,"series":[{"name":"series0"},{"name":"series1"}]}]}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

// Ensure the handler merges results from the same statement.
func TestHandler_Query_MergeEmptyResults(t *testing.T) {
	h := NewHandler(false)
	h.StatementExecutor.ExecuteStatementFn = func(stmt influxql.Statement, ctx *query.ExecutionContext) error {
		ctx.Results <- &query.Result{StatementID: 1, Series: models.Rows{}}
		ctx.Results <- &query.Result{StatementID: 1, Series: models.Rows([]*models.Row{{Name: "series1"}})}
		return nil
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewJSONRequest("GET", "/query?db=foo&q=SELECT+*+FROM+bar", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	} else if body := strings.TrimSpace(w.Body.String()); body != `{"results":[{"statement_id":1,"series":[{"name":"series1"}]}]}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

// Ensure the handler can parse chunked and chunk size query parameters.
func TestHandler_Query_Chunked(t *testing.T) {
	h := NewHandler(false)
	h.StatementExecutor.ExecuteStatementFn = func(stmt influxql.Statement, ctx *query.ExecutionContext) error {
		if ctx.ChunkSize != 2 {
			t.Fatalf("unexpected chunk size: %d", ctx.ChunkSize)
		}
		ctx.Results <- &query.Result{StatementID: 1, Series: models.Rows([]*models.Row{{Name: "series0"}})}
		ctx.Results <- &query.Result{StatementID: 1, Series: models.Rows([]*models.Row{{Name: "series1"}})}
		return nil
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewJSONRequest("GET", "/query?db=foo&q=SELECT+*+FROM+bar&chunked=true&chunk_size=2", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	} else if w.Body.String() != `{"results":[{"statement_id":1,"series":[{"name":"series0"}]}]}
{"results":[{"statement_id":1,"series":[{"name":"series1"}]}]}
` {
		t.Fatalf("unexpected body: %s", w.Body.String())
	}
}

// Ensure the handler can accept an async query.
func TestHandler_Query_Async(t *testing.T) {
	done := make(chan struct{})
	h := NewHandler(false)
	h.StatementExecutor.ExecuteStatementFn = func(stmt influxql.Statement, ctx *query.ExecutionContext) error {
		if stmt.String() != `SELECT * FROM bar` {
			t.Fatalf("unexpected query: %s", stmt.String())
		} else if ctx.Database != `foo` {
			t.Fatalf("unexpected db: %s", ctx.Database)
		}
		ctx.Results <- &query.Result{StatementID: 1, Series: models.Rows([]*models.Row{{Name: "series0"}})}
		ctx.Results <- &query.Result{StatementID: 2, Series: models.Rows([]*models.Row{{Name: "series1"}})}
		close(done)
		return nil
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewJSONRequest("GET", "/query?db=foo&q=SELECT+*+FROM+bar&async=true", nil))
	if w.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", w.Code)
	} else if body := strings.TrimSpace(w.Body.String()); body != `` {
		t.Fatalf("unexpected body: %s", body)
	}

	// Wait to make sure the async query runs and completes.
	timer := time.NewTimer(100 * time.Millisecond)
	defer timer.Stop()

	select {
	case <-timer.C:
		t.Fatal("timeout while waiting for async query to complete")
	case <-done:
	}
}

// Ensure the handler returns a status 400 if the query is not passed in.
func TestHandler_Query_ErrQueryRequired(t *testing.T) {
	h := NewHandler(false)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewJSONRequest("GET", "/query", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", w.Code)
	} else if body := strings.TrimSpace(w.Body.String()); body != `{"error":"missing required parameter \"q\""}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

// Ensure the handler returns a status 400 if the query cannot be parsed.
func TestHandler_Query_ErrInvalidQuery(t *testing.T) {
	h := NewHandler(false)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewJSONRequest("GET", "/query?q=SELECT", nil))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", w.Code)
	} else if body := strings.TrimSpace(w.Body.String()); body != `{"error":"error parsing query: found EOF, expected identifier, string, number, bool at line 1, char 8"}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

// Ensure the handler returns an appropriate 401 or 403 status when authentication or authorization fails.
func TestHandler_Query_ErrAuthorize(t *testing.T) {
	h := NewHandler(true)
	h.QueryAuthorizer.AuthorizeQueryFn = func(u meta.User, q *influxql.Query, db string) error {
		return errors.New("marker")
	}
	h.MetaClient.AdminUserExistsFn = func() bool { return true }
	h.MetaClient.AuthenticateFn = func(u, p string) (meta.User, error) {

		users := []meta.UserInfo{
			{
				Name:  "admin",
				Hash:  "admin",
				Admin: true,
			},
			{
				Name: "user1",
				Hash: "abcd",
				Privileges: map[string]influxql.Privilege{
					"db0": influxql.ReadPrivilege,
				},
			},
		}

		for _, user := range users {
			if u == user.Name {
				if p == user.Hash {
					return &user, nil
				}
				return nil, meta.ErrAuthenticate
			}
		}
		return nil, meta.ErrUserNotFound
	}

	for i, tt := range []struct {
		user     string
		password string
		query    string
		code     int
	}{
		{
			query: "/query?q=SHOW+DATABASES",
			code:  http.StatusUnauthorized,
		},
		{
			user:     "user1",
			password: "abcd",
			query:    "/query?q=SHOW+DATABASES",
			code:     http.StatusForbidden,
		},
		{
			user:     "user2",
			password: "abcd",
			query:    "/query?q=SHOW+DATABASES",
			code:     http.StatusUnauthorized,
		},
	} {
		w := httptest.NewRecorder()
		r := MustNewJSONRequest("GET", tt.query, nil)
		params := r.URL.Query()
		if tt.user != "" {
			params.Set("u", tt.user)
		}
		if tt.password != "" {
			params.Set("p", tt.password)
		}
		r.URL.RawQuery = params.Encode()

		h.ServeHTTP(w, r)
		if w.Code != tt.code {
			t.Errorf("%d. unexpected status: got=%d exp=%d\noutput: %s", i, w.Code, tt.code, w.Body.String())
		}
	}
}

// Ensure the handler returns a status 200 if an error is returned in the result.
func TestHandler_Query_ErrResult(t *testing.T) {
	h := NewHandler(false)
	h.StatementExecutor.ExecuteStatementFn = func(stmt influxql.Statement, ctx *query.ExecutionContext) error {
		return errors.New("measurement not found")
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewJSONRequest("GET", "/query?db=foo&q=SHOW+SERIES+from+bin", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	} else if body := strings.TrimSpace(w.Body.String()); body != `{"results":[{"statement_id":0,"error":"measurement not found"}]}` {
		t.Fatalf("unexpected body: %s", body)
	}
}

// Ensure that closing the HTTP connection causes the query to be interrupted.
func TestHandler_Query_CloseNotify(t *testing.T) {
	// Avoid leaking a goroutine when this fails.
	done := make(chan struct{})
	defer close(done)

	interrupted := make(chan struct{})
	h := NewHandler(false)
	h.StatementExecutor.ExecuteStatementFn = func(stmt influxql.Statement, ctx *query.ExecutionContext) error {
		select {
		case <-ctx.Done():
		case <-done:
		}
		close(interrupted)
		return nil
	}

	s := httptest.NewServer(h)
	defer s.Close()

	// Parse the URL and generate a query request.
	u, err := url.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	u.Path = "/query"

	values := url.Values{}
	values.Set("q", "SELECT * FROM cpu")
	values.Set("db", "db0")
	values.Set("rp", "rp0")
	values.Set("chunked", "true")
	u.RawQuery = values.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		t.Fatal(err)
	}

	// Perform the request and retrieve the response.
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}

	// Validate that the interrupted channel has NOT been closed yet.
	timer := time.NewTimer(100 * time.Millisecond)
	select {
	case <-interrupted:
		timer.Stop()
		t.Fatal("query interrupted unexpectedly")
	case <-timer.C:
	}

	// Close the response body which should abort the query in the handler.
	resp.Body.Close()

	// The query should abort within 100 milliseconds.
	timer.Reset(100 * time.Millisecond)
	select {
	case <-interrupted:
		timer.Stop()
	case <-timer.C:
		t.Fatal("timeout while waiting for query to abort")
	}
}

// Ensure the prometheus remote write works
func TestHandler_PromWrite(t *testing.T) {
	req := &remote.WriteRequest{
		Timeseries: []*remote.TimeSeries{
			{
				Labels: []*remote.LabelPair{
					{Name: "host", Value: "a"},
					{Name: "region", Value: "west"},
				},
				Samples: []*remote.Sample{
					{TimestampMs: 1, Value: 1.2},
					{TimestampMs: 2, Value: math.NaN()},
				},
			},
		},
	}

	data, err := proto.Marshal(req)
	if err != nil {
		t.Fatal("couldn't marshal prometheus request")
	}
	compressed := snappy.Encode(nil, data)

	b := bytes.NewReader(compressed)
	h := NewHandler(false)
	h.MetaClient.DatabaseFn = func(name string) *meta.DatabaseInfo {
		return &meta.DatabaseInfo{}
	}
	called := false
	h.PointsWriter.WritePointsFn = func(db, rp string, _ models.ConsistencyLevel, _ meta.User, points []models.Point) error {
		called = true
		point := points[0]
		if point.UnixNano() != int64(time.Millisecond) {
			t.Fatalf("Exp point time %d but got %d", int64(time.Millisecond), point.UnixNano())
		}
		tags := point.Tags()
		expectedTags := models.Tags{models.Tag{Key: []byte("host"), Value: []byte("a")}, models.Tag{Key: []byte("region"), Value: []byte("west")}}
		if !reflect.DeepEqual(tags, expectedTags) {
			t.Fatalf("tags don't match\n\texp: %v\n\tgot: %v", expectedTags, tags)
		}

		fields, err := point.Fields()
		if err != nil {
			t.Fatal(err.Error())
		}
		expFields := models.Fields{"f64": 1.2}
		if !reflect.DeepEqual(fields, expFields) {
			t.Fatalf("fields don't match\n\texp: %v\n\tgot: %v", expFields, fields)
		}
		return nil
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewRequest("POST", "/api/v1/prom/write?db=foo", b))
	if !called {
		t.Fatal("WritePoints: expected call")
	}
	if w.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", w.Code)
	}
}

// Ensure Prometheus remote read requests are converted to the correct InfluxQL query and
// data is returned
func TestHandler_PromRead(t *testing.T) {
	req := &remote.ReadRequest{
		Queries: []*remote.Query{{
			Matchers: []*remote.LabelMatcher{
				{Type: remote.MatchType_EQUAL, Name: "eq", Value: "a"},
				{Type: remote.MatchType_NOT_EQUAL, Name: "neq", Value: "b"},
				{Type: remote.MatchType_REGEX_MATCH, Name: "regex", Value: "c"},
				{Type: remote.MatchType_REGEX_NO_MATCH, Name: "neqregex", Value: "d"},
			},
			StartTimestampMs: 1,
			EndTimestampMs:   2,
		}},
	}
	data, err := proto.Marshal(req)
	if err != nil {
		t.Fatal("couldn't marshal prometheus request")
	}
	compressed := snappy.Encode(nil, data)
	b := bytes.NewReader(compressed)

	h := NewHandler(false)
	h.StatementExecutor.ExecuteStatementFn = func(stmt influxql.Statement, ctx *query.ExecutionContext) error {
		if stmt.String() != `SELECT f64 FROM foo.._ WHERE eq = 'a' AND neq != 'b' AND regex =~ /c/ AND neqregex !~ /d/ AND time >= '1970-01-01T00:00:00.001Z' AND time <= '1970-01-01T00:00:00.002Z' GROUP BY *` {
			t.Fatalf("unexpected query: %s", stmt.String())
		} else if ctx.Database != `foo` {
			t.Fatalf("unexpected db: %s", ctx.Database)
		}
		row := &models.Row{
			Name:    "_",
			Tags:    map[string]string{"foo": "bar"},
			Columns: []string{"time", "f64"},
			Values:  [][]interface{}{{time.Unix(23, 0), 1.2}},
		}
		ctx.Results <- &query.Result{StatementID: 1, Series: models.Rows([]*models.Row{row})}
		return nil
	}

	w := httptest.NewRecorder()

	h.ServeHTTP(w, MustNewJSONRequest("POST", "/api/v1/prom/read?db=foo", b))
	if w.Code != http.StatusOK {
		t.Fatalf("unexpected status: %d", w.Code)
	}

	reqBuf, err := snappy.Decode(nil, w.Body.Bytes())
	if err != nil {
		t.Fatal(err.Error())
	}

	var resp remote.ReadResponse
	if err := proto.Unmarshal(reqBuf, &resp); err != nil {
		t.Fatal(err.Error())
	}

	expLabels := []*remote.LabelPair{{Name: "foo", Value: "bar"}}
	expSamples := []*remote.Sample{{TimestampMs: 23000, Value: 1.2}}

	ts := resp.Results[0].Timeseries[0]

	if !reflect.DeepEqual(expLabels, ts.Labels) {
		t.Fatalf("unexpected labels\n\texp: %v\n\tgot: %v", expLabels, ts.Labels)
	}
	if !reflect.DeepEqual(expSamples, ts.Samples) {
		t.Fatalf("unexpectd samples\n\texp: %v\n\tgot: %v", expSamples, ts.Samples)
	}
}

// Ensure the handler handles ping requests correctly.
// TODO: This should be expanded to verify the MetaClient check in servePing is working correctly
func TestHandler_Ping(t *testing.T) {
	h := NewHandler(false)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewRequest("GET", "/ping", nil))
	if w.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", w.Code)
	}
	h.ServeHTTP(w, MustNewRequest("HEAD", "/ping", nil))
	if w.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", w.Code)
	}
}

// Ensure the handler returns the version correctly from the different endpoints.
func TestHandler_Version(t *testing.T) {
	h := NewHandler(false)
	h.StatementExecutor.ExecuteStatementFn = func(stmt influxql.Statement, ctx *query.ExecutionContext) error {
		return nil
	}
	tests := []struct {
		method   string
		endpoint string
		body     io.Reader
	}{
		{
			method:   "GET",
			endpoint: "/ping",
			body:     nil,
		},
		{
			method:   "GET",
			endpoint: "/query?db=foo&q=SELECT+*+FROM+bar",
			body:     nil,
		},
		{
			method:   "POST",
			endpoint: "/write",
			body:     bytes.NewReader(make([]byte, 10)),
		},
		{
			method:   "GET",
			endpoint: "/notfound",
			body:     nil,
		},
	}

	for _, test := range tests {
		w := httptest.NewRecorder()
		h.ServeHTTP(w, MustNewRequest(test.method, test.endpoint, test.body))
		if v := w.HeaderMap["X-Influxdb-Version"]; len(v) > 0 {
			if v[0] != "0.0.0" {
				t.Fatalf("unexpected version: %s", v)
			}
		} else {
			t.Fatalf("Header entry 'X-Influxdb-Version' not present")
		}

		if v := w.HeaderMap["X-Influxdb-Build"]; len(v) > 0 {
			if v[0] != "OSS" {
				t.Fatalf("unexpected BuildType: %s", v)
			}
		} else {
			t.Fatalf("Header entry 'X-Influxdb-Build' not present")
		}
	}
}

// Ensure the handler handles status requests correctly.
func TestHandler_Status(t *testing.T) {
	h := NewHandler(false)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewRequest("GET", "/status", nil))
	if w.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", w.Code)
	}
	h.ServeHTTP(w, MustNewRequest("HEAD", "/status", nil))
	if w.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", w.Code)
	}
}

// Ensure write endpoint can handle bad requests
func TestHandler_HandleBadRequestBody(t *testing.T) {
	b := bytes.NewReader(make([]byte, 10))
	h := NewHandler(false)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewRequest("POST", "/write", b))
	if w.Code != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", w.Code)
	}
}

func TestHandler_Write_EntityTooLarge_ContentLength(t *testing.T) {
	b := bytes.NewReader(make([]byte, 100))
	h := NewHandler(false)
	h.Config.MaxBodySize = 5
	h.MetaClient.DatabaseFn = func(name string) *meta.DatabaseInfo {
		return &meta.DatabaseInfo{}
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewRequest("POST", "/write?db=foo", b))
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("unexpected status: %d", w.Code)
	}
}

func TestHandler_Write_SuppressLog(t *testing.T) {
	var buf bytes.Buffer
	c := httpd.NewConfig()
	c.SuppressWriteLog = true
	h := NewHandlerWithConfig(c)
	h.CLFLogger = log.New(&buf, "", log.LstdFlags)
	h.MetaClient.DatabaseFn = func(name string) *meta.DatabaseInfo {
		return &meta.DatabaseInfo{}
	}
	h.PointsWriter.WritePointsFn = func(database, retentionPolicy string, consistencyLevel models.ConsistencyLevel, user meta.User, points []models.Point) error {
		return nil
	}

	b := strings.NewReader("cpu,host=server01 value=2\n")
	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewRequest("POST", "/write?db=foo", b))
	if w.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", w.Code)
	}

	// If the log has anything in it, this failed.
	if buf.Len() > 0 {
		t.Fatalf("expected no bytes to be written to the log, got %d", buf.Len())
	}
}

// onlyReader implements io.Reader only to ensure Request.ContentLength is not set
type onlyReader struct {
	r io.Reader
}

func (o onlyReader) Read(p []byte) (n int, err error) {
	return o.r.Read(p)
}

func TestHandler_Write_EntityTooLarge_NoContentLength(t *testing.T) {
	b := onlyReader{bytes.NewReader(make([]byte, 100))}
	h := NewHandler(false)
	h.Config.MaxBodySize = 5
	h.MetaClient.DatabaseFn = func(name string) *meta.DatabaseInfo {
		return &meta.DatabaseInfo{}
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewRequest("POST", "/write?db=foo", b))
	if w.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("unexpected status: %d", w.Code)
	}
}

// TestHandler_Write_NegativeMaxBodySize verifies no error occurs if MaxBodySize is < 0
func TestHandler_Write_NegativeMaxBodySize(t *testing.T) {
	b := bytes.NewReader([]byte(`foo n=1`))
	h := NewHandler(false)
	h.Config.MaxBodySize = -1
	h.MetaClient.DatabaseFn = func(name string) *meta.DatabaseInfo {
		return &meta.DatabaseInfo{}
	}
	called := false
	h.PointsWriter.WritePointsFn = func(_, _ string, _ models.ConsistencyLevel, _ meta.User, _ []models.Point) error {
		called = true
		return nil
	}

	w := httptest.NewRecorder()
	h.ServeHTTP(w, MustNewRequest("POST", "/write?db=foo", b))
	if !called {
		t.Fatal("WritePoints: expected call")
	}
	if w.Code != http.StatusNoContent {
		t.Fatalf("unexpected status: %d", w.Code)
	}
}

// Ensure X-Forwarded-For header writes the correct log message.
func TestHandler_XForwardedFor(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(false)
	h.CLFLogger = log.New(&buf, "", 0)

	req := MustNewRequest("GET", "/query", nil)
	req.Header.Set("X-Forwarded-For", "192.168.0.1")
	req.RemoteAddr = "127.0.0.1"
	h.ServeHTTP(httptest.NewRecorder(), req)

	parts := strings.Split(buf.String(), " ")
	if parts[0] != "192.168.0.1,127.0.0.1" {
		t.Errorf("unexpected host ip address: %s", parts[0])
	}
}

func TestHandler_XRequestId(t *testing.T) {
	var buf bytes.Buffer
	h := NewHandler(false)
	h.CLFLogger = log.New(&buf, "", 0)

	cases := []map[string]string{
		{"X-Request-Id": "abc123", "Request-Id": ""},          // X-Request-Id is used.
		{"X-REQUEST-ID": "cde", "Request-Id": ""},             // X-REQUEST-ID is used.
		{"X-Request-Id": "", "Request-Id": "foobarzoo"},       // Request-Id is used.
		{"X-Request-Id": "abc123", "Request-Id": "foobarzoo"}, // X-Request-Id takes precedence.
		{"X-Request-Id": "", "Request-Id": ""},                // v1 UUID generated.
	}

	for _, c := range cases {
		t.Run(fmt.Sprint(c), func(t *testing.T) {
			buf.Reset()
			req := MustNewRequest("GET", "/ping", nil)
			req.RemoteAddr = "127.0.0.1"

			// Set the relevant request ID headers
			var allEmpty = true
			for k, v := range c {
				req.Header.Set(k, v)
				if v != "" {
					allEmpty = false
				}
			}

			w := httptest.NewRecorder()
			h.ServeHTTP(w, req)

			// Split up the HTTP log line. The request ID is currently located in
			// index 12. If the log line gets changed in the future, this test
			// will likely break and the index will need to be updated.
			parts := strings.Split(buf.String(), " ")
			i := 12

			// If neither header is set then we expect a v1 UUID to be generated.
			if allEmpty {
				if got, exp := len(parts[i]), 36; got != exp {
					t.Fatalf("got ID of length %d, expected one of length %d", got, exp)
				}
			} else if c["X-Request-Id"] != "" {
				if got, exp := parts[i], c["X-Request-Id"]; got != exp {
					t.Fatalf("got ID of %q, expected %q", got, exp)
				}
			} else if c["X-REQUEST-ID"] != "" {
				if got, exp := parts[i], c["X-REQUEST-ID"]; got != exp {
					t.Fatalf("got ID of %q, expected %q", got, exp)
				}
			} else {
				if got, exp := parts[i], c["Request-Id"]; got != exp {
					t.Fatalf("got ID of %q, expected %q", got, exp)
				}
			}

			// Check response headers
			if got, exp := w.Header().Get("Request-Id"), parts[i]; got != exp {
				t.Fatalf("Request-Id header was %s, expected %s", got, exp)
			} else if got, exp := w.Header().Get("X-Request-Id"), parts[i]; got != exp {
				t.Fatalf("X-Request-Id header was %s, expected %s", got, exp)
			}
		})
	}
}

func TestThrottler_Handler(t *testing.T) {
	t.Run("OK", func(t *testing.T) {
		throttler := httpd.NewThrottler(2, 98)

		// Send the total number of concurrent requests to the channel.
		var concurrentN int32
		concurrentCh := make(chan int)

		h := throttler.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			atomic.AddInt32(&concurrentN, 1)
			concurrentCh <- int(atomic.LoadInt32(&concurrentN))
			time.Sleep(1 * time.Millisecond)
			atomic.AddInt32(&concurrentN, -1)
		}))

		// Execute requests concurrently.
		const n = 100
		for i := 0; i < n; i++ {
			go func() { h.ServeHTTP(nil, nil) }()
		}

		// Read the number of concurrent requests for every execution.
		for i := 0; i < n; i++ {
			if v := <-concurrentCh; v > 2 {
				t.Fatalf("concurrent requests exceed maximum: %d", v)
			}
		}
	})

	t.Run("ErrTimeout", func(t *testing.T) {
		throttler := httpd.NewThrottler(2, 1)
		throttler.EnqueueTimeout = 1 * time.Millisecond

		resp := make(chan struct{})
		h := throttler.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp <- struct{}{}
		}))

		pending := make(chan struct{}, 2)

		// First two requests should execute immediately.
		go func() { pending <- struct{}{}; h.ServeHTTP(nil, nil) }()
		go func() { pending <- struct{}{}; h.ServeHTTP(nil, nil) }()

		<-pending
		<-pending

		// Third request should be enqueued but timeout.
		w := httptest.NewRecorder()
		h.ServeHTTP(w, nil)
		if w.Code != http.StatusServiceUnavailable {
			t.Fatalf("unexpected status code: %d", w.Code)
		} else if body := w.Body.String(); body != "request throttled, exceeds timeout\n" {
			t.Fatalf("unexpected response body: %q", body)
		}

		// Allow 2 existing requests to complete.
		<-resp
		<-resp
	})

	t.Run("ErrFull", func(t *testing.T) {
		delay := 100 * time.Millisecond
		if os.Getenv("CI") != "" {
			delay = 2 * time.Second
		}

		throttler := httpd.NewThrottler(2, 1)

		resp := make(chan struct{})
		h := throttler.Handler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			resp <- struct{}{}
		}))

		// First two requests should execute immediately and third should be queued.
		go func() { h.ServeHTTP(nil, nil) }()
		go func() { h.ServeHTTP(nil, nil) }()
		go func() { h.ServeHTTP(nil, nil) }()
		time.Sleep(delay)

		// Fourth request should fail when trying to enqueue.
		w := httptest.NewRecorder()
		h.ServeHTTP(w, nil)
		if w.Code != http.StatusServiceUnavailable {
			t.Fatalf("unexpected status code: %d", w.Code)
		} else if body := w.Body.String(); body != "request throttled, queue full\n" {
			t.Fatalf("unexpected response body: %q", body)
		}

		// Allow 3 existing requests to complete.
		<-resp
		<-resp
		<-resp
	})
}

// NewHandler represents a test wrapper for httpd.Handler.
type Handler struct {
	*httpd.Handler
	MetaClient        *internal.MetaClientMock
	StatementExecutor HandlerStatementExecutor
	QueryAuthorizer   HandlerQueryAuthorizer
	PointsWriter      HandlerPointsWriter
}

// NewHandler returns a new instance of Handler.
func NewHandler(requireAuthentication bool) *Handler {
	config := httpd.NewConfig()
	config.AuthEnabled = requireAuthentication
	config.SharedSecret = "super secret key"
	return NewHandlerWithConfig(config)
}

func NewHandlerWithConfig(config httpd.Config) *Handler {
	h := &Handler{
		Handler: httpd.NewHandler(config),
	}

	h.MetaClient = &internal.MetaClientMock{}

	h.Handler.MetaClient = h.MetaClient
	h.Handler.QueryExecutor = query.NewExecutor()
	h.Handler.QueryExecutor.StatementExecutor = &h.StatementExecutor
	h.Handler.QueryAuthorizer = &h.QueryAuthorizer
	h.Handler.PointsWriter = &h.PointsWriter
	h.Handler.Version = "0.0.0"
	h.Handler.BuildType = "OSS"
	return h
}

// HandlerStatementExecutor is a mock implementation of Handler.StatementExecutor.
type HandlerStatementExecutor struct {
	ExecuteStatementFn func(stmt influxql.Statement, ctx *query.ExecutionContext) error
}

func (e *HandlerStatementExecutor) ExecuteStatement(stmt influxql.Statement, ctx *query.ExecutionContext) error {
	return e.ExecuteStatementFn(stmt, ctx)
}

// HandlerQueryAuthorizer is a mock implementation of Handler.QueryAuthorizer.
type HandlerQueryAuthorizer struct {
	AuthorizeQueryFn func(u meta.User, query *influxql.Query, database string) error
}

func (a *HandlerQueryAuthorizer) AuthorizeQuery(u meta.User, query *influxql.Query, database string) error {
	return a.AuthorizeQueryFn(u, query, database)
}

type HandlerPointsWriter struct {
	WritePointsFn func(database, retentionPolicy string, consistencyLevel models.ConsistencyLevel, user meta.User, points []models.Point) error
}

func (h *HandlerPointsWriter) WritePoints(database, retentionPolicy string, consistencyLevel models.ConsistencyLevel, user meta.User, points []models.Point) error {
	return h.WritePointsFn(database, retentionPolicy, consistencyLevel, user, points)
}

// MustNewRequest returns a new HTTP request. Panic on error.
func MustNewRequest(method, urlStr string, body io.Reader) *http.Request {
	r, err := http.NewRequest(method, urlStr, body)
	if err != nil {
		panic(err.Error())
	}
	return r
}

// MustNewRequest returns a new HTTP request with the content type set. Panic on error.
func MustNewJSONRequest(method, urlStr string, body io.Reader) *http.Request {
	r := MustNewRequest(method, urlStr, body)
	r.Header.Set("Accept", "application/json")
	return r
}

// MustJWTToken returns a new JWT token and signed string or panics trying.
func MustJWTToken(username, secret string, expired bool) (*jwt.Token, string) {
	token := jwt.New(jwt.GetSigningMethod("HS512"))
	token.Claims.(jwt.MapClaims)["username"] = username
	if expired {
		token.Claims.(jwt.MapClaims)["exp"] = time.Now().Add(-time.Second).Unix()
	} else {
		token.Claims.(jwt.MapClaims)["exp"] = time.Now().Add(time.Minute * 10).Unix()
	}
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		panic(err)
	}
	return token, signed
}
