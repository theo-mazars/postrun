package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
	"strings"

	"github.com/google/uuid"

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
  apiKey := r.Context().Value("api_key").(string)
  domainId := r.Context().Value("domain_id").(string)
  domain := r.Context().Value("domain").(string)

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
  if !strings.HasSuffix(email.From.Address, domain) {
    w.WriteHeader(http.StatusUnauthorized)
    return
  }

  sender := smtp.NewSender(domain)
  response, err := sender.Send(email)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    fmt.Println(err)
    return
  }
  err = s.logEmail(response, email, apiKey, domainId)
  if err != nil {
    fmt.Println(err)
  }
  if response.Code != 250 {
    w.WriteHeader(http.StatusBadRequest)
    w.Write([]byte(response.Message))
    return
  }

  w.WriteHeader(http.StatusOK)
  w.Write([]byte(response.Message))
}

func (s *Server) logEmail(status *smtp.SMTPResponse, email *smtp.Email, key string, domain string) error {
  _, err := s.pool.Exec(`
    INSERT INTO emails
      (domain_id, api_key_id, message_id, from_address, to_address, subject, smtp_response_code, smtp_response_text)
    VALUES
      ($1, $2, $3, $4, $5, $6, $7, $8)
  `,
    domain,
    key,
    email.Id,
    email.From.Address,
    email.To.Address,
    email.Subject,
    status.Code,
    status.Message,
  )
  return err
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
    Id: uuid.New(),
    From: from,
    To: to,
    Subject: s.Subject,
    Body: s.Text,
  }, nil
}
