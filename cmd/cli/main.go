package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/theo-mazars/postrun/internal/config"
	"github.com/theo-mazars/postrun/internal/dns"
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
    fmt.Printf("Loaded Config: %+v\n", cfg)
    fmt.Println("TODO: actually send email")
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
}

func main() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
