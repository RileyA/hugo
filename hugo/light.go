package hugo

type LightState struct {
	On        bool      `json:"on"`
	Bri       uint8     `json:"bri"`
	Hue       uint16    `json:"hue"`
	Sat       uint8     `json:"sat"`
	Effect    string    `json:"effect"`
	Xy        []float64 `json:"xy"`
	Ct        uint16    `json:"ct"`
	Alert     string    `json:"alert"`
	ColorMode string    `json:"colormode"`
	Reachable bool      `json:"reachable"`
}

type Light struct {
	State            LightState `json:"state"`
	Type             string     `json:"type"`
	Name             string     `json:"name"`
	ModelId          string     `json:"modelid"`
	ManufacturerName string     `json:"manufacturername"`
	UniqueId         string     `json:"uniqueid"`
	SwVersion        string     `json:"swversion"`
}
