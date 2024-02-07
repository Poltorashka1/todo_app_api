package main

import (
	"fmt"
	"github.com/go-chi/chi"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"log/slog"
	"os"
	_ "web/docs"
	"web/internal/config"
	"web/internal/logging"
	"web/internal/server/middleware"
	"web/internal/server/server"
	"web/internal/storage"
	"web/internal/storage/sqlite"
)

// Todo Exclude tag from query

// Todo delete h.AllTags param and make it in struct or other ideas
// Todo delete task after 2 days of due date
// Todo work from due date format make it more simple
// Todo Update task handler + route + storage

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
	// init middlewares
	initMiddlewares(httpServer)
	// Init routes
	initRoutes(httpServer)
	// Start server
	httpServer.Start(cfg, log)
}

func initMiddlewares(server *server.Server) {
	router := server.Router
	// recovery from panic
	router.Use(middleware.PanicRecovery(server.Log))
	// request execution time
	router.Use(middleware.HandlerExecutionTime(server.Log))
}

// initRoutes init routes for server
func initRoutes(server *server.Server) {
	router := server.Router
	router.Route("/task", func(r chi.Router) {
		// get all tasks
		r.Get("/", server.TaskHandlers.GetTasksHandler)
		// get task by id
		r.Get("/{id:\\d*}", server.TaskHandlers.GetTaskHandler)
		// get all tasks by tag or tags list with different modes.
		// full - returns tasks who have  includes specified tag or tags list in the task tags.
		// short - returns tasks who have only specified tag or tags list included in the task tags.
		// tags - in query using , as separator
		// due - in query format: 2006-01-02T15:04:05Z
		r.Route("/tag", func(r chi.Router) {
			r.Get("/{mode:(?:short|full)}/", server.TaskHandlers.GetTasksByModeAndTagHandler)
			r.Get("/", server.TaskHandlers.GetTasksByTagHandler)
		})
		// get tasks by due date
		r.Get("/{due:[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}(?::|%3A)[0-9]{2}(?::|%3A)[0-9]{2}Z}", server.TaskHandlers.GetTasksByDueDateHandler)

		// create new task using request body data
		// request body example:
		// {"text": "text", "tags": ["tag", "tag", "tag"], "due": "2021-01-01 00:00:00"}
		r.Post("/", server.TaskHandlers.CreateTaskHandler)

		// delete task by id
		r.Delete("/{id:[0-9]*}", server.TaskHandlers.DeleteTaskHandler)
		// delete all tasks
		r.Delete("/", server.TaskHandlers.DeleteTasksHandler)
	})
	router.Route("/tag", func(r chi.Router) {
		// get all tags
		r.Get("/", server.TaskHandlers.GetTagsHandler)
		// get tag by name
		r.Get("/{name:[A-Za-z]+}", server.TaskHandlers.GetTagHandler)

		// create new tag using request body data
		// request body example:
		// {"name": "name"}
		r.Post("/", server.TaskHandlers.CreateTagHandler)

		// delete all tags
		r.Delete("/", server.TaskHandlers.DeleteTagsHandler)
		// delete tag by name
		r.Delete("/{name:[A-Za-z]+}", server.TaskHandlers.DeleteTagHandler)
	})
	router.MethodNotAllowed(server.TaskHandlers.MethodNotAllowedHandler)
	router.NotFound(server.TaskHandlers.NotFoundHandler)
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

	log.Info("Successfully connected to database", slog.String("type", cfg.Type))
	return db
}
