package network

import "encoding/json"

func Encode(v any) ([]byte, error) {
	return json.Marshal(v)
}

func Decode(data []byte, v any) error {
	return json.Unmarshal(data, v)
}

func DecodePayload[T any](env Envelope) (T, error) {
	var out T
	err := json.Unmarshal(env.Payload, &out)
	return out, err
}
