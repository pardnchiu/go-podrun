package model

import "time"

type Pod struct {
	ID        int64     `json:"id"`
	UID       string    `json:"uid"`
	PodID     string    `json:"pod_id"`
	PodName   string    `json:"pod_name"`
	LocalDir  string    `json:"local_dir"`
	RemoteDir string    `json:"remote_dir"`
	File      string    `json:"file"`
	Target    string    `json:"target"`
	Status    string    `json:"status"`
	Hostname  string    `json:"hostname"`
	IP        string    `json:"ip"`
	Replicas  int       `json:"replicas"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Dismiss   int       `json:"dismiss"`
}
