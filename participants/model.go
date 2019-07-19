package participants

type Participant struct {
	Id    uint32 `json:"id"`
	Name  string `json:"name"`
	Score uint16 `json:"score"`
}
