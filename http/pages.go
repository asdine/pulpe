package http

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/blankrobot/pulpe"
	"github.com/julienschmidt/httprouter"
)

// RegisterPageHandler register the routes for serving pages.
func RegisterPageHandler(mux *ServeMux, connect Connector, dir string, lazy bool) {
	pattern := filepath.Join(dir, "*.tmpl.html")

	h := pageHandler{
		lazy:    lazy,
		dir:     dir,
		connect: connect,
	}

	if !lazy {
		h.templates = template.Must(template.ParseGlob(pattern))
	}

	mux.HandleFunc("/join", h.handleRegister)
	mux.HandleFunc("/login", h.handleLogin)
	mux.HandleFunc("/logout", h.handleLogout)

	router := httprouter.New()
	router.GET("/:owner", h.handleBoardPage)
	router.GET("/:owner/:board", h.handleBoardPage)
	router.GET("/:owner/:board/:list/:card", h.handleBoardPage)
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			h.handleIndex(w, r)
			return
		}

		router.ServeHTTP(w, r)
	})
}

type pageHandler struct {
	templates *template.Template
	lazy      bool
	dir       string
	connect   Connector
}

func (h *pageHandler) render(wr io.Writer, name string, data interface{}) {
	if h.lazy {
		template.Must(template.ParseFiles(filepath.Join(h.dir, name))).Execute(wr, data)
	} else {
		h.templates.ExecuteTemplate(wr, name, data)
	}
}

func (h *pageHandler) handleIndex(w http.ResponseWriter, r *http.Request) {
	session := h.connect(r)
	defer session.Close()

	user, err := session.Authenticate()
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	http.Redirect(w, r, "/"+user.Login, http.StatusFound)
}

func (h *pageHandler) handleRegister(w http.ResponseWriter, r *http.Request) {
	h.render(w, "register.tmpl.html", map[string]string{
		"Title": "Join",
	})
}

func (h *pageHandler) handleLogin(w http.ResponseWriter, r *http.Request) {
	h.render(w, "login.tmpl.html", map[string]string{
		"Title": "Sign in",
	})
}

func (h *pageHandler) handleLogout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("pulpesid")
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	session := h.connect(r)
	defer session.Close()

	err = session.UserSessionService().DeleteSession(cookie.Value)
	if err != nil && err != pulpe.ErrUserSessionUnknownID {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "pulpesid",
		Expires: time.Now().UTC(),
		Path:    "/",
	})

	http.Redirect(w, r, "/login", http.StatusFound)
}

func (h *pageHandler) handleBoardPage(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	session := h.connect(r)
	defer session.Close()

	user, err := session.Authenticate()
	if err != nil {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	if user.Login != ps.ByName("owner") {
		http.NotFound(w, r)
		return
	}

	if ps.ByName("board") != "" {
		h.render(w, "board.tmpl.html", map[string]string{
			"Title": "",
		})
		return
	}

	boards, err := session.BoardService().Boards()
	if err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if len(boards) == 0 {
		h.render(w, "board.tmpl.html", map[string]string{
			"Title": "",
		})
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/%s/%s", user.Login, boards[0].Slug), http.StatusFound)
}
