package main

import (
	"fmt"
	"github.com/go-chi/chi"
	_ "web/docs"

	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log/slog"
	"os"
	"web/internal/config"
	"web/internal/logging"
	"web/internal/server/server"
	"web/internal/server/server/handlers"
	"web/internal/storage"
	"web/internal/storage/sqlite"
)

// Todo get by tag + due date
// Todo delete h.AllTags param and make it in struct or other ideas
// Todo delete task after 2 days of due date
// Todo work from due date format make it more simple

// @title Swagger Todo App Application
// @version 1.0
// @description API Server for Todo Application
// @host localhost:8000
// @BasePath /

func main() {
	// Setup logger
	log := logging.SetupLogger()

	// Read config file and create new Config{}
	cfg := config.NewConfig("config/config.yaml", log)

	// Create Sql connection
	SqlDataBase := SqlConnect(cfg, log)

	// Create new server
	httpServer := server.NewServer(&SqlDataBase, log)
	// Init handlers
	allHandlers := handlers.NewHandlers(httpServer)
	httpServer.InitHandlers(allHandlers)
	// Init routes
	initRoutes(httpServer)
	// Start server
	httpServer.Start(cfg, log)

}

// initRoutes init routes for server
func initRoutes(server *server.Server) {
	router := server.Router
	router.Route("/task", func(r chi.Router) {
		// get all tasks
		r.Get("/", server.Handlers.GetTasksHandler)
		// get task by id
		r.Get("/{id:\\d*}", server.Handlers.GetTaskHandler)
		// get all tasks by tag or tags list with different modes.
		// full - returns tasks who have  includes specified tag or tags list in the task tags
		// short - returns tasks who have only specified tag or tags list included in the task tags
		r.Route("/{mode:(?:short|full)}", func(r chi.Router) {
			// for tag in URL
			//r.Get("/{tag}", server.Handlers.GetTasksByTagHandler)
			// for tag in query
			r.Get("/", server.Handlers.GetTasksByTagHandler)
		}) //r.Get("/{mode:(?:short|full)}/{tag:[A-Za-z0-9?]+}", server.Handlers.GetTasksByTagHandler)

		// get tasks by due date
		r.Get("/{due:[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}(?::|%3A)[0-9]{2}(?::|%3A)[0-9]{2}Z}", server.Handlers.GetTasksByDueDateHandler)

		// create new task using request body data
		// request body example:
		// {"text": "text", "tags": ["tag", "tag", "tag"], "due": "2021-01-01 00:00:00"}
		r.Post("/", server.Handlers.CreateTaskHandler)

		// delete task by id
		r.Delete("/{id:[0-9]*}", server.Handlers.DeleteTaskHandler)
		// delete all tasks
		r.Delete("/", server.Handlers.DeleteTasksHandler)
	})
	router.Route("/tag", func(r chi.Router) {
		// get all tags
		r.Get("/", server.Handlers.GetTagsHandler)
		// get tag by name
		r.Get("/{name:[A-Za-z]+}", server.Handlers.GetTagHandler)

		// create new tag using request body data
		// request body example:
		// {"name": "name"}
		r.Post("/", server.Handlers.CreateTagHandler)

		// delete all tags
		r.Delete("/", server.Handlers.DeleteTagsHandler)
		// delete tag by name
		r.Delete("/{name:[A-Za-z]+}", server.Handlers.DeleteTagHandler)
	})
	router.MethodNotAllowed(server.Handlers.MethodNotAllowedHandler)
	router.NotFound(server.Handlers.NotFoundHandler)
	router.Get("/swagger/*", httpSwagger.Handler())

}

// SqlConnect connect to database from config
func SqlConnect(cfg *config.Config, log *slog.Logger) storage.Storage {
	var sqlStorage storage.Storage

	switch cfg.Type {
	case "sqlite":
		sqlStorage = &sqlite.StoreSqlite{}
	default:
		log.Error(fmt.Sprintf("Unknown database type: %v", cfg.Type))
		os.Exit(1)
	}

	db := sqlStorage.Connect(cfg, log)

	log.Info(fmt.Sprintf("Successfully connected to database: %v", cfg.Type))
	return db
}
