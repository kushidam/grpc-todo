package main

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	todov1 "github.com/kushidam/grpc-todo/gen/todo/v1"
)

func TestCreateTodo(t *testing.T) {
	
	tests := []struct {
		name           string
		request        *connect.Request[todov1.CreateTodoRequest]
		expectedTitle  string
	}{
		{
			name: "Create Todo with valid input",
			request: &connect.Request[todov1.CreateTodoRequest]{
				Msg: &todov1.CreateTodoRequest{
					Title: "test1 title",
				},
			},
			expectedTitle: "test1 title",
		},
		{
			name: "Create Todo with empty string",
			request: &connect.Request[todov1.CreateTodoRequest]{
				Msg: &todov1.CreateTodoRequest{
					Title: "",
				},
			},
			expectedTitle: "",
		},
		// シナリオのテストケースをここに追加
	}

	Todoer := &TodoServer{}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actualResult, err := Todoer.CreateTodo(context.TODO(), test.request)

			if err != nil {
				t.Fatalf("CreateTodo failed with error: %v", err)
			}
			if actualResult.Msg.Item.Title != test.expectedTitle {
				t.Errorf("actualResult [%v] is not equal to expected [%v]", actualResult.Msg.Item.Title, test.expectedTitle)
			}
		})
	}
}