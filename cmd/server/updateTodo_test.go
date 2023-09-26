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


	reqItem := &todov1.UpdateTodoRequest{
		Id:     testUuid,
		Status: todov1.TodoItem_STATUS_NOSTARTED,
	}
	
	reqTodo := &connect.Request[todov1.UpdateTodoRequest]{
		Msg: reqItem,
	}

	expected := todov1.TodoItem_STATUS_COMPLETED
	actualResult, err := Todoer.UpdateTodo(context.TODO(), reqTodo)

	if err != nil {
		t.Fatalf("UpdateTodo failed with error: %v", err)
	}
	if actualResult.Msg.Item.Status != expected {
		t.Errorf("actualResult [%v] is not equal to expected [%v]", actualResult.Msg.Item.Title, expected)
	}
}