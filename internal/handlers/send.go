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
    smtp2httpCode(w, response)
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

func returnHTTPError(w http.ResponseWriter, status int, response *ErrorResponse) {
  body, err := json.Marshal(response)
  if err != nil {
    w.WriteHeader(http.StatusInternalServerError)
    return
  }
  w.Header().Set("Content-Type", "application/json")
  w.WriteHeader(status)
  w.Write(body)
}

func smtp2httpCode(w http.ResponseWriter, response *smtp.SMTPResponse) {
  details := fmt.Sprintf("[%d]: %s", response.Code, response.Message)

  switch response.Code {
  case 421:
    returnHTTPError(w, 503, &ErrorResponse{"service_unavailable", "SMTP Service not available", details})
  case 450:
    returnHTTPError(w, 503, &ErrorResponse{"mailbox_unavailable", "Mailbox busy or temporarily blocked", details})
  case 451:
    returnHTTPError(w, 503, &ErrorResponse{"requested_action_aborted", "Local error", details})
  case 452:
    returnHTTPError(w, 422, &ErrorResponse{"mailbox_full", "Insufficient storage", details})
  case 454:
    returnHTTPError(w, 503, &ErrorResponse{"temporary_failure", "Temporary authentication failure", details})
  case 500:
    returnHTTPError(w, 502, &ErrorResponse{"syntax_error", "Command unrecognized", details})
  case 501:
    returnHTTPError(w, 502, &ErrorResponse{"syntax_error", "Parameter or Argument error", details})
  case 502:
    returnHTTPError(w, 502, &ErrorResponse{"command_unimplemented", "Command not implemented", details})
  case 503:
    returnHTTPError(w, 502, &ErrorResponse{"command_bad_sequence", "Bad sequence of commands", details})
  case 504:
    returnHTTPError(w, 502, &ErrorResponse{"command_parameter_unimplemented", "Command parameter not implemented", details})
  case 521:
    returnHTTPError(w, 400, &ErrorResponse{"mail_not_accepted", "Server does not accept mails", details})
  case 530:
    returnHTTPError(w, 401, &ErrorResponse{"auth_required", "Authentication required", details})
  case 550:
    returnHTTPError(w, 400, &ErrorResponse{"mailbox_not_found", "Mailbox not found", details})
  case 551:
    returnHTTPError(w, 400, &ErrorResponse{"user_not_local", "User not local", details})
  case 552:
    returnHTTPError(w, 422, &ErrorResponse{"exceed_storage", "Exceeded storage allocation", details})
  case 553:
    returnHTTPError(w, 400, &ErrorResponse{"mailbox_not_allowed", "Mailbox name not allowed", details})
  case 554:
    returnHTTPError(w, 422, &ErrorResponse{"transaction_failed", "spam/policy rejection", details})
  case 556:
    returnHTTPError(w, 400, &ErrorResponse{"domain_reject", "Domain does not accept mail", details})
  default:
    returnHTTPError(w, 500, &ErrorResponse{"internal_server_error", "Unknown SMTP Error", details})
  }
}
