package command

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/pardnchiu/go-podrun/internal/model"
	"github.com/pardnchiu/go-podrun/internal/utils"
)

const (
	Reset = "\033[0m"
	Hint  = "\033[90m"
	Ok    = "\033[32m"
	Error = "\033[31m"
	Warn  = "\033[33m"
)

func (p *PodmanArg) ComposeCMD() (*model.Pod, error) {
	d := &model.Pod{
		UID:       p.UID,
		PodID:     filepath.Base(p.RemoteDir),
		PodName:   filepath.Base(p.RemoteDir),
		LocalDir:  p.LocalDir,
		RemoteDir: p.RemoteDir,
		Target:    p.Target,
		File:      p.File,
		Status:    "starting",
		Hostname:  p.Hostname,
		IP:        p.IP,
		Replicas:  1,
	}

	switch p.Command {
	case "up":
		return p.up(d)
	case "clear":
		return p.clear(d)
	case "down", "ps", "logs", "restart", "exec", "build":
		return p.runCMD(d)
	}
	return nil, fmt.Errorf("unsupported command: %s", p.Command)
}

func (p *PodmanArg) up(d *model.Pod) (*model.Pod, error) {
	fmt.Println("[+] create folder if not exist")
	if err := utils.SSHRun("mkdir", "-p", p.RemoteDir); err != nil {
		return nil, err
	}

	// * 同步檔案夾資料
	fmt.Println("[*] syncing files")
	fmt.Println(Hint + "──────────────────────────────────────────────────")
	if err := p.RsyncToRemote(); err != nil {
		return nil, err
	}
	fmt.Println("──────────────────────────────────────────────────" + Reset)

	// * 調整 docker-compose.yml 內容
	fmt.Println("[*] modifying compose file (remove ports)")
	if err := p.ModifyComposeFile(); err != nil {
		return nil, fmt.Errorf("[x] failed to modify compose file: %w", err)
	}

	// * 關閉舊的容器 (if exists)
	fmt.Println("[*] cleaning up old containers")
	_, _ = utils.SSEOutput(fmt.Sprintf(
		"cd '%s' && podman compose down -v >/dev/null 2>&1",
		p.RemoteDir,
	))
	removePod(d.UID)

	// * 執行動作
	fmt.Printf("[*] executing: podman compose %s\n", strings.Join(p.RemoteArgs, " "))
	fmt.Println(Hint + "──────────────────────────────────────────────────")
	remoteCmd := fmt.Sprintf("cd '%s' && podman compose %s 2>&1", p.RemoteDir, shellJoin(p.RemoteArgs))
	if !p.Detach {
		remoteCmd = fmt.Sprintf(`
				cleanup() {
					echo "[*] stopping containers"
					cd '%s' && podman compose down
				}
				trap cleanup INT TERM
				%s
			`, p.RemoteDir, remoteCmd)
	}
	if err := utils.SSHRun(remoteCmd); err != nil {
		return nil, err
	}
	fmt.Println(Hint + "──────────────────────────────────────────────────" + Reset)

	// * 取得 Pod 資訊
	projectName := filepath.Base(p.RemoteDir)
	podInfo, err := utils.SSEOutput(fmt.Sprintf(
		"podman pod ps --filter 'name=pod_%s' --format '{{.ID}}\t{{.Name}}'",
		projectName,
	))
	if err == nil && podInfo != "" {
		parts := strings.Split(strings.TrimSpace(podInfo), "\t")
		if len(parts) >= 2 {
			d.PodID = parts[0]
			d.PodName = parts[1]
		}
	}

	// * 輸出結果
	if p.Detach {
		fmt.Println("[*] service ports:")
		fmt.Println(Ok + "──────────────────────────────────────────────────")
		output, _ := utils.SSEOutput(fmt.Sprintf(
			"cd '%s' && podman ps --filter 'label=io.podman.compose.project=%s' --format 'table {{.Names}}\t{{.Ports}}'",
			p.RemoteDir,
			projectName),
		)
		fmt.Println(output)
		fmt.Printf("Pod ID: %s\n", d.PodID)
		fmt.Printf("Pod Name: %s\n", d.PodName)
		fmt.Printf("Hostname: %s\n", d.Hostname)
		fmt.Printf("IP: %s\n", d.IP)
		fmt.Println("──────────────────────────────────────────────────" + Reset)
	}

	// *  發送 Pod 資訊到 API
	fmt.Println("[*] syncing pod info to database")
	jsonData, err := json.Marshal(d)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal json: %w", err)
	}

	resp, err := http.Post(
		"http://localhost:8080/pod/upsert",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to upsert pod: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to upsert pod: status %d, body: %s", resp.StatusCode, string(body))
	}

	return d, nil
}

