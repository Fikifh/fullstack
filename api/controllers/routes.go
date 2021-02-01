package controllers

import "github.com/fikifh/fullstack/api/middlewares"

func (s *Server) initializeRoutes() {
	//Home Route
	s.Router.HandleFunc("/", middlewares.SetMiddlewareJSON(s.Home)).Methods("GET")

	//Login route
	s.Router.HandleFunc("/login", middlewares.SetMiddlewareJSON(s.Login)).Methods("POST")

	//Users routes
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.CreateUser)).Methods("POST")
	s.Router.HandleFunc("/users", middlewares.SetMiddlewareJSON(s.GetUsers)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/users/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdateUser))).Methods("PUT")
	s.Router.HandleFunc("/", middlewares.SetMiddlewareAuthentication(s.DeleteUser)).Methods("DELETE")

	//Post routes
	s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.CreatedPost)).Methods("POST")
	s.Router.HandleFunc("/posts", middlewares.SetMiddlewareJSON(s.getPosts)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(s.GetUser)).Methods("GET")
	s.Router.HandleFunc("/posts/{id}", middlewares.SetMiddlewareJSON(middlewares.SetMiddlewareAuthentication(s.UpdatePost))).Methods("PUT")
	s.Router.HandleFunc("/", middlewares.SetMiddlewareAuthentication(s.DeletePost)).Methods("DELETE")

}
