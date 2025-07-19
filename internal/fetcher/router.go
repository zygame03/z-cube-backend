package fetcher

import (
	"errors"
	"fmt"
	"sync"
)

var (
	ErrRouteNotFound  = errors.New("route not found")
	ErrRouteDuplicate = errors.New("route already exists")
	ErrRouteInvalid   = errors.New("invalid route")
)

type RouteView struct {
	Name string
	URL  string
}

type Router struct {
	baseURL      string
	routes       map[string]*Route
	allRoutesNum int
	enabledNum   int
	mu           sync.RWMutex
}

func NewRouter(url string) *Router {
	return &Router{
		baseURL: url,
		routes:  make(map[string]*Route),
	}
}

func (r *Router) EnabledCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()

	return r.enabledNum
}

func (r *Router) AddRoutes(routes []*Route) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	addedCount := 0
	skippedCount := 0

	for _, route := range routes {
		if !validateRoute(route) {
			skippedCount++
			continue
		}
		if _, ok := r.routes[route.Name]; ok {
			skippedCount++
			continue
		}

		r.routes[route.Name] = route
		r.allRoutesNum++
		if route.Enabled {
			r.enabledNum++
		}
		addedCount++
	}

	return nil
}

func (r *Router) AddRoute(route *Route) error {
	if !validateRoute(route) {
		return ErrRouteInvalid
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.routes[route.Name]
	if ok {
		return nil
	}
	r.routes[route.Name] = route
	r.allRoutesNum++
	if route.Enabled {
		r.enabledNum++
	}

	return nil
}

func (r *Router) FetchableRoutes() []*RouteView {
	r.mu.RLock()
	defer r.mu.RUnlock()
	var list []*RouteView
	for _, v := range r.routes {
		if v.Enabled {
			list = append(list, &RouteView{
				Name: v.Name,
				URL:  fmt.Sprintf("%s%s", r.baseURL, v.Path),
			})
		}
	}
	return list
}

func (r *Router) Enable(path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	v, ok := r.routes[path]
	if ok {
		if !v.Enabled {
			v.Enabled = true
			r.enabledNum++
		}
		return nil
	}
	return ErrRouteNotFound
}

func (r *Router) Disable(path string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	v, ok := r.routes[path]
	if ok {
		if v.Enabled {
			v.Enabled = false
			r.enabledNum--
		}
		return nil
	}
	return ErrRouteNotFound
}

func validateRoute(route *Route) bool {
	if route.Name == "" || route.Path == "" {
		return false
	}
	return true
}
