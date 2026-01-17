package utils

import (
	"fmt"
	"net"
	"os"
	"os/exec"
)

func IsDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

func FileExist(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func GetMAC() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, e := range ifaces {
		if e.Flags&net.FlagLoopback == 0 && len(e.HardwareAddr) > 0 {
			return e.HardwareAddr.String(), nil
		}
	}
	return "", fmt.Errorf("no valid MAC address found")
}

func CmdExec(main string, args ...string) error {
	cmd := exec.Command(main, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}
