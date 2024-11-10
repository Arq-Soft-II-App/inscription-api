// main.go
package main

import (
	"inscription-api/src/config/builder"
	"inscription-api/src/config/envs"

	"go.uber.org/zap"
)

func main() {
	env := envs.LoadEnvs(".env")
	port := env.Get("PORT")
	if port == "" {
		port = "8000"
	}

	app := builder.BuildApp()
	defer app.DisconnectDB()

	router := app.GetRouter()
	app.Logger.Info("Iniciando servidor", zap.String("port", port))
	if err := router.Run(":" + port); err != nil {
		app.Logger.Fatal("Error al iniciar el servidor", zap.Error(err))
	}
}
