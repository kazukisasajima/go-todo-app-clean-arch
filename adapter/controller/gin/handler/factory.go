package handler

import "sync"

var (
	serverHandler *ServerHandler
	once 		sync.Once
)

type ServerHandler struct {
	*TaskHandler
	*UserHandler
}

func NewHandler() *ServerHandler {
	once.Do(func() {
		serverHandler = &ServerHandler{}
	})
	return serverHandler
}

func (h *ServerHandler) Register(i interface{}) *ServerHandler {
	switch interfaceType := i.(type) {
	case *TaskHandler:
		serverHandler.TaskHandler = interfaceType
	case *UserHandler:
		serverHandler.UserHandler = interfaceType
	}
	return serverHandler
}
