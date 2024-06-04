package library

type ErrorPayload struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type CreateSuggestionParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type StarSuggestionParams struct {
	Star int `json:"star"`
}

type WithReasonParams struct {
	Reason string `json:"reason"`
}
