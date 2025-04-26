package client

import (
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

type SensitiveData struct {
	SensitiveFields []string `json:"sensitive_fields"`
}

var (
	SensitiveDataInstance  *SensitiveData
	defaultSensitiveFields = []string{
		"data.phoneNumber",
		"data.email",
	}
)

func NewSensitiveData() *SensitiveData {
	return &SensitiveData{}
}

func init() {
	SensitiveDataInstance = NewSensitiveData()
	SensitiveDataInstance.SetSensitiveFields(defaultSensitiveFields)
}

func (s *SensitiveData) IsSensitiveField(field string) bool {
	for _, sensitiveField := range s.SensitiveFields {
		if sensitiveField == field {
			return true
		}
	}
	return false
}

func (s *SensitiveData) AddSensitiveField(field string) {
	s.SensitiveFields = append(s.SensitiveFields, field)
}

func (s *SensitiveData) RemoveSensitiveField(field string) {
	for i, sensitiveField := range s.SensitiveFields {
		if sensitiveField == field {
			s.SensitiveFields = append(s.SensitiveFields[:i], s.SensitiveFields[i+1:]...)
			return
		}
	}
}

func (s *SensitiveData) GetSensitiveFields() []string {
	return s.SensitiveFields
}

func (s *SensitiveData) SetSensitiveFields(fields []string) {
	s.SensitiveFields = fields
}

func (s *SensitiveData) ClearSensitiveFields() {
	s.SensitiveFields = []string{}
}

func (s *SensitiveData) FilterSensitiveFields(data string) string {
	for _, sensitiveField := range s.SensitiveFields {
		value := gjson.Get(data, sensitiveField).Value()
		if value == nil {
			continue
		}

		data, _ = sjson.Set(data, sensitiveField, "********")
	}
	return data
}
