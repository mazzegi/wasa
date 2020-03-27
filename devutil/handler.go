package devutil

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type HandlerOption func(h *MakeServeHandler) error

func WithHook(fnc func()) HandlerOption {
	return func(h *MakeServeHandler) error {
		h.hooks = append(h.hooks, fnc)
		return nil
	}
}

type MakeServeHandler struct {
	src   string
	dist  string
	wasm  string
	hooks []func()
}

func NewMakeServeHandler(src, dist, wasm string, opts ...HandlerOption) (*MakeServeHandler, error) {
	h := &MakeServeHandler{
		src:  src,
		dist: dist,
		wasm: wasm,
	}
	for _, opt := range opts {
		err := opt(h)
		if err != nil {
			return nil, err
		}
	}
	for _, hf := range h.hooks {
		hf()
	}
	//build once
	err := Make(src, dist)
	if err != nil {
		return nil, errors.Wrap(err, "initial build")
	}
	return h, nil
}

func (h *MakeServeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, h.wasm) {
		for _, hf := range h.hooks {
			hf()
		}
		err := Make(h.src, h.dist)
		if err != nil {
			serr := fmt.Sprintf("error building wasm: %v", err)
			fmt.Println(err)
			http.Error(w, serr, http.StatusInternalServerError)
			return
		}
	}
	http.FileServer(http.Dir(h.dist)).ServeHTTP(w, r)
}
