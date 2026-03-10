package registry

import (
	"context"
	"testing"

	"github.com/go-kratos/kratos/v2/registry"
)

func TestNew(t *testing.T) {
	r := New()
	if r == nil {
		t.Fatal("New() returned nil")
	}
}

func TestEndpoint(t *testing.T) {
	r := New()
	cases := []string{
		"order-svc:9000",
		"user-svc:8080",
		"",
	}
	for _, tc := range cases {
		if got := r.Endpoint(tc); got != tc {
			t.Errorf("Endpoint(%q) = %q, want %q", tc, got, tc)
		}
	}
}

func TestRegister(t *testing.T) {
	r := New()
	err := r.Register(context.Background(), &registry.ServiceInstance{ID: "test"})
	if err != nil {
		t.Errorf("Register() returned error: %v", err)
	}
}

func TestDeregister(t *testing.T) {
	r := New()
	err := r.Deregister(context.Background(), &registry.ServiceInstance{ID: "test"})
	if err != nil {
		t.Errorf("Deregister() returned error: %v", err)
	}
}

func TestGetService(t *testing.T) {
	r := New()
	instances, err := r.GetService(context.Background(), "order-svc")
	if err != nil {
		t.Errorf("GetService() returned error: %v", err)
	}
	if instances != nil {
		t.Errorf("GetService() = %v, want nil", instances)
	}
}

func TestWatch(t *testing.T) {
	r := New()
	w, err := r.Watch(context.Background(), "order-svc")
	if err != nil {
		t.Errorf("Watch() returned error: %v", err)
	}
	if w == nil {
		t.Fatal("Watch() returned nil watcher")
	}

	instances, err := w.Next()
	if err != nil {
		t.Errorf("watcher.Next() returned error: %v", err)
	}
	if instances != nil {
		t.Errorf("watcher.Next() = %v, want nil", instances)
	}

	if err := w.Stop(); err != nil {
		t.Errorf("watcher.Stop() returned error: %v", err)
	}
}
