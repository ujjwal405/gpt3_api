package model

type User_Query struct {
	Query string `json:"query"`
}
type Query_Response struct {
	Answer string `json:"answer"`
}
type User_Search struct {
	Document string `json:"document"`
	Query    string `json:"query"`
}
type Search_Response struct {
	Answer string  `json:"answer"`
	Score  float64 `json:"score"`
}
