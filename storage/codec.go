package storage

import (
	"encoding/json"
	"gradebook/domain"
)

func EncodeGradebook(record domain.Gradebook) ([]byte, error) { return json.Marshal(record) }
func DecodeGradebook(data []byte) (domain.Gradebook, error) {
	var record domain.Gradebook
	err := json.Unmarshal(data, &record)
	return record, err
}
func EncodeAudit(event domain.AuditEvent) ([]byte, error) { return json.Marshal(event) }
func DecodeAudit(data []byte) (domain.AuditEvent, error) {
	var event domain.AuditEvent
	err := json.Unmarshal(data, &event)
	return event, err
}
