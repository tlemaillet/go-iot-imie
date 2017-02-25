package common

type SensorData struct {
	SensorName string `json:"sensor"`
	Measurement string `json:"type"`
	Time int64 `json:"receivedAt"`
	Value string `json:"value"`
}

type Sensor struct {
	Id int `json:"id,omitempty"`
	DisplayName string `json:"displayname"`
	Uuid string `json:"uuid,omitempty"`
	Vendor string `json:"vendor"`
	Product string `json:"product"`
	Version int `json:"version"`
	Enable bool `json:"enable,omitempty"`
}