package main

import (
	"context"
	"testing"

	"connectrpc.com/connect"
	todov1 "github.com/kushidam/grpc-todo/gen/todo/v1"
)

func TestCreateTodo(t *testing.T) {
	Todoer := &TodoServer{}
	// testUuid := uuid.New().String()
	
	reqItem := &todov1.CreateTodoRequest{
		Title:   "test1 title",
	}
	reqTodo := &connect.Request[todov1.CreateTodoRequest]{
		Msg: reqItem,
	}

	expected := "test1 title"
	actualResult, err := Todoer.CreateTodo(context.TODO(), reqTodo)

	if err != nil {
		t.Fatalf("CreateTodo failed with error: %v", err)
	}
	if actualResult.Msg.Item.Title != expected {
		t.Errorf("actualResult [%v] is not equal to expected [%v]", actualResult.Msg.Item.Title, expected)
	}
}