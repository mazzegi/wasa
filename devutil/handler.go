package devutil

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

type MakeServeHandler struct {
	src  string
	dist string
	wasm string
}

func NewMakeServeHandler(src, dist, wasm string) (*MakeServeHandler, error) {
	//build once
	err := Make(src, dist)
	if err != nil {
		return nil, errors.Wrap(err, "initial build")
	}

	h := &MakeServeHandler{
		src:  src,
		dist: dist,
		wasm: wasm,
	}
	return h, nil
}

func (h *MakeServeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, h.wasm) {
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
