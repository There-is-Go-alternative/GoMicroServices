package database

import (
	account "github.com/There-is-Go-alternative/GoMicroServices/account/domain"
	chats "github.com/There-is-Go-alternative/GoMicroServices/chats/domain"
)

type MockDatabase struct {
	Chats    []chats.Chat
	Messages []chats.Message
}

func (d *MockDatabase) CreateChat(chat chats.Chat) (chats.Chat, error) {
	d.Chats = append(d.Chats, chat)
	return chat, nil
}

func (d *MockDatabase) CreateMsg(message chats.Message) (chats.Message, error) {
	d.Messages = append(d.Messages, message)
	return message, nil
}

func (d MockDatabase) GetAllChatsOfUser(user_ID account.AccountID) ([]chats.Chat, error) {
	result := []chats.Chat{}
	for _, current_chat := range d.Chats {
		for _, current_user_id := range current_chat.UsersIDs {
			if current_user_id == user_ID {
				result = append(result, current_chat)
			}
		}
	}
	return result, nil
}
