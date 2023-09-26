package main

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/google/uuid"
	todov1 "github.com/kushidam/grpc-todo/gen/todo/v1"
	"github.com/kushidam/grpc-todo/gen/todo/v1/todov1connect"
)

type TodoServer struct {
	items sync.Map
}

func (s *TodoServer) CreateTodo(
	ctx context.Context,
	req *connect.Request[todov1.CreateTodoRequest],
) (*connect.Response[todov1.CreateTodoResponse], error) {
	id := uuid.New().String()
	item := &todov1.TodoItem{
		Id:     id,
		Title:  req.Msg.Title,
		Status: todov1.TodoItem_STATUS_NOSTARTED,
	}
	s.items.Store(id, item)
	res := connect.NewResponse(&todov1.CreateTodoResponse{
		Item: item,
	})
	return res, nil
}

func (s *TodoServer) UpdateTodo(
	ctx context.Context,
	req *connect.Request[todov1.UpdateTodoRequest],
) (*connect.Response[todov1.UpdateTodoResponse], error) {
	reqKeyId := req.Msg.Id;
	todoItemBefore, _ := s.items.Load(reqKeyId)
	if todoItemBefore == nil {
		return nil, errors.New("Todo item not found") //TODO非機能要件
	}

	todoItemAfter, ok := todoItemBefore.(*todov1.TodoItem)
	if !ok {
		return nil, errors.New("Failed to assert TodoItem type") //TODO非機能要件
	}

	if todoItemAfter.Status == todov1.TodoItem_STATUS_NOSTARTED {
		todoItemAfter.Status = todov1.TodoItem_STATUS_COMPLETED
	} else if todoItemAfter.Status == todov1.TodoItem_STATUS_COMPLETED {
		todoItemAfter.Status = todov1.TodoItem_STATUS_NOSTARTED
	} else {
		todoItemAfter.Status = todov1.TodoItem_STATUS_UNKNOWN_UNSPECIFIED
	}

	if !s.items.CompareAndSwap(reqKeyId, todoItemBefore, todoItemAfter) {
		return nil, errors.New("Failed to CompareAndSwap") //TODO非機能要件
	}

	res := connect.NewResponse(&todov1.UpdateTodoResponse{
		Item: todoItemAfter, 
	})
	return res, nil
}

func (s *TodoServer) DeleteTodo(
	ctx context.Context,
	req *connect.Request[todov1.DeleteTodoRequest],
) (*connect.Response[todov1.DeleteTodoResponse], error) {
	res := connect.NewResponse(&todov1.DeleteTodoResponse{
		
	})
	return res, nil
}

func main() {
	Todoer := &TodoServer{}
	mux := http.NewServeMux()
	path, handler := todov1connect.NewTodoServiceHandler(Todoer)
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}