package main

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

//go:embed frontend/dist
var frontendAssets embed.FS

// --- Post model ---

type Post struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
}

type PostStore struct {
	mu    sync.RWMutex
	posts map[int]Post
	next  int
}

func NewPostStore() *PostStore {
	return &PostStore{
		posts: map[int]Post{},
		next:  1,
	}
}

func (s *PostStore) List() []Post {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]Post, 0, len(s.posts))
	for i := len(s.posts); i >= 1; i-- {
		if p, ok := s.posts[i]; ok {
			result = append(result, p)
		}
	}
	return result
}

func (s *PostStore) Get(id int) (Post, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	p, ok := s.posts[id]
	return p, ok
}

func (s *PostStore) Create(title, content, author string) Post {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.next
	s.next++
	p := Post{
		ID:        id,
		Title:     title,
		Content:   content,
		Author:    author,
		CreatedAt: time.Now(),
	}
	s.posts[id] = p
	return p
}

func (s *PostStore) Update(id int, title, content string) (Post, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	p, ok := s.posts[id]
	if !ok {
		return Post{}, false
	}
	p.Title = title
	p.Content = content
	s.posts[id] = p
	return p, true
}

func (s *PostStore) Delete(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, ok := s.posts[id]
	if ok {
		delete(s.posts, id)
	}
	return ok
}

// --- Helpers ---

type apiError struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func extractID(path string) (int, error) {
	parts := strings.Split(strings.TrimRight(path, "/"), "/")
	if len(parts) == 0 {
		return 0, fmt.Errorf("no id")
	}
	return strconv.Atoi(parts[len(parts)-1])
}

// --- Main ---

func main() {
	store := NewPostStore()

	// Seed data
	store.Create("환영합니다!", "test-runner 게시판이 오픈했습니다.\n자유롭게 글을 남겨주세요.", "운영자")

	mux := http.NewServeMux()

	// Health
	mux.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, 200, map[string]string{"status": "healthy", "project": "test-runner"})
	})

	// Posts CRUD
	mux.HandleFunc("/api/posts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			writeJSON(w, 200, store.List())
		case "POST":
			var p struct {
				Title   string `json:"title"`
				Content string `json:"content"`
				Author  string `json:"author"`
			}
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				writeJSON(w, 400, apiError{"invalid request body"})
				return
			}
			if p.Title == "" {
				writeJSON(w, 400, apiError{"title is required"})
				return
			}
			if p.Author == "" {
				p.Author = "익명"
			}
			post := store.Create(p.Title, p.Content, p.Author)
			writeJSON(w, 201, post)
		default:
			writeJSON(w, 405, apiError{"method not allowed"})
		}
	})

	// Single post by ID
	mux.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		id, err := extractID(r.URL.Path)
		if err != nil {
			writeJSON(w, 400, apiError{"invalid id"})
			return
		}

		switch r.Method {
		case "GET":
			post, ok := store.Get(id)
			if !ok {
				writeJSON(w, 404, apiError{"post not found"})
				return
			}
			writeJSON(w, 200, post)

		case "PUT":
			var p struct {
				Title   string `json:"title"`
				Content string `json:"content"`
			}
			if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
				writeJSON(w, 400, apiError{"invalid request body"})
				return
			}
			post, ok := store.Update(id, p.Title, p.Content)
			if !ok {
				writeJSON(w, 404, apiError{"post not found"})
				return
			}
			writeJSON(w, 200, post)

		case "DELETE":
			if store.Delete(id) {
				writeJSON(w, 200, map[string]string{"status": "deleted"})
			} else {
				writeJSON(w, 404, apiError{"post not found"})
			}

		default:
			writeJSON(w, 405, apiError{"method not allowed"})
		}
	})

	// Frontend SPA
	frontendFS, err := fs.Sub(frontendAssets, "frontend/dist")
	if err != nil {
		log.Fatalf("Failed to load frontend assets: %v", err)
	}
	mux.Handle("/", http.FileServer(http.FS(frontendFS)))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
