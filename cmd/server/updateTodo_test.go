package main

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	todov1 "github.com/kushidam/grpc-todo/gen/todo/v1"
)

func TestUpdateTodo(t *testing.T) {
	Todoer := &TodoServer{}
	testUuid := uuid.New().String()

	testItem := &todov1.TodoItem{
		Id:     testUuid,
		Title:  "Test Title",
		Status: todov1.TodoItem_STATUS_NOSTARTED,
	}
	Todoer.items.Store(testUuid, testItem)

	tests := []struct {
		name           string
		request        *connect.Request[todov1.UpdateTodoRequest]
		requestStatus  todov1.TodoItem_Status
		expectedStatus todov1.TodoItem_Status
		expectError    bool
	}{
		{"Update to COMPLETED", &connect.Request[todov1.UpdateTodoRequest]{Msg: &todov1.UpdateTodoRequest{Id: testUuid,},}, todov1.TodoItem_STATUS_COMPLETED, todov1.TodoItem_STATUS_COMPLETED, false,},
		{"Update to NOSTARTED", &connect.Request[todov1.UpdateTodoRequest]{Msg: &todov1.UpdateTodoRequest{Id: testUuid,},}, todov1.TodoItem_STATUS_NOSTARTED, todov1.TodoItem_STATUS_NOSTARTED, false,},
		{"Todo item not found", &connect.Request[todov1.UpdateTodoRequest]{Msg: &todov1.UpdateTodoRequest{Id: "123456",},}, todov1.TodoItem_STATUS_COMPLETED, todov1.TodoItem_STATUS_COMPLETED, true,},
		// シナリオのテストケースをここに追加
	}

	// テストケースをループで実行
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			actualResult, err := Todoer.UpdateTodo(context.TODO(), test.request)

			// エラーが期待通りの場合
			if test.expectError && err == nil {
				t.Errorf("Expected an error, but got nil")
			}

			// エラーが期待通りでない場合
			if !test.expectError && err != nil {
				t.Errorf("Expected no error, UpdateTodo failed with error: %v", err)
			}

			// ステータスが期待通りでない場合
			if actualResult != nil && actualResult.Msg.Item.Status != test.expectedStatus {
				t.Errorf("actualResult [%v] is not equal to expected [%v]", actualResult.Msg.Item.Status, test.expectedStatus)
			}
		})
	}
}