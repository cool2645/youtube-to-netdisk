package carrier

type InfoDict struct {
	ID                 string               `json:"id"`
	Creator            string               `json:"creator"`
	Uploader           string               `json:"uploader"`
	Description        string               `json:"description"`
	Title              string               `json:"title"`
	RequestedSubtitles map[string]Subtitles `json:"requested_subtitles"`
}

type Subtitles struct {
	Extension string `json:"ext"`
	URL       string `json:"url"`
}
