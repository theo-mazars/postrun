package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/theo-mazars/postrun/internal/smtp"
)

type SendRequest struct {
  From    string `json:"from" validate:"nonzero"`
	To      string `json:"to" validate:"nonzero"`
	Subject string `json:"subject" validate:"nonzero"`
	Text    string `json:"text" validate:"nonzero"`
}

func (s *Server) SendHandler(w http.ResponseWriter, r *http.Request) {
  sendRequest := &SendRequest{}

  err := json.NewDecoder(r.Body).Decode(&sendRequest)
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    return
  }

  email, err := sendRequest.validate()
  if err != nil {
    w.WriteHeader(http.StatusBadRequest)
    fmt.Println(err)
    return
  }
  if !strings.HasSuffix(email.From.Address, r.Context().Value("domain").(string)) {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  sender := smtp.NewSender(r.Context().Value("domain").(string))
  response, err := sender.Send(email)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Println(err)
    return
  }
  if response.Code != 250 {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte(response.Message))
    return
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(response.Message))
}

func (s *SendRequest) validate() (*smtp.Email, error) {
  from, err := mail.ParseAddress(s.From)
  if err != nil {
    return nil, fmt.Errorf("from: %s", err)
  }

  to, err := mail.ParseAddress(s.To)
  if err != nil {
    return nil, fmt.Errorf("from: %s", err)
  }

  return &smtp.Email{
    From: from,
    To: to,
    Subject: s.Subject,
    Body: s.Text,
  }, nil
}
