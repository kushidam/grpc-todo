package main

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	todov1 "github.com/kushidam/grpc-todo/gen/todo/v1"
)

func TestDeleteTodo(t *testing.T) {
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
		request        *connect.Request[todov1.DeleteTodoRequest]
		expectResId	   string
		expectError    bool
	}{
		{"Delete completed successfully", &connect.Request[todov1.DeleteTodoRequest]{ Msg: &todov1.DeleteTodoRequest{ Id: testUuid, }, }, testUuid, false, },
		{"Todo item not found",           &connect.Request[todov1.DeleteTodoRequest]{ Msg: &todov1.DeleteTodoRequest{ Id: "123456", }, },"NoData",  true, },
		// シナリオのテストケースをここに追加
	}
	
	// テストケースをループで実行
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			actualResult, err := Todoer.DeleteTodo(context.TODO(), test.request)

			// エラーが期待通りの場合
			if test.expectError && err == nil {
				t.Errorf("Expected an error, but got nil")
			}

			// エラーが期待通りでない場合
			if !test.expectError && err != nil {
				t.Errorf("Expected no error, DeleteTodo failed with error: %v", err)
			}

			// ステータスが期待通りでない場合
			if actualResult != nil && actualResult.Msg.Id != test.expectResId {
				t.Errorf("actualResult [%v] is not equal to expected [%v]", actualResult.Msg.Id, test.expectResId)
			}
		})
	}
}