package main

import (
	"context"
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
		Status: todov1.TodoItem_STATUS_OPEN,
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
	res := connect.NewResponse(&todov1.UpdateTodoResponse{
		
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
	Toder := &TodoServer{}
	mux := http.NewServeMux()
	path, handler := todov1connect.NewTodoServiceHandler(Toder)
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:8080",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}