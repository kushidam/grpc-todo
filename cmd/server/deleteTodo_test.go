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


	reqItem := &todov1.DeleteTodoRequest{
		Id:     testUuid,
	}
	
	reqTodo := &connect.Request[todov1.DeleteTodoRequest]{
		Msg: reqItem,
	}

	expected := testUuid
	actualResult, err := Todoer.DeleteTodo(context.TODO(), reqTodo)

	if err != nil {
		t.Fatalf("DeleteTodo failed with error: %v", err)
	}
	if actualResult.Msg.Id != expected {
		t.Errorf("actualResult [%v] is not equal to expected [%v]", actualResult.Msg.Id, expected)
	}
}