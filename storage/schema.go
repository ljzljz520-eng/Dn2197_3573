package storage

import "go.etcd.io/bbolt"

func EnsureSchema(db *bbolt.DB) error {
	return db.Update(func(tx *bbolt.Tx) error {
		for _, name := range [][]byte{gradesBucket, auditBucket, metaBucket} {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}

func BucketNames() []string { return []string{"gradebooks", "audits", "metadata"} }

func IsSchemaReady(names []string) bool {
	if len(names) != 3 {
		return false
	}
	seen := map[string]bool{}
	for _, name := range names {
		seen[name] = true
	}
	return seen["gradebooks"] && seen["audits"] && seen["metadata"]
}
