// +build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {
	opts := []entc.Option{
		entc.TemplateFiles("template/ent.tmpl"),
	}
	err := entc.Generate("./schema", &gen.Config{
		Templates: entgql.AllTemplates,
	}, opts...)
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
