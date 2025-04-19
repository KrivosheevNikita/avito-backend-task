package handlers

import (
	"net/http"
	"pvz-service/internal/middleware"
	"pvz-service/internal/transport/handlers/auth"
	"pvz-service/internal/transport/handlers/product"
	"pvz-service/internal/transport/handlers/pvz"
	"pvz-service/internal/transport/handlers/reception"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/dummyLogin", auth.DummyLoginHandler).Methods("POST")
	r.HandleFunc("/register", auth.RegisterHandler).Methods("POST")
	r.HandleFunc("/login", auth.LoginHandler).Methods("POST")

	r.Handle("/pvz",
		middleware.RoleCheck("moderator")(
			http.HandlerFunc(pvz.CreatePVZHandler),
		),
	).Methods("POST")

	r.Handle("/pvz",
		middleware.RoleCheck("employee", "moderator")(
			http.HandlerFunc(pvz.GetPVZHandler),
		),
	).Methods("GET")

	r.Handle("/pvz/{pvzId}/close_last_reception",
		middleware.RoleCheck("employee")(
			http.HandlerFunc(pvz.CloseLastReceptionHandler),
		),
	).Methods("POST")

	r.Handle("/pvz/{pvzId}/delete_last_product",
		middleware.RoleCheck("employee")(
			http.HandlerFunc(product.DeleteLastProductHandler),
		),
	).Methods("POST")

	r.Handle("/receptions",
		middleware.RoleCheck("employee")(
			http.HandlerFunc(reception.CreateReceptionHandler),
		),
	).Methods("POST")

	r.Handle("/products",
		middleware.RoleCheck("employee")(
			http.HandlerFunc(product.AddProductHandler),
		),
	).Methods("POST")

	return r
}
