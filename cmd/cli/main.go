package main

import (
	"fmt"
	"net/mail"
	"os"

	"github.com/spf13/cobra"

	"github.com/theo-mazars/postrun/internal/config"
	"github.com/theo-mazars/postrun/internal/dns"
	"github.com/theo-mazars/postrun/internal/smtp"
)

var cfgFile string

var rootCmd = &cobra.Command{
  Use: "postrun",
  Short: "Transactional email SMTP sender",
  Long: "Lightweight transactional email sender with HTTP API",
}

var versionCmd = &cobra.Command{
  Use: "version",
  Short: "Print version",
  Run: func (cmd *cobra.Command, args []string) {
    fmt.Println("postrun v0.1.0")
  },
}

var sendCmd = &cobra.Command{
  Use: "send",
  Short: "Send an email",
  Run: func(cmd *cobra.Command, args []string) {
    cfg := config.Load(cfgFile)

    fromStr, _ := cmd.Flags().GetString("from")
    toStr, _ := cmd.Flags().GetString("to")
    subject, _ := cmd.Flags().GetString("subject")
    body, _ := cmd.Flags().GetString("body")

    from, err := mail.ParseAddress(fromStr)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Invalid From address: %v\n", err)
    }
    to, err := mail.ParseAddress(toStr)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Invalid To address: %v\n", err)
    }

    email := &smtp.Email{
      From: from,
      To: to,
      Subject: subject,
      Body: body,
    }

    sender := smtp.NewSender(cfg.SMTP.Domain)
    result, err := sender.Send(email)
    if err != nil {
      fmt.Fprintf(os.Stderr, "Send failed: %v\n", err)
      os.Exit(1)
    }

    fmt.Printf("Result: success=%v code=%d message=%s\n", result.Success, result.Code, result.Message)
  },
}

var mxCmd = &cobra.Command{
  Use: "mx [email]",
  Short: "Lookup MX records for email domain",
  Args: cobra.ExactArgs(1),
  Run: func(cmd *cobra.Command, args []string) {
    email := args[0]
    records, err := dns.LookupMX(email)
    if err != nil {
      fmt.Fprintf(os.Stderr, "MX lookup failed: %v\n", err)
      os.Exit(1)
    }
    for _, mx := range records {
      fmt.Printf("• %3d\t%s\n", mx.Pref, mx.Host)
    }
  },
}

func init() {
  rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path")
  rootCmd.AddCommand(versionCmd)
  rootCmd.AddCommand(sendCmd)
  rootCmd.AddCommand(mxCmd)
  sendCmd.Flags().String("to", "", "recipient email")
  sendCmd.Flags().String("from", "", "sender email")
  sendCmd.Flags().String("subject", "Test", "email subject")
  sendCmd.Flags().String("body", "Hello from postrung", "email body")
  sendCmd.MarkFlagRequired("to")
  sendCmd.MarkFlagRequired("from")
}

func main() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
