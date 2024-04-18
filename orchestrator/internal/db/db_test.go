package db

import "testing"

func TestInsertUser(t *testing.T) {
	//Создаем нового юзера
	var user UserJson
	user.Login = "test"
	user.Password = "12345678"
	err := InsertUser(user)
	if err != nil {
		t.Fatal("Unexpected error: %v", err)
		return
	}
	//Изменяем пароль и смотрим как отреагирует функция
	user.Password = "87654321"
	err = ValidateUser(user)
	//Если не ошибка то функция не работает
	if err == nil {
		t.Fatal("Expected an error")
		return
	}
	//Все работает
}
