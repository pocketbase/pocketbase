package router_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/pocketbase/pocketbase/tools/hook"
	"github.com/pocketbase/pocketbase/tools/router"
)

func TestRouter(t *testing.T) {
	calls := ""

	r := router.NewRouter(func(w http.ResponseWriter, r *http.Request) (*router.Event, router.EventCleanupFunc) {
		return &router.Event{
				Response: w,
				Request:  r,
			},
			func() {
				calls += ":cleanup"
			}
	})

	r.BindFunc(func(e *router.Event) error {
		calls += "root_m:"

		err := e.Next()

		if err != nil {
			calls += "/error"
		}

		return err
	})

	r.Any("/any", func(e *router.Event) error {
		calls += "/any"
		return nil
	})

	r.GET("/a", func(e *router.Event) error {
		calls += "/a"
		return nil
	})

	g1 := r.Group("/a/b").BindFunc(func(e *router.Event) error {
		calls += "a_b_group_m:"
		return e.Next()
	})
	g1.GET("/1", func(e *router.Event) error {
		calls += "/1_get"
		return nil
	}).BindFunc(func(e *router.Event) error {
		calls += "1_get_m:"
		return e.Next()
	})
	g1.POST("/1", func(e *router.Event) error {
		calls += "/1_post"
		return nil
	})
	g1.GET("/{param}", func(e *router.Event) error {
		calls += "/" + e.Request.PathValue("param")
		return errors.New("test") // should be normalized to an ApiError
	})

	mux, err := r.BuildMux()
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := ts.Client()

	scenarios := []struct {
		method string
		path   string
		calls  string
	}{
		{http.MethodGet, "/any", "root_m:/any:cleanup"},
		{http.MethodOptions, "/any", "root_m:/any:cleanup"},
		{http.MethodPatch, "/any", "root_m:/any:cleanup"},
		{http.MethodPut, "/any", "root_m:/any:cleanup"},
		{http.MethodPost, "/any", "root_m:/any:cleanup"},
		{http.MethodDelete, "/any", "root_m:/any:cleanup"},
		// ---
		{http.MethodPost, "/a", "root_m:/error:cleanup"}, // missing
		{http.MethodGet, "/a", "root_m:/a:cleanup"},
		{http.MethodHead, "/a", "root_m:/a:cleanup"}, // auto registered with the GET
		{http.MethodGet, "/a/b/1", "root_m:a_b_group_m:1_get_m:/1_get:cleanup"},
		{http.MethodHead, "/a/b/1", "root_m:a_b_group_m:1_get_m:/1_get:cleanup"},
		{http.MethodPost, "/a/b/1", "root_m:a_b_group_m:/1_post:cleanup"},
		{http.MethodGet, "/a/b/456", "root_m:a_b_group_m:/456/error:cleanup"},
	}

	for _, s := range scenarios {
		t.Run(s.method+"_"+s.path, func(t *testing.T) {
			calls = "" // reset

			req, err := http.NewRequest(s.method, ts.URL+s.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if calls != s.calls {
				t.Fatalf("Expected calls\n%q\ngot\n%q", s.calls, calls)
			}
		})
	}
}

func TestRouterUnbind(t *testing.T) {
	calls := ""

	r := router.NewRouter(func(w http.ResponseWriter, r *http.Request) (*router.Event, router.EventCleanupFunc) {
		return &router.Event{
				Response: w,
				Request:  r,
			},
			func() {
				calls += ":cleanup"
			}
	})
	r.Bind(&hook.Handler[*router.Event]{
		Id: "root_1",
		Func: func(e *router.Event) error {
			calls += "root_1:"
			return e.Next()
		},
	})
	r.Bind(&hook.Handler[*router.Event]{
		Id: "root_2",
		Func: func(e *router.Event) error {
			calls += "root_2:"
			return e.Next()
		},
	})
	r.Bind(&hook.Handler[*router.Event]{
		Id: "root_3",
		Func: func(e *router.Event) error {
			calls += "root_3:"
			return e.Next()
		},
	})
	r.GET("/action", func(e *router.Event) error {
		calls += "root_action"
		return nil
	}).Unbind("root_1")

	ga := r.Group("/group_a")
	ga.Unbind("root_1")
	ga.Bind(&hook.Handler[*router.Event]{
		Id: "group_a_1",
		Func: func(e *router.Event) error {
			calls += "group_a_1:"
			return e.Next()
		},
	})
	ga.Bind(&hook.Handler[*router.Event]{
		Id: "group_a_2",
		Func: func(e *router.Event) error {
			calls += "group_a_2:"
			return e.Next()
		},
	})
	ga.Bind(&hook.Handler[*router.Event]{
		Id: "group_a_3",
		Func: func(e *router.Event) error {
			calls += "group_a_3:"
			return e.Next()
		},
	})
	ga.GET("/action", func(e *router.Event) error {
		calls += "group_a_action"
		return nil
	}).Unbind("root_2", "group_b_1", "group_a_1")

	gb := r.Group("/group_b")
	gb.Unbind("root_2")
	gb.Bind(&hook.Handler[*router.Event]{
		Id: "group_b_1",
		Func: func(e *router.Event) error {
			calls += "group_b_1:"
			return e.Next()
		},
	})
	gb.Bind(&hook.Handler[*router.Event]{
		Id: "group_b_2",
		Func: func(e *router.Event) error {
			calls += "group_b_2:"
			return e.Next()
		},
	})
	gb.Bind(&hook.Handler[*router.Event]{
		Id: "group_b_3",
		Func: func(e *router.Event) error {
			calls += "group_b_3:"
			return e.Next()
		},
	})
	gb.GET("/action", func(e *router.Event) error {
		calls += "group_b_action"
		return nil
	}).Unbind("group_b_3", "group_a_3", "root_3")

	mux, err := r.BuildMux()
	if err != nil {
		t.Fatal(err)
	}

	ts := httptest.NewServer(mux)
	defer ts.Close()

	client := ts.Client()

	scenarios := []struct {
		method string
		path   string
		calls  string
	}{
		{http.MethodGet, "/action", "root_2:root_3:root_action:cleanup"},
		{http.MethodGet, "/group_a/action", "root_3:group_a_2:group_a_3:group_a_action:cleanup"},
		{http.MethodGet, "/group_b/action", "root_1:group_b_1:group_b_2:group_b_action:cleanup"},
	}

	for _, s := range scenarios {
		t.Run(s.method+"_"+s.path, func(t *testing.T) {
			calls = "" // reset

			req, err := http.NewRequest(s.method, ts.URL+s.path, nil)
			if err != nil {
				t.Fatal(err)
			}

			_, err = client.Do(req)
			if err != nil {
				t.Fatal(err)
			}

			if calls != s.calls {
				t.Fatalf("Expected calls\n%q\ngot\n%q", s.calls, calls)
			}
		})
	}
}
