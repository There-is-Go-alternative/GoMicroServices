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

func (d MockDatabase) GetAllMessagesOfChat(chat_ID chats.ChatID) ([]chats.Message, error) {
	result := []chats.Message{}
	for _, current_message := range d.Messages {
		if current_message.ChatID == chat_ID {
			result = append(result, current_message)
		}
	}
	return result, nil
}

func (d MockDatabase) GetMessage(message_id string) (chats.Message, error) {
	for _, current_message := range d.Messages {
		if current_message.ID == string(message_id) {
			return current_message
		}
	}
	return nil, nil
}

func (d MockDatabase) GetChat(chat_id string) (chats.Chat, error) {
	for _, current_chat := range d.Chats {
		if current_chat.ID == string(chat_id) {
			return current_chat
		}
	}
	return nil, nil
}

func (d *MockDatabase) DeleteChat(chat_id string) bool {
	for index, current_chat := range d.Chats {
		if current_chat.ID == string(chat_id) {
			d.Chats[index] = d.Chats[len(d.Chats)-1]
			d.Chats = d.Chats[:len(d.Chats)-1]
			return true
		}
	}
	return false
}

func (d *MockDatabase) DeleteMessage(message_id string) bool {
	for index, current_message := range d.Messages {
		if current_message.ID == string(message_id) {
			d.Messages[index] = d.Messages[len(d.Messages)-1]
			d.Messages = d.Messages[:len(d.Messages)-1]
			return true
		}
	}
	return false
}
