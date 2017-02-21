package common

type SensorData struct {
	SensorName string `json:"sensor"`
	Measurement string `json:"type"`
	Time int64 `json:"receivedAt"`
	Value string `json:"value"`
}