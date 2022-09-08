//go:build ignore
// +build ignore

package main

import (
	"log"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
)

func main() {

	ex, err := entgql.NewExtension(
		entgql.WithConfigPath("./graphql/gqlgen.yml"),
		// Generate GQL schema from the Ent's schema.
	//	entgql.WithSchemaGenerator(),
	// Generate the filters to a separate schema
	// file and load it in the gqlgen.yml config.
	//entgql.WithSchemaPath("./ent.graphql"),
	//entgql.WithWhereFilters(true),
	)
	if err != nil {
		log.Fatalf("creating entgql extension: %v", err)
	}
	opts := []entc.Option{
		entc.Extensions(ex),
	}
	err = entc.Generate("./ent/_schema", &gen.Config{
		Features: []gen.Feature{
			// gen.FeaturePrivacy,
			gen.FeatureUpsert,
		},
		// Templates: []*gen.Template{
		// 	gen.MustParse(gen.NewTemplate("static").
		// 		Funcs(template.FuncMap{"title": strings.ToTitle}).
		// 		ParseFiles("template/static.tmpl")),
		// },
	}, opts...)
	if err != nil {
		log.Fatalf("running ent codegen: %v", err)
	}
}
