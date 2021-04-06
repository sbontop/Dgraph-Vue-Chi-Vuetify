# Dgraph-Vue-Chi-Vuetify
Full Stack Boilerplate Project built ontop tech stacks including graph-based database called DGraph, a robust and secure Backend in Go and an interactive Frontend in Vue.js

# Installation

## Frontend 
1. npm install -g @vue/cli
2. cd frontend/buyers-app/
3. npm install
4. vue-cli-service serve

## Backend
1. cd backend
2. go get "github.com/go-chi/chi"
3. go get "github.com/go-chi/chi/v5/middleware"
4. go get "github.com/go-chi/render"
5. go get "github.com/dgraph-io/dgo/v2"
6. go get "github.com/dgraph-io/dgo/v2/protos/api"
7. go get "google.golang.org/grpc"
8. go run main.go

## Dgraph Database
1. docker run --rm -it -p 8000:8000 -p 8080:8080 -p 9080:9080 dgraph/standalone:v20.11.3

