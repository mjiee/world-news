package main

import (
	"context"
	"fmt"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type GreetRequest struct {
	Message string `json:"msg"`
}

type GreetResponse struct {
	Message string `json:"msg"`
}

// Greet returns a greeting for the given name
func (a *App) Greet(req *GreetRequest) *GreetResponse {
	return &GreetResponse{
		Message: fmt.Sprintf("Hello %s, It's show time!", req.Message),
	}
}
