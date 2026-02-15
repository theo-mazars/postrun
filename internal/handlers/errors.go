package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/theo-mazars/postrun/internal/smtp"
)

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
