package hub

import "encoding/json"

type VideoData struct {
	Url     string  `json:"url"`
	Time    float32 `json:"time"`
	Playing bool    `json:"playing"`
}

// Update will be used to abstract some side-effects when connecting to a memory-database
// for now keep it as is
func (v *VideoData) Update(d VideoData) {
	v.Playing = d.Playing
	v.Time = d.Time
	v.Url = d.Url
}

func unmarshalVideo(marshalled interface{}) VideoData {
	res := VideoData{}
	bytes, _ := marshalled.([]byte)

	err := json.Unmarshal(bytes, &res)
	if err != nil {
		panic(err)
	}

	return res
}
