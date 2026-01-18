package utils

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
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

func CMDRun(main string, args ...string) error {
	cmd := exec.Command(main, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func CMDOutput(main string, args ...string) (string, error) {
	out, err := exec.Command(main, args...).Output()
	return string(out), err
}

const (
	remoteServer = "podrun@10.7.22.101"
	password     = "passwd"
)

func SSHRun(args ...string) error {
	command := strings.Join(args, " ")
	cmdArgs := []string{
		"-p", password,
		"ssh",
		"-tt",
		"-o", "StrictHostKeyChecking=no",
		"-o", "LogLevel=QUIET",
		remoteServer,
		command,
	}
	cmd := exec.Command("sshpass", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

func SSEOutput(args ...string) (string, error) {
	command := strings.Join(args, " ")
	cmdArgs := []string{
		"-p", password,
		"ssh",
		"-o", "StrictHostKeyChecking=no",
		"-o", "LogLevel=QUIET",
		remoteServer,
		command,
	}
	return CMDOutput("sshpass", cmdArgs...)
}

func GetHostName() string {
	if host, err := os.Hostname(); err == nil {
		return host
	}
	return "unknown"
}
