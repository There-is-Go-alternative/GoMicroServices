package database

import (
	chats "github.com/There-is-Go-alternative/GoMicroServices/chats/domain"
)

type MockDatabase struct {
	Chats []chats.Chat
	Msgs  []chats.Message
}

func (d *MockDatabase) CreateChat(chat chats.Chat) (chats.Chat, error) {
	d.Chats = append(d.Chats, chat)
	return chat, nil
}

func (d *MockDatabase) CreateMsg(message chats.Message) (chats.Message, error) {
	d.Msgs = append(d.Msgs, message)
	return message, nil
}

func (d MockDatabase) GetAllChatsOfUser(user_ID string) ([]chats.Chat, error) {
	result := []chats.Chat{}
	for _, current_chat := range d.Chats {
		for _, elem := range current_chat.UsersIDs {
			if elem == user_ID {
				result = append(result, current_chat)
			}
		}
	}
	return result, nil
}
