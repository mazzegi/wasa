package devutil

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pkg/errors"
)

type ServerConfig struct {
	Bind     string
	MainGo   string
	DistDir  string
	WasmExec string
}

type Server struct {
	config     ServerConfig
	httpServer *http.Server
}

func NewServer(c ServerConfig) (*Server, error) {
	err := os.MkdirAll(c.DistDir, os.ModePerm)
	if err != nil {
		return nil, err
	}
	bs, err := exec.Command("cp", c.WasmExec, c.DistDir).CombinedOutput()
	if err != nil {
		return nil, errors.Errorf("cp wasm_exec.js: %s", string(bs))
	}

	s := &Server{
		config: c,
		httpServer: &http.Server{
			Addr: c.Bind,
		},
	}
	return s, nil
}

func (s *Server) ListenAndServe() error {

	mux := http.NewServeMux()
	mux.HandleFunc("/", s.serveFiles)
	s.httpServer.Handler = mux
	fmt.Printf("start listening on (%s) (mgo=%s) (dist=%s)\n", s.config.Bind, s.config.MainGo, s.config.DistDir)
	return s.httpServer.ListenAndServe()
}

func (s *Server) Close() {
	s.httpServer.Close()
}

func (s *Server) serveFiles(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" || r.URL.Path == "/index.html" {
		s.serveIndex(w, r)
		return
	} else if r.URL.Path == "/"+wasmFile {
		s.serveWASM(w, r)
		return
	}
	http.FileServer(http.Dir(s.config.DistDir)).ServeHTTP(w, r)
}

func (s *Server) serveWASM(w http.ResponseWriter, r *http.Request) {
	wasmPath := filepath.Join(s.config.DistDir, wasmFile)
	err := buildWASM(s.config.MainGo, wasmPath)
	if err != nil {
		serr := fmt.Sprintf("error building wasm: %v", err)
		fmt.Println(err)
		http.Error(w, serr, http.StatusInternalServerError)
		return
	}
	http.ServeFile(w, r, wasmPath)
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index.html").Parse(indexHTMLTemplate)
	if err != nil {
		serr := fmt.Sprintf("error building index.html: %v", err)
		fmt.Println(err)
		http.Error(w, serr, http.StatusInternalServerError)
		return
	}
	data := struct {
	}{}
	t.Execute(w, data)
}