func (p *PodmanArg) clear(d *model.Pod) (*model.Pod, error) {
	// * 停止並移除容器和 volumes
	fmt.Println("[*] remove containers and volumes")
	fmt.Println(Hint + "──────────────────────────────────────────────────")
	downCmd := fmt.Sprintf(
		"cd '%s' && podman compose down -v 2>&1 | grep -v 'no container\\|no pod' || true",
		p.RemoteDir,
	)
	removePod(d.UID)
	if err := utils.SSHRun(downCmd); err != nil {
		return nil, fmt.Errorf("failed to remove containers: %w", err)
	}
	fmt.Println("──────────────────────────────────────────────────" + Reset)

	// * 移除映像
	fmt.Println("[*] clean images")
	fmt.Println(Hint + "──────────────────────────────────────────────────")
	imageCmd := fmt.Sprintf(
		"cd '%s' && podman compose down --rmi all 2>&1 | grep -v 'no container\\|no pod\\|no image' || true",
		p.RemoteDir,
	)
	if err := utils.SSHRun(imageCmd); err != nil {
		return nil, fmt.Errorf("failed to remove images: %w", err)
	}
	fmt.Println("──────────────────────────────────────────────────" + Reset)

	// * 移除資料夾
	fmt.Println("[*] remove project folder")
	fmt.Println(Hint + "──────────────────────────────────────────────────")
	removeCmd := fmt.Sprintf(
		"podman run --rm --privileged -v '%s:/parent' alpine:latest sh -c 'rm -rf /parent/%s'",
		filepath.Dir(p.RemoteDir),
		filepath.Base(p.RemoteDir),
	)
	if err := utils.SSHRun(removeCmd); err != nil {
		return nil, fmt.Errorf("failed to remove folder: %w", err)
	}
	fmt.Println(Hint + "──────────────────────────────────────────────────" + Reset)
	return d, nil
}

func (p *PodmanArg) runCMD(d *model.Pod) (*model.Pod, error) {
	fmt.Printf("[*] executing: podman compose %s\n", strings.Join(p.RemoteArgs, " "))
	fmt.Println(Hint + "──────────────────────────────────────────────────")
	if err := utils.SSHRun(fmt.Sprintf(
		"cd '%s' && podman compose %s",
		p.RemoteDir,
		shellJoin(p.RemoteArgs)),
	); err != nil {
		return nil, err
	}
	fmt.Println(Hint + "──────────────────────────────────────────────────" + Reset)

	if p.Command == "down" {
		removePod(d.UID)
	}
	return d, nil
}

func (p *PodmanArg) RsyncToRemote() error {
	env, err := utils.GetENV()
	if err != nil {
		return err
	}

	cmdArgs := []string{
		"-p", env.Password,
		"rsync",
		"-avz", "--delete",
		"--exclude=node_modules/", "--exclude=vendor/", "--exclude=__pycache__/",
		"--exclude=*.pyc", "--exclude=.venv/", "--exclude=venv/", "--exclude=env/",
		"--exclude=.env.local", "--exclude=.git/", "--exclude=.gitignore",
		"--exclude=*.log", "--exclude=.DS_Store", "--exclude=Thumbs.db",
		"-e", "ssh -o StrictHostKeyChecking=no",
		p.LocalDir + "/",
		fmt.Sprintf("%s:%s/", env.Remote, p.RemoteDir),
	}
	return utils.CMDRun("sshpass", cmdArgs...)
}

func (p *PodmanArg) ModifyComposeFile() error {
	composeFile := "docker-compose.yml"
	output, _ := utils.SSEOutput("test -f '%s/%s' || echo 'notfound'", p.RemoteDir, composeFile)
	if strings.TrimSpace(output) == "notfound" {
		composeFile = "docker-compose.yaml"
	}

	// 移除 ports
	sedCmds := []string{
		`sed -i -E 's/(["\x27]?)[0-9]+:([0-9]+)(["\x27]?)/\1\2\3/g' '%s/%s'`,
		`sed -i -E 's/(["\x27]?)\$\{[^}]+\}:([0-9]+)(["\x27]?)/\1\2\3/g' '%s/%s'`,
		`sed -i -E 's/\$\{[^}]+:[?][^}]+\}://g' '%s/%s'`,
	}

	for _, cmdTemplate := range sedCmds {
		cmd := fmt.Sprintf(cmdTemplate, p.RemoteDir, composeFile)
		if err := utils.SSHRun(cmd); err != nil {
			return err
		}
	}

	// 強制為所有相對路徑 volume 加入 :z（如果沒有）
	awkCmd := fmt.Sprintf(`
		awk '
		/^\s+- \.\/[^:]+:[^:]+$/ { print $0 ":z"; next }
		/^\s+- \.\/[^:]+:[^:]+:[^z]*$/ { gsub(/:([^z:]+)$/, ":\\1,z"); print; next }
		/^\s+- \.\/[^:]+:[^:]+:.*z/ { print; next }
		{ print }
		' '%s/%s' > '%s/%s.tmp' && mv '%s/%s.tmp' '%s/%s'
	`, p.RemoteDir, composeFile, p.RemoteDir, composeFile, p.RemoteDir, composeFile, p.RemoteDir, composeFile)

	return utils.SSHRun(awkCmd)
}

func shellJoin(args []string) string {
	escaped := make([]string, len(args))
	for i, arg := range args {
		if strings.ContainsAny(arg, " \t\n'\"\\$`") {
			escaped[i] = "'" + strings.ReplaceAll(arg, "'", "'\"'\"'") + "'"
		} else {
			escaped[i] = arg
		}
	}
	return strings.Join(escaped, " ")
}

func removePod(uid string) {
	jsonData, err := json.Marshal(&model.Pod{
		Dismiss: 1,
	})
	// * slience, if wrong, just wrong
	if err != nil {
		return
	}

	resp, err := http.Post(
		"http://localhost:8080/pod/update/"+uid,
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	// * slience, if wrong, just wrong
	if err != nil {
		return
	}
	defer resp.Body.Close()

	// * slience, if wrong, just wrong
	if resp.StatusCode != http.StatusOK {
		return
	}
}
