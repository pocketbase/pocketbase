package router

import (
	"errors"
	"fmt"
	"net/http"
	"slices"
	"testing"

	"github.com/pocketbase/pocketbase/tools/hook"
)

func TestRouterGroupGroup(t *testing.T) {
	t.Parallel()

	g0 := RouterGroup[*Event]{}

	g1 := g0.Group("test1")
	g2 := g0.Group("test2")

	if total := len(g0.children); total != 2 {
		t.Fatalf("Expected %d child groups, got %d", 2, total)
	}

	if g1.Prefix != "test1" {
		t.Fatalf("Expected g1 with prefix %q, got %q", "test1", g1.Prefix)
	}
	if g2.Prefix != "test2" {
		t.Fatalf("Expected g2 with prefix %q, got %q", "test2", g2.Prefix)
	}
}

func TestRouterGroupBindFunc(t *testing.T) {
	t.Parallel()

	g := RouterGroup[*Event]{}

	calls := ""

	// append one function
	g.BindFunc(func(e *Event) error {
		calls += "a"
		return nil
	})

	// append multiple functions
	g.BindFunc(
		func(e *Event) error {
			calls += "b"
			return nil
		},
		func(e *Event) error {
			calls += "c"
			return nil
		},
	)

	if total := len(g.Middlewares); total != 3 {
		t.Fatalf("Expected %d middlewares, got %v", 3, total)
	}

	for _, h := range g.Middlewares {
		_ = h.Func(nil)
	}

	if calls != "abc" {
		t.Fatalf("Expected calls sequence %q, got %q", "abc", calls)
	}
}

func TestRouterGroupBind(t *testing.T) {
	t.Parallel()

	g := RouterGroup[*Event]{
		// mock excluded middlewares to check whether the entry will be deleted
		excludedMiddlewares: map[string]struct{}{"test2": {}},
	}

	calls := ""

	// append one handler
	g.Bind(&hook.Handler[*Event]{
		Func: func(e *Event) error {
			calls += "a"
			return nil
		},
	})

	// append multiple handlers
	g.Bind(
		&hook.Handler[*Event]{
			Id: "test1",
			Func: func(e *Event) error {
				calls += "b"
				return nil
			},
		},
		&hook.Handler[*Event]{
			Id: "test2",
			Func: func(e *Event) error {
				calls += "c"
				return nil
			},
		},
	)

	if total := len(g.Middlewares); total != 3 {
		t.Fatalf("Expected %d middlewares, got %v", 3, total)
	}

	for _, h := range g.Middlewares {
		_ = h.Func(nil)
	}

	if calls != "abc" {
		t.Fatalf("Expected calls %q, got %q", "abc", calls)
	}

	// ensures that the previously excluded middleware was removed
	if len(g.excludedMiddlewares) != 0 {
		t.Fatalf("Expected test2 to be removed from the excludedMiddlewares list, got %v", g.excludedMiddlewares)
	}
}

func TestRouterGroupUnbind(t *testing.T) {
	t.Parallel()

	g := RouterGroup[*Event]{}

	calls := ""

	// anonymous middlewares
	g.Bind(&hook.Handler[*Event]{
		Func: func(e *Event) error {
			calls += "a"
			return nil // unused value
		},
	})

	// middlewares with id
	g.Bind(&hook.Handler[*Event]{
		Id: "test1",
		Func: func(e *Event) error {
			calls += "b"
			return nil // unused value
		},
	})
	g.Bind(&hook.Handler[*Event]{
		Id: "test2",
		Func: func(e *Event) error {
			calls += "c"
			return nil // unused value
		},
	})
	g.Bind(&hook.Handler[*Event]{
		Id: "test3",
		Func: func(e *Event) error {
			calls += "d"
			return nil // unused value
		},
	})

	// remove
	g.Unbind("") // should be no-op
	g.Unbind("test1", "test3")

	if total := len(g.Middlewares); total != 2 {
		t.Fatalf("Expected %d middlewares, got %v", 2, total)
	}

	for _, h := range g.Middlewares {
		if err := h.Func(nil); err != nil {
			continue
		}
	}

	if calls != "ac" {
		t.Fatalf("Expected calls %q, got %q", "ac", calls)
	}

	// ensure that the ids were added in the exclude list
	excluded := []string{"test1", "test3"}
	if len(g.excludedMiddlewares) != len(excluded) {
		t.Fatalf("Expected excludes %v, got %v", excluded, g.excludedMiddlewares)
	}
	for id := range g.excludedMiddlewares {
		if !slices.Contains(excluded, id) {
			t.Fatalf("Expected %q to be marked as excluded", id)
		}
	}
}

