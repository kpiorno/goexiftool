package dtos

type Data map[string]interface{}

//Response defines the response from exiftool
type Response struct {
	Tags []*Data `json:"tags"`
}
