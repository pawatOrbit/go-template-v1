package cmd

import (
	"fmt"
	"net"

	"github.com/spf13/cobra"
)

type BuildInfo struct {
	Service     string
	ServiceName string
	Version     string
}

func InitServeCommandGroup(rootCmd *cobra.Command) {
	rootCmd.AddGroup(&cobra.Group{ID: "serve", Title: "Serve:"})
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("local IP not found")
}
