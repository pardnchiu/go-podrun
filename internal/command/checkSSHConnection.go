package command

import (
	"github.com/pardnchiu/go-podrun/internal/utils"
)

const (
	RemoteServer = "podrun@10.7.22.101"
	Password     = "passwd"
)

func CheckSSHConnection() error {
	if err := utils.CMDRun("sshpass",
		"-p", Password,
		"ssh",
		"-o", "ConnectTimeout=3",
		"-o", "StrictHostKeyChecking=no",
		"-q", RemoteServer,
		"exit"); err != nil {
		return err
	}
	return nil
}
