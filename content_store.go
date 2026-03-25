package main

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

type contentSnapshot struct {
	Projects       []Project
	Certifications []Certification
	LoadedAt       time.Time
	Ready          bool
	Refreshing     bool
}

type contentStore struct {
	mu         sync.RWMutex
	snapshot   contentSnapshot
	refreshing atomic.Bool
}

var appContentStore = &contentStore{}

func cloneProjects(in []Project) []Project {
	if len(in) == 0 {
		return nil
	}
	out := make([]Project, len(in))
	copy(out, in)
	return out
}

func cloneCerts(in []Certification) []Certification {
	if len(in) == 0 {
		return nil
	}
	out := make([]Certification, len(in))
	copy(out, in)
	return out
}

func (c *contentStore) Snapshot() contentSnapshot {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return contentSnapshot{
		Projects:       cloneProjects(c.snapshot.Projects),
		Certifications: cloneCerts(c.snapshot.Certifications),
		LoadedAt:       c.snapshot.LoadedAt,
		Ready:          c.snapshot.Ready,
		Refreshing:     c.refreshing.Load(),
	}
}

func (c *contentStore) Refresh() {
	if !c.refreshing.CompareAndSwap(false, true) {
		return
	}
	defer c.refreshing.Store(false)

	var projects []Project
	var certs []Certification

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		projects = fetchProjectsFromNotion()
	}()
	go func() {
		defer wg.Done()
		certs = fetchCertificationsFromNotion()
	}()
	wg.Wait()

	c.mu.Lock()
	c.snapshot.Projects = cloneProjects(projects)
	c.snapshot.Certifications = cloneCerts(certs)
	c.snapshot.LoadedAt = time.Now()
	c.snapshot.Ready = true
	c.mu.Unlock()
}

func (c *contentStore) WarmUp(timeout time.Duration) {
	if timeout <= 0 {
		c.Refresh()
		return
	}

	done := make(chan struct{})
	go func() {
		defer close(done)
		c.Refresh()
	}()

	t := time.NewTimer(timeout)
	defer t.Stop()

	select {
	case <-done:
	case <-t.C:
	}
}

func (c *contentStore) StartAutoRefresh(ctx context.Context, interval time.Duration) {
	if interval <= 0 {
		return
	}

	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				c.Refresh()
			}
		}
	}()
}
