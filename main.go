package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

// User is structure represeting User.
// It could be stored in DB or fetched from REST, etc.
type User struct {
	ID    int
	Name  string
	Email string
}

// UserObject defines `User` object equivalent in graphql
var UserObject = graphql.NewObject(graphql.ObjectConfig{
	Name: "User",
	Fields: graphql.Fields{
		"id":    &graphql.Field{Type: graphql.Int},
		"name":  &graphql.Field{Type: graphql.String},
		"email": &graphql.Field{Type: graphql.String},
	},
})

// ListOfUsers is the list of users
var ListOfUsers = []User{
	{
		ID:    1,
		Name:  "Sai Varun Reddy",
		Email: "saivarunvishal@gmail.com",
	},
	{
		ID:    2,
		Name:  "Tony Stark",
		Email: "tony@starkindustries.inc",
	},
	{
		ID:   3,
		Name: "Mr. X",
	},
}

// LoginQuery is for login
var LoginQuery = graphql.Field{
	Name: "Login",
	Type: graphql.String,
	Args: graphql.FieldConfigArgument{
		"email": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
		"password": &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.String),
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		email := p.Args["email"].(string)
		password := p.Args["password"].(string)
		if email == "varun" && password == "123" {
			return "success", nil
		}
		return nil, errors.New("Invalid Login")
	},
}

// UserField contains user query type
var UserField = graphql.Field{
	Name: "User",
	Type: graphql.NewList(UserObject),
	Args: graphql.FieldConfigArgument{
		"limit": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		// Default limit to all users
		limit := len(ListOfUsers)
		if p.Args["limit"] != nil {
			limit = p.Args["limit"].(int)
		}
		return ListOfUsers[0:limit], nil
	},
}

// QueryFieldsContainer contains all `Query` types
var QueryFieldsContainer = graphql.Fields{
	"Users": &UserField,
}

// MutationContainer will have all mutation fields
var MutationContainer = graphql.Fields{
	"Login": &LoginQuery,
}

func main() {
	var queryType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Query",
			Fields: QueryFieldsContainer,
		},
	)
	var mutationType = graphql.NewObject(
		graphql.ObjectConfig{
			Name:   "Mutation",
			Fields: MutationContainer,
		},
	)
	var schema, err = graphql.NewSchema(
		graphql.SchemaConfig{
			Query:    queryType,
			Mutation: mutationType,
		},
	)
	if err != nil {
		log.Println(err)
	}
	h := handler.New(&handler.Config{
		Schema:   &schema,
		Pretty:   true,
		GraphiQL: true,
	})

	http.Handle("/graphql", h)
	http.ListenAndServe(":8088", nil)
}
