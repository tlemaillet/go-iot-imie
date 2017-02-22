package common

type SensorData struct {
	SensorName string `json:"sensor"`
	Measurement string `json:"type"`
	Time int64 `json:"receivedAt"`
	Value string `json:"value"`
}


type Sensor struct {
	DisplayName string `json:"displayname"`
	Vendor string `json:"vendor"`
	Product string `json:"product"`
	Version int `json:"version"`
}