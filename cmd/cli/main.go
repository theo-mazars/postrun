package main
import (
  "fmt"
  "os"

  "github.com/spf13/cobra"

  "github.com/theo-mazars/postrun/internal/config"
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

func init() {
  rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file path")
  rootCmd.AddCommand(versionCmd)
  rootCmd.AddCommand(sendCmd)
}

func main() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
