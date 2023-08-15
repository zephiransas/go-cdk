package domain

import "app/domain/vo"

type Todo struct {
	UserId vo.SubId `json:"-" dynamodbav:"user_id"`
	Id     int      `json:"id" dynamodbav:"id"`
	Title  string   `json:"title" dynamodbav:"title"`
	Done   bool     `json:"done" dynamodbav:"done"`
}
