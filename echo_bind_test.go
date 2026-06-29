package params

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type echoBindRequest struct {
	Name      String `json:"name"      form:"name"      query:"name"`
	Age       Int    `json:"age"       form:"age"       query:"age"`
	Score     Float  `json:"score"     form:"score"     query:"score"`
	Active    Bool   `json:"active"    form:"active"    query:"active"`
	Token     UUID   `json:"token"     form:"token"     query:"token"`
	CreatedAt Time   `json:"createdAt" form:"createdAt" query:"createdAt"`
}

var (
	testUUID      = uuid.MustParse("550e8400-e29b-41d4-a716-446655440000")
	testUUIDStr   = "550e8400-e29b-41d4-a716-446655440000"
	testTimeStr   = "2025-09-09T13:20:25Z"
	testTimeParsed, _ = time.Parse(time.RFC3339, testTimeStr)
)

func newEcho() *echo.Echo {
	e := echo.New()
	e.Binder = &echo.DefaultBinder{}
	return e
}

func assertFull(t *testing.T, req echoBindRequest) {
	t.Helper()
	require.True(t, req.Name.Present())
	require.Equal(t, "Alice", req.Name.Value())
	require.True(t, req.Age.Present())
	require.Equal(t, 30, req.Age.Value())
	require.True(t, req.Score.Present())
	require.InDelta(t, 9.5, req.Score.Value(), 0.0001)
	require.True(t, req.Active.Present())
	require.True(t, req.Active.Value())
	require.True(t, req.Token.Present())
	require.Equal(t, testUUID, req.Token.Value())
	require.True(t, req.CreatedAt.Present())
	require.Equal(t, testTimeParsed.UTC(), req.CreatedAt.Value().UTC())
}

// --- Query params ---

func TestEchoBind_QueryParams(t *testing.T) {
	e := newEcho()
	q := make(url.Values)
	q.Set("name", "Alice")
	q.Set("age", "30")
	q.Set("score", "9.5")
	q.Set("active", "true")
	q.Set("token", testUUIDStr)
	q.Set("createdAt", testTimeStr)

	r := httptest.NewRequest(http.MethodGet, "/?"+q.Encode(), nil)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)

	var req echoBindRequest
	require.NoError(t, c.Bind(&req))
	assertFull(t, req)
}

func TestEchoBind_QueryParams_MissingFields(t *testing.T) {
	e := newEcho()
	r := httptest.NewRequest(http.MethodGet, "/?name=Bob", nil)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)

	var req echoBindRequest
	require.NoError(t, c.Bind(&req))
	require.True(t, req.Name.Present())
	require.Equal(t, "Bob", req.Name.Value())
	require.False(t, req.Age.Present())
	require.False(t, req.Score.Present())
	require.False(t, req.Active.Present())
	require.False(t, req.Token.Present())
	require.False(t, req.CreatedAt.Present())
}

// --- Form params ---

func TestEchoBind_FormParams(t *testing.T) {
	e := newEcho()
	f := make(url.Values)
	f.Set("name", "Alice")
	f.Set("age", "30")
	f.Set("score", "9.5")
	f.Set("active", "true")
	f.Set("token", testUUIDStr)
	f.Set("createdAt", testTimeStr)

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)

	var req echoBindRequest
	require.NoError(t, c.Bind(&req))
	assertFull(t, req)
}

func TestEchoBind_FormParams_MissingFields(t *testing.T) {
	e := newEcho()
	f := make(url.Values)
	f.Set("age", "42")

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(f.Encode()))
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)

	var req echoBindRequest
	require.NoError(t, c.Bind(&req))
	require.False(t, req.Name.Present())
	require.True(t, req.Age.Present())
	require.Equal(t, 42, req.Age.Value())
}

// --- JSON body ---

func TestEchoBind_JSON(t *testing.T) {
	e := newEcho()
	body := `{"name":"Alice","age":30,"score":9.5,"active":true,"token":"` + testUUIDStr + `","createdAt":"` + testTimeStr + `"}`

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)

	var req echoBindRequest
	require.NoError(t, c.Bind(&req))
	assertFull(t, req)
}

func TestEchoBind_JSON_NullFields(t *testing.T) {
	e := newEcho()
	body := `{"name":null,"age":null,"score":null,"active":null,"token":null,"createdAt":null}`

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)

	var req echoBindRequest
	require.NoError(t, c.Bind(&req))
	require.False(t, req.Name.Present())
	require.False(t, req.Age.Present())
	require.False(t, req.Score.Present())
	require.False(t, req.Active.Present())
	require.False(t, req.Token.Present())
	require.False(t, req.CreatedAt.Present())
}

func TestEchoBind_JSON_MissingFields(t *testing.T) {
	e := newEcho()
	body := `{"name":"Charlie"}`

	r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)

	var req echoBindRequest
	require.NoError(t, c.Bind(&req))
	require.True(t, req.Name.Present())
	require.Equal(t, "Charlie", req.Name.Value())
	require.False(t, req.Age.Present())
}

// --- Mix: Form + Query ---

func TestEchoBind_FormAndQuery(t *testing.T) {
	e := newEcho()
	// name and age via query, score and active via form
	q := make(url.Values)
	q.Set("name", "Alice")
	q.Set("age", "30")

	f := make(url.Values)
	f.Set("score", "9.5")
	f.Set("active", "true")
	f.Set("token", testUUIDStr)
	f.Set("createdAt", testTimeStr)

	r := httptest.NewRequest(http.MethodPost, "/?"+q.Encode(), strings.NewReader(f.Encode()))
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)

	var req echoBindRequest
	require.NoError(t, c.Bind(&req))
	assertFull(t, req)
}

// --- Mix: JSON + Query ---
// Echo does NOT bind query params automatically for POST/PUT with JSON body (by design, see issue #1670).
// To combine JSON body with query params, call BindQueryParams and BindBody explicitly.

func TestEchoBind_JSONAndQuery(t *testing.T) {
	e := newEcho()
	q := make(url.Values)
	q.Set("name", "Alice")

	body := `{"age":30,"score":9.5,"active":true,"token":"` + testUUIDStr + `","createdAt":"` + testTimeStr + `"}`

	r := httptest.NewRequest(http.MethodPost, "/?"+q.Encode(), strings.NewReader(body))
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)

	binder := &echo.DefaultBinder{}
	var req echoBindRequest
	require.NoError(t, binder.BindQueryParams(c, &req))
	require.NoError(t, binder.BindBody(c, &req))
	assertFull(t, req)
}

func TestEchoBind_JSONAndQuery_AutoBindSkipsQuery(t *testing.T) {
	// Verify that c.Bind on POST with JSON does NOT include query params — expected Echo behavior.
	e := newEcho()
	q := make(url.Values)
	q.Set("name", "Alice")

	body := `{"age":30}`

	r := httptest.NewRequest(http.MethodPost, "/?"+q.Encode(), strings.NewReader(body))
	r.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	w := httptest.NewRecorder()
	c := e.NewContext(r, w)

	var req echoBindRequest
	require.NoError(t, c.Bind(&req))
	require.False(t, req.Name.Present(), "query params are not bound for POST+JSON via c.Bind")
	require.True(t, req.Age.Present())
	require.Equal(t, 30, req.Age.Value())
}