func TestRouterGroupRoute(t *testing.T) {
	t.Parallel()

	group := RouterGroup[*Event]{}

	sub := group.Group("sub")

	var called bool
	route := group.Route(http.MethodPost, "/test", func(e *Event) error {
		called = true
		return nil
	})

	// ensure that the route was registered only to the main one
	// ---
	if len(sub.children) != 0 {
		t.Fatalf("Expected no sub children, got %d", len(sub.children))
	}

	if len(group.children) != 2 {
		t.Fatalf("Expected %d group children, got %d", 2, len(group.children))
	}
	// ---

	// check the registered route
	// ---
	if route != group.children[1] {
		t.Fatalf("Expected group children %v, got %v", route, group.children[1])
	}

	if route.Method != http.MethodPost {
		t.Fatalf("Expected route method %q, got %q", http.MethodPost, route.Method)
	}

	if route.Path != "/test" {
		t.Fatalf("Expected route path %q, got %q", "/test", route.Path)
	}

	route.Action(nil)
	if !called {
		t.Fatal("Expected route action to be called")
	}
}

func TestRouterGroupRouteAliases(t *testing.T) {
	t.Parallel()

	group := RouterGroup[*Event]{}

	testErr := errors.New("test")

	testAction := func(e *Event) error {
		return testErr
	}

	scenarios := []struct {
		route        *Route[*Event]
		expectMethod string
		expectPath   string
	}{
		{
			group.Any("/test", testAction),
			"",
			"/test",
		},
		{
			group.GET("/test", testAction),
			http.MethodGet,
			"/test",
		},
		{
			group.SEARCH("/test", testAction),
			"SEARCH",
			"/test",
		},
		{
			group.POST("/test", testAction),
			http.MethodPost,
			"/test",
		},
		{
			group.DELETE("/test", testAction),
			http.MethodDelete,
			"/test",
		},
		{
			group.PATCH("/test", testAction),
			http.MethodPatch,
			"/test",
		},
		{
			group.PUT("/test", testAction),
			http.MethodPut,
			"/test",
		},
		{
			group.HEAD("/test", testAction),
			http.MethodHead,
			"/test",
		},
		{
			group.OPTIONS("/test", testAction),
			http.MethodOptions,
			"/test",
		},
	}

	for i, s := range scenarios {
		t.Run(fmt.Sprintf("%d_%s_%s", i, s.expectMethod, s.expectPath), func(t *testing.T) {
			if s.route.Method != s.expectMethod {
				t.Fatalf("Expected method %q, got %q", s.expectMethod, s.route.Method)
			}

			if s.route.Path != s.expectPath {
				t.Fatalf("Expected path %q, got %q", s.expectPath, s.route.Path)
			}

			if err := s.route.Action(nil); !errors.Is(err, testErr) {
				t.Fatal("Expected test action")
			}
		})
	}
}

func TestRouterGroupHasRoute(t *testing.T) {
	t.Parallel()

	group := RouterGroup[*Event]{}

	group.Any("/any", nil)

	group.GET("/base", nil)
	group.DELETE("/base", nil)

	sub := group.Group("/sub1")
	sub.GET("/a", nil)
	sub.POST("/a", nil)

	sub2 := sub.Group("/sub2")
	sub2.GET("/b", nil)
	sub2.GET("/b/{test}", nil)

	// special cases to test the normalizations
	group.GET("/c/", nil)          // the same as /c/{test...}
	group.GET("/d/{test...}", nil) // the same as /d/

	scenarios := []struct {
		method   string
		path     string
		expected bool
	}{
		{
			http.MethodGet,
			"",
			false,
		},
		{
			"",
			"/any",
			true,
		},
		{
			http.MethodPost,
			"/base",
			false,
		},
		{
			http.MethodGet,
			"/base",
			true,
		},
		{
			http.MethodDelete,
			"/base",
			true,
		},
		{
			http.MethodGet,
			"/sub1",
			false,
		},
		{
			http.MethodGet,
			"/sub1/a",
			true,
		},
		{
			http.MethodPost,
			"/sub1/a",
			true,
		},
		{
			http.MethodDelete,
			"/sub1/a",
			false,
		},
		{
			http.MethodGet,
			"/sub2/b",
			false,
		},
		{
			http.MethodGet,
			"/sub1/sub2/b",
			true,
		},
		{
			http.MethodGet,
			"/sub1/sub2/b/{test}",
			true,
		},
		{
			http.MethodGet,
			"/sub1/sub2/b/{test2}",
			false,
		},
		{
			http.MethodGet,
			"/c/{test...}",
			true,
		},
		{
			http.MethodGet,
			"/d/",
			true,
		},
	}

	for _, s := range scenarios {
		t.Run(s.method+"_"+s.path, func(t *testing.T) {
			has := group.HasRoute(s.method, s.path)

			if has != s.expected {
				t.Fatalf("Expected %v, got %v", s.expected, has)
			}
		})
	}
}
