package database

import "context"

func (s *SQLite) ListPodRecords(ctx context.Context, uid string) ([]ContainerRecord, error) {
	rows, err := s.db.QueryContext(ctx, `
  SELECT
    pods.local_dir,
    pods.remote_dir,
    pods.file,
    pods.hostname,
    pods.ip,
    records.content
  FROM records
  LEFT JOIN pods ON records.pod_id = pods.id
  WHERE pods.dismiss = 0 AND pods.uid = ?
  `, uid)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []ContainerRecord
	for rows.Next() {
		var r ContainerRecord
		if err := rows.Scan(&r.LocalDir, &r.RemoteDir, &r.File,
			&r.Hostname, &r.IP, &r.Content); err != nil {
			return nil, err
		}
		records = append(records, r)
	}

	return records, rows.Err()
}
