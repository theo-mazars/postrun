package dns

import (
  "fmt"
  "net"
  "sort"
  "strings"
)

func LookupMX(email string) ([]net.MX, error) {
  // TODO: properly handle email address formatting
  parts := strings.Split(email, "@")
  if len(parts) != 2 {
    return nil, fmt.Errorf("invalid email format")
  }

  domain := parts[1]

  if domain == "localhost" {
    records := make([]net.MX, 1)
    records[0] = net.MX{
      Host: "localhost",
      Pref: 0,
    }
    return records, nil
  }

  mxRecords, err := net.LookupMX(domain)
  if err != nil {
    return nil, fmt.Errorf("MX Lookup failed: %w", err)
  }
  if len(mxRecords) == 0 {
    return nil, fmt.Errorf("No MX records found for %s", domain)
  }

  records := make([]net.MX, len(mxRecords))
  for i, mx := range mxRecords {
    records[i] = net.MX{
      Host: strings.TrimSuffix(mx.Host, "."),
      Pref: mx.Pref,
    }
  }

  sort.Slice(records, func(i int, j int) bool {
    return records[i].Pref < records[j].Pref
  })

  return records, nil
}
