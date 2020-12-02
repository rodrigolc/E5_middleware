package MyRPC

import (
	"encoding/json"
)

//Marshaller
type Marshaller struct{}

func (m *Marshaller) Marshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func (m *Marshaller) Unmarshal(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
