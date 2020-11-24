package MyRPC

import (
	"encoding/json"
)

//Marshaller
type Marshaller struct{}

func (m *Marshaller) Marshal(obj interface{}) ([]byte, error) {
	return json.Marshal(obj)
}
func (m *Marshaller) Unmarshal(data []byte, pointer interface{}) error {
	return json.Unmarshal(data, pointer)
}
