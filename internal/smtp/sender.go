package smtp

import (
	"bufio"
	"fmt"
	"net"
	"net/mail"
	"time"

	"github.com/theo-mazars/postrun/internal/dns"
)

type Email struct {
	From    *mail.Address
	To      *mail.Address
	Subject string
	Body    string
}

type SMTPResponse struct {
	Success bool
	Code    int
	Message string
}

type Sender struct {
	helloDomain string
}

func NewSender(helloDomain string) *Sender {
	return &Sender{helloDomain}
}

func (s *Sender) Send(email *Email) (*SMTPResponse, error) {
  startTime := time.Now()
	mxRecords, err := dns.LookupMX(email.To.Address)
	if err != nil {
		return nil, err
	}

	for _, mx := range mxRecords {
		sess, err := s.openSession(mx.Host)
		if err != nil {
			fmt.Printf("%v", err)
			continue
		}

		code, msg, err := sess.init(s.helloDomain)
		if err != nil {
			return nil, err
		}
		if code != 250 {
			return &SMTPResponse{false, code, msg}, nil
		}

		res, err := sess.mailFrom(email.From.Address)
		if err != nil {
			return nil, err
		}
		if res != nil && !res.Success {
			return res, nil
		}

		res, err = sess.rcptTo(email.To.Address)
		if err != nil {
			return nil, err
		}
		if res != nil && !res.Success {
			return res, nil
		}

		res, err = sess.data(email)
		if err != nil {
			return nil, err
		}

		sess.close()
		fmt.Printf("* Elapsed time: %s\n", time.Since(startTime))
		return res, nil
	}

	return nil, fmt.Errorf("all MX hosts failed")
}

func (s *Sender) openSession(host string) (*session, error) {
	conn, err := net.DialTimeout("tcp", host+":25", 10*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}

	conn.SetDeadline(time.Now().Add(60 * time.Second))

	reader := bufio.NewReader(conn)
	writer := bufio.NewWriter(conn)

	sess := &session{host, conn, reader, writer}
	code, _, err := sess.readResponse()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("greeting failed: %w", err)
	}
	if code != 220 {
		conn.Close()
		return nil, fmt.Errorf("unexpected greeting: %d", code)
	}

	return sess, nil
}
