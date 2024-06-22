package models

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

const UserRole string = "user"
const AssistantRole string = "assistant"
