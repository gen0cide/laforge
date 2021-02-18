// +build ignore

package main

import (
	"log"

	"github.com/facebook/ent/entc"
	"github.com/facebook/ent/entc/gen"
	"github.com/facebookincubator/ent-contrib/entgql"
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
