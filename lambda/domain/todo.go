package domain

type Todo struct {
	UserId string `json:"-" dynamodbav:"user_id"`
	Id     int    `json:"id" dynamodbav:"id"`
	Title  string `json:"title" dynamodbav:"title"`
	Done   bool   `json:"done" dynamodbav:"done"`
}
