package server

import (
	"fmt"
	"net/http"

	"github.com/kigorek/go_final_project/pkg/api"
)

func Run(port string) error {
	api.Init()
	fmt.Println("Start HTTP Server")

	http.Handle("/", http.FileServer(http.Dir("./web")))

	err := http.ListenAndServe(port, nil)

	if err != nil {
		return fmt.Errorf("Ошибка при запуске сервера: %w", err)
	}

	return nil
}
