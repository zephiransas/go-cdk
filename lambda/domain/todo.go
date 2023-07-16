package domain

type Todo struct {
	Id    string `json:"id" dynamodbav:"id"`
	Title string `json:"title" dynamodbav:"title"`
	Done  bool   `json:"done" dynamodbav:"done"`
}
