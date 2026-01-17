package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func commandExec(main string, args ...string) error {
	cmd := exec.Command(main, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func installPackage(pkg string, args ...string) error {
	if _, err := exec.LookPath(pkg); err != nil {
		return err
	}

	err := commandExec(pkg, args...)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	relyPackages := []string{"sshpass", "rsync", "ssh", "curl", "unzip"}
	var missPackages []string

	for _, e := range relyPackages {
		if _, err := exec.LookPath(e); err != nil {
			missPackages = append(missPackages, e)
		}
	}

	if len(missPackages) > 0 {
		fmt.Println("[!] missing packages:", strings.Join(missPackages, ", "))

		goos := runtime.GOOS
		if goos != "linux" && goos != "darwin" {
			fmt.Println("[x] only support RHEL / Debian and macOS")
			return
		}

		switch goos {
		case "darwin":
			if _, err := exec.LookPath("brew"); err != nil {
				// /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
				if _, err := exec.Command("bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)").Output(); err != nil {
					fmt.Println("[x] failed to install brew", err)
					return
				}
			}
			args := append([]string{"install", "-y"}, missPackages...)
			err := installPackage("brew", args...)
			if err != nil {
				fmt.Println("[x] failed to install packages", err)
				return
			}
		case "linux":
			args := append([]string{"install", "-y"}, missPackages...)
			err := installPackage("apt", args...)
			if err != nil {
				fmt.Println("[x] failed to install packages", err)
				return
			}
			err = installPackage("yum", args...)
			if err != nil {
				fmt.Println("[x] failed to install packages", err)
				return
			}
			err = installPackage("dnf", args...)
			if err != nil {
				fmt.Println("[x] failed to install packages", err)
				return
			}
		}
		return
	}
}
