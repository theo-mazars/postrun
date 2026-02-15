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
  From    string `json:"from"`
	To      string `json:"to"`
	Subject string `json:"subject"`
	Text    string `json:"text"`
}

type ErrorResponse struct {
  Code string `json:"code"`
  Message string `json:"message"`
  Details string `json:"details"`
}

func (s *Server) SendHandler(w http.ResponseWriter, r *http.Request) {
  sendRequest := &SendRequest{}
  apiKey := r.Context().Value("api_key").(string)
  domainId := r.Context().Value("domain_id").(string)
  domain := r.Context().Value("domain").(string)

  err := json.NewDecoder(r.Body).Decode(&sendRequest)
  if err != nil {
    returnHTTPError(w, 400, &ErrorResponse{"bad_request", "Bad Request", "missing body"})
    return
  }

  email, err := sendRequest.validate()
  if err != nil {
    returnHTTPError(w, 400, &ErrorResponse{"bad_request", "Bad Request", err.Error()})
    return
  }
  if !strings.HasSuffix(email.From.Address, domain) {
    returnHTTPError(w, 401, &ErrorResponse{"bad_domain", "Bad Domain", fmt.Sprintf("%s: api key cannot send for this domain", domain)})
    return
  }

  sender := smtp.NewSender(domain)
  response, err := sender.Send(email)
  if err != nil {
    returnHTTPError(w, 500, &ErrorResponse{"internal_server_error", "Internal Server Error", "If you're admin, check logs"})
    fmt.Println(err)
    return
  }
  err = s.logEmail(response, email, apiKey, domainId)
  if err != nil {
    fmt.Println(err)
  }
  if response.Code != 250 {
    smtp2httpCode(w, response)
    return
  }

  w.WriteHeader(http.StatusOK)
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
    return nil, fmt.Errorf("to: %s", err)
  }

  if s.Subject == "" {
    return nil, fmt.Errorf("subject: empty subject is not allowed")
  }
  if s.Text == "" {
    return nil, fmt.Errorf("text: empty text is not allowed")
  }

  return &smtp.Email{
    Id: uuid.New(),
    From: from,
    To: to,
    Subject: s.Subject,
    Body: s.Text,
  }, nil
}
