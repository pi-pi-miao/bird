package bird_controller


const (
	Success                   = "200"
	ReadRequestError          = "10001"
	UnmarshalRequestBodyError = "10002"
	PasswordError             = "10003"
)

type BirdResponse struct {
	Code string  `json:"code"`
	Data string  `json:"data"`
}
