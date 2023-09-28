package main

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"connectrpc.com/connect"
	"golang.org/x/exp/slog"
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
	slog.Info("CreateTodo", req )
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
	slog.Info("UpdateTodo", req )
    reqKeyId := req.Msg.Id
    todoItemBefore, _ := s.items.Load(reqKeyId)
    if todoItemBefore == nil {
        errMsg := "Todo item not found"
        slog.Warn(errMsg, "request_id", reqKeyId)
        return nil, errors.New(errMsg)
    }

    todoItemAfter, ok := todoItemBefore.(*todov1.TodoItem)
    if !ok {
		errMsg := "Failed to assert TodoItem type"
        slog.Warn(errMsg, "request_id", reqKeyId)
        return nil, errors.New(errMsg)
    }

    if todoItemAfter.Status == todov1.TodoItem_STATUS_NOSTARTED {
        todoItemAfter.Status = todov1.TodoItem_STATUS_COMPLETED
    } else if todoItemAfter.Status == todov1.TodoItem_STATUS_COMPLETED {
        todoItemAfter.Status = todov1.TodoItem_STATUS_NOSTARTED
    } else {
        todoItemAfter.Status = todov1.TodoItem_STATUS_UNKNOWN_UNSPECIFIED
    }

    if !s.items.CompareAndSwap(reqKeyId, todoItemBefore, todoItemAfter) {
		errMsg := "Failed to CompareAndSwap"
        slog.Warn(errMsg, "todoBefore", todoItemBefore, "todoAfter", todoItemAfter)
        return nil, errors.New(errMsg)
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
	slog.Info("DeleteTodo", req )
    reqKeyId := req.Msg.Id
    _, ok := s.items.LoadAndDelete(reqKeyId)
    if !ok {
		errMsg := "Todo item not found"
        slog.Warn(errMsg, "request_id", reqKeyId)
        return nil, errors.New(errMsg)
    }

    res := connect.NewResponse(&todov1.DeleteTodoResponse{
        Id: reqKeyId,
    })
    return res, nil
}

func main() {
	Todoer := &TodoServer{}
	mux := http.NewServeMux()
	path, handler := todov1connect.NewTodoServiceHandler(Todoer)
	mux.Handle(path, handler)
	http.ListenAndServe(
		"localhost:50051",
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	)
}