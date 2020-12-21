// +build ignore

package main

import (
	"log"

	"github.com/facebook/ent/entc"
	"github.com/facebook/ent/entc/gen"
	"github.com/facebookincubator/ent-contrib/entgql"
)

func main() {
    err := entc.Generate("./schema", &gen.Config{
        Templates: entgql.AllTemplates,
    })
    if err != nil {
        log.Fatalf("running ent codegen: %v", err)
    }
}