// config/builder/builder.go
package builder

import (
	"inscription-api/src/clients"
	"inscription-api/src/config/db"
	"inscription-api/src/config/log"
	"inscription-api/src/controllers"
	"inscription-api/src/middlewares"
	"inscription-api/src/routes"
	"inscription-api/src/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type AppBuilder struct {
	db                    *gorm.DB
	Logger                *zap.Logger
	inscriptionRepo       clients.InscriptionRepository
	inscriptionService    services.InscriptionService
	inscriptionController *controllers.InscriptionController
	routes                *gin.Engine
}

func NewAppBuilder() *AppBuilder {
	return &AppBuilder{}
}

func BuildApp() *AppBuilder {
	return NewAppBuilder().
		BuildLogger().
		BuildDBConnection().
		BuildInscriptionRepo().
		BuildInscriptionService().
		BuildInscriptionController().
		BuildRouter()
}

func (b *AppBuilder) BuildLogger() *AppBuilder {
	b.Logger = log.GetLogger()
	b.Logger.Info("[INSCRIPTION-API] Logger inicializado")
	return b
}

func (b *AppBuilder) BuildDBConnection() *AppBuilder {
	var err error
	b.db, err = db.ConnectDB(b.Logger)
	if err != nil {
		b.Logger.Fatal("[INSCRIPTION-API] Error al conectar a la base de datos", zap.Error(err))
	}
	b.Logger.Info("[INSCRIPTION-API] Conexión a la base de datos establecida")
	return b
}

func (b *AppBuilder) DisconnectDB() {
	if b.db != nil {
		sqlDB, err := b.db.DB()
		if err != nil {
			b.Logger.Error("[INSCRIPTION-API] Error al obtener la conexión SQL", zap.Error(err))
		} else {
			if err := sqlDB.Close(); err != nil {
				b.Logger.Error("[INSCRIPTION-API] Error al desconectar de la base de datos", zap.Error(err))
			} else {
				b.Logger.Info("[INSCRIPTION-API] Conexión a la base de datos cerrada")
			}
		}
	}
	_ = b.Logger.Sync()
}

func (b *AppBuilder) BuildInscriptionRepo() *AppBuilder {
	b.inscriptionRepo = clients.NewGormInscriptionRepository(b.db, b.Logger)
	b.Logger.Info("[INSCRIPTION-API] Repositorio de inscripciones inicializado")
	return b
}

func (b *AppBuilder) BuildInscriptionService() *AppBuilder {
	b.inscriptionService = services.NewInscriptionService(b.inscriptionRepo, b.Logger)
	b.Logger.Info("[INSCRIPTION-API] Servicio de inscripciones inicializado")
	return b
}

func (b *AppBuilder) BuildInscriptionController() *AppBuilder {
	b.inscriptionController = controllers.NewInscriptionController(b.inscriptionService, b.Logger)
	b.Logger.Info("[INSCRIPTION-API] Controlador de inscripciones inicializado")
	return b
}

func (b *AppBuilder) BuildRouter() *AppBuilder {
	b.routes = gin.New()
	b.routes.Use(gin.Recovery())
	b.routes.Use(middlewares.LoggerMiddleware(b.Logger))
	b.routes.Use(middlewares.APIKeyAuthMiddleware(b.Logger))
	routes.SetupRoutes(b.routes, b.inscriptionController)
	b.Logger.Info("[INSCRIPTION-API] Router inicializado")
	return b
}

func (b *AppBuilder) GetRouter() *gin.Engine {
	return b.routes
}
