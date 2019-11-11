package weather

import (
	"fmt"
	"net/http"
)

func (wa *WebApi) City(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "City Hello, %s!", r.URL.Path[1:])
}

func (wa *WebApi) Detail(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Detail Hello, %s!", r.URL.Path[1:])
}
