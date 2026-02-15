package handlers

import "net/http"

func (s *Server) SendHandler(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello World!"))
}
