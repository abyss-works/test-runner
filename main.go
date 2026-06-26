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

// extractPostIDAndRest returns the post ID and the remaining path after /api/posts/{id}
func extractPostIDAndRest(path string) (int, string, error) {
	trimmed := strings.TrimRight(path, "/")
	prefix := "/api/posts"
	if !strings.HasPrefix(trimmed, prefix) {
		return 0, "", fmt.Errorf("invalid path")
	}
	rest := strings.TrimPrefix(trimmed, prefix) // e.g. "/3" or "/3/comments" or "/3/comments/5"
	rest = strings.TrimLeft(rest, "/")
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) == 0 || parts[0] == "" {
		return 0, "", fmt.Errorf("no id")
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", err
	}
	suffix := ""
	if len(parts) > 1 {
		suffix = "/" + parts[1]
	}
	return id, suffix, nil
}

// --- Comment model ---

type Comment struct {
	ID        int       `json:"id"`
	PostID    int       `json:"postId"`
	Content   string    `json:"content"`
	Author    string    `json:"author"`
	CreatedAt time.Time `json:"createdAt"`
}

type CommentStore struct {
	mu       sync.RWMutex
	comments map[int][]Comment // postID -> comments
	next     int
}

func NewCommentStore() *CommentStore {
	return &CommentStore{
		comments: map[int][]Comment{},
		next:     1,
	}
}

func (s *CommentStore) List(postID int) []Comment {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.comments[postID]
}

func (s *CommentStore) Create(postID int, content, author string) Comment {
	s.mu.Lock()
	defer s.mu.Unlock()
	id := s.next
	s.next++
	c := Comment{
		ID:        id,
		PostID:    postID,
		Content:   content,
		Author:    author,
		CreatedAt: time.Now(),
	}
	s.comments[postID] = append(s.comments[postID], c)
	return c
}

func (s *CommentStore) Delete(postID, commentID int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	cs, ok := s.comments[postID]
	if !ok {
		return false
	}
	for i, c := range cs {
		if c.ID == commentID {
			s.comments[postID] = append(cs[:i], cs[i+1:]...)
			return true
		}
	}
	return false
}

// --- Main ---

func main() {
	store := NewPostStore()
	comments := NewCommentStore()

	// Seed data
	p := store.Create("환영합니다!", "test-runner 게시판이 오픈했습니다.\n자유롭게 글을 남겨주세요.", "운영자")
	comments.Create(p.ID, "첫 댓글입니다! 게시판 잘 사용하겠습니다.", "방문객")

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

	// Single post by ID + Comments
	mux.HandleFunc("/api/posts/", func(w http.ResponseWriter, r *http.Request) {
		postID, suffix, err := extractPostIDAndRest(r.URL.Path)
		if err != nil {
			writeJSON(w, 400, apiError{"invalid path"})
			return
		}

		// Comment routes: /api/posts/{id}/comments [/ {commentId}]
		if strings.HasPrefix(suffix, "/comments") {
			commentPath := strings.TrimPrefix(suffix, "/comments") // "" or "/{commentId}"

			if commentPath == "" {
				// GET /api/posts/{id}/comments — list
				// POST /api/posts/{id}/comments — create
				switch r.Method {
				case "GET":
					writeJSON(w, 200, comments.List(postID))
				case "POST":
					var c struct {
						Content string `json:"content"`
						Author  string `json:"author"`
					}
					if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
						writeJSON(w, 400, apiError{"invalid request body"})
						return
					}
					if c.Content == "" {
						writeJSON(w, 400, apiError{"content is required"})
						return
					}
					if c.Author == "" {
						c.Author = "익명"
					}
					// Verify post exists
					if _, ok := store.Get(postID); !ok {
						writeJSON(w, 404, apiError{"post not found"})
						return
					}
					comment := comments.Create(postID, c.Content, c.Author)
					writeJSON(w, 201, comment)
				default:
					writeJSON(w, 405, apiError{"method not allowed"})
				}
				return
			}

			// DELETE /api/posts/{id}/comments/{commentId}
			commentPath = strings.TrimLeft(commentPath, "/")
			commentID, err := strconv.Atoi(commentPath)
			if err != nil {
				writeJSON(w, 400, apiError{"invalid comment id"})
				return
			}
			if r.Method == "DELETE" {
				if comments.Delete(postID, commentID) {
					writeJSON(w, 200, map[string]string{"status": "deleted"})
				} else {
					writeJSON(w, 404, apiError{"comment not found"})
				}
			} else {
				writeJSON(w, 405, apiError{"method not allowed"})
			}
			return
		}

		// Single post routes: /api/posts/{id}
		switch r.Method {
		case "GET":
			post, ok := store.Get(postID)
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
			post, ok := store.Update(postID, p.Title, p.Content)
			if !ok {
				writeJSON(w, 404, apiError{"post not found"})
				return
			}
			writeJSON(w, 200, post)

		case "DELETE":
			if store.Delete(postID) {
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
