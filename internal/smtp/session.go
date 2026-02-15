package smtp

import (
	"bufio"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"strings"
	"time"
)

type session struct {
	host   string
	conn   net.Conn
	reader *bufio.Reader
	writer *bufio.Writer
}

func (sess *session) close() {
  err := sess.sendCmd("QUIT")
  if err != nil {
    fmt.Fprintf(os.Stderr, "cmd: QUIT: %v", err)
  }
 	code, msg, err := sess.readResponse()
	if code != 221 {
	  fmt.Fprintf(os.Stderr, "cmd: QUIT: %v: %v", err, msg)
	}
	sess.conn.Close()
}

func (sess *session) readResponse() (int, string, error) {
	var msg strings.Builder
	var code int

	for {
		line, err := sess.reader.ReadString('\n')
		if err != nil {
			return 0, "", err
		}

		line = strings.TrimSpace(line)
		if len(line) < 3 {
			return 0, "", fmt.Errorf("invalid response: %s", line)
		}

		fmt.Sscanf(line[:3], "%d", &code)
		msg.WriteString(line[4:] + "\n")

		if len(line) > 3 && line[3] == ' ' {
			fmt.Printf("<-- [%s] %s\n", line[:3], line[4:])
			break
		}
		fmt.Printf("<    %s  %s\n", line[:3], line[4:])
	}

	return code, strings.TrimSpace(msg.String()), nil
}

func (sess *session) sendCmd(cmd string) error {
	fmt.Printf("--> %s\n", cmd)

	_, err := sess.writer.WriteString(cmd + "\r\n")
	if err != nil {
		return err
	}

	return sess.writer.Flush()
}

func (sess *session) init(helloDomain string) (int, string, error) {
	code, msg, err := sess.ehlo(helloDomain)
	if err != nil {
		return code, msg, err
	}

	code, msg, err = sess.startTls(helloDomain)
	if err != nil {
		fmt.Printf("%d - %s\n", code, msg)
		return code, msg, err
	}

	return code, msg, nil
}

func (sess *session) ehlo(helloDomain string) (int, string, error) {
	if err := sess.sendCmd("EHLO " + helloDomain); err != nil {
		return 0, "", err
	}

	code, msg, err := sess.readResponse()
	if err != nil {
		return code, msg, err
	}

	return code, msg, nil
}

func (sess *session) startTls(helloDomain string) (int, string, error) {
	if err := sess.sendCmd("STARTTLS"); err != nil {
		return 0, "", err
	}

	code, msg, err := sess.readResponse()
	if err != nil {
		return code, msg, err
	}
	if code == 220 {
		tlsConn := tls.Client(sess.conn, &tls.Config{
			ServerName:         sess.host,
			InsecureSkipVerify: true,
		})
		if err := tlsConn.Handshake(); err != nil {
			return 0, "", fmt.Errorf("TLS handshake failed: %w", err)
		}

		sess.reader = bufio.NewReader(tlsConn)
		sess.writer = bufio.NewWriter(tlsConn)
		sess.conn = tlsConn

		return sess.ehlo(helloDomain)
	}

	return code, msg, nil
}

func (sess *session) mailFrom(from string) (*SMTPResponse, error) {
	if err := sess.sendCmd(fmt.Sprintf("MAIL FROM:<%s>", from)); err != nil {
		return nil, err
	}

	code, msg, err := sess.readResponse()
	if err != nil {
		return nil, err
	}

	if code != 250 {
		return &SMTPResponse{false, code, msg}, nil
	}
	return &SMTPResponse{true, code, msg}, nil
}

func (sess *session) rcptTo(to string) (*SMTPResponse, error) {
	if err := sess.sendCmd(fmt.Sprintf("RCPT TO:<%s>", to)); err != nil {
		return nil, err
	}

	code, msg, err := sess.readResponse()
	if err != nil {
		return nil, err
	}

	if code != 250 {
		return &SMTPResponse{false, code, msg}, nil
	}
	return &SMTPResponse{true, code, msg}, nil
}

func (sess *session) data(email *Email) (*SMTPResponse, error) {
	if err := sess.sendCmd("DATA"); err != nil {
		return nil, err
	}

	code, msg, err := sess.readResponse()
	if err != nil {
		return nil, err
	}
	if code != 354 {
		return &SMTPResponse{false, code, msg}, nil
	}

	message := buildMessage(email)
	if _, err := sess.writer.WriteString(message); err != nil {
		return nil, err
	}
	if err := sess.writer.Flush(); err != nil {
		return nil, err
	}

	if err := sess.sendCmd("\r\n."); err != nil {
		return nil, err
	}

	code, msg, err = sess.readResponse()
	if err != nil {
		return nil, err
	}
	if code != 250 {
		return &SMTPResponse{false, code, msg}, nil
	}

	return &SMTPResponse{true, code, msg}, nil
}

func buildMessage(email *Email) string {
	var msg strings.Builder

	// Email Headers
	msg.WriteString(fmt.Sprintf("From: %s\r\n", email.From))
	msg.WriteString(fmt.Sprintf("To: %s\r\n", email.To))
	msg.WriteString(fmt.Sprintf("Subject: %s\r\n", email.Subject))
	msg.WriteString(fmt.Sprintf("Message-ID: <%s.%s>\r\n", email.Id.String(), email.From.Address))
	msg.WriteString(fmt.Sprintf("Date: %s\r\n", time.Now().Format(time.RFC1123Z)))
	msg.WriteString("MIME-Version: 1.0\r\n")
	msg.WriteString("Content-Type: text/plain; charset=UTF-8\r\n")
	msg.WriteString("\r\n")

	// Email Body
	for line := range strings.SplitSeq(email.Body, "\n") {
		if strings.HasPrefix(line, ".") {
			msg.WriteString(".")
		}
		msg.WriteString(line + "\r\n")
	}

	return msg.String()
}
