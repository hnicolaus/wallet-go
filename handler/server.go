package handler

import (
	"github.com/WalletService/usecase"
)

type Server struct {
	Usecase usecase.UsecaseInterface
}

type NewServerOptions struct {
}

func NewServer(usecase usecase.UsecaseInterface) *Server {
	return &Server{
		Usecase: usecase,
	}
}
