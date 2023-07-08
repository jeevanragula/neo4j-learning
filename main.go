package main

import (
	"context"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
	"log"
	"os"
	"time"
)

func main() {
	password := os.Getenv("password")
	ctx := context.Background()
	driver, err := neo4j.NewDriverWithContext(
		"neo4j://localhost:7687",               // (1)
		neo4j.BasicAuth("neo4j", password, ""), // (2)
	)
	err = driver.VerifyConnectivity(ctx)
	if err != nil {
		panic(err)
	}

	// Open a new Session
	session := driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func(driver neo4j.DriverWithContext, ctx context.Context) {
		err := driver.Close(ctx)
		if err != nil {
			panic(err)
		}
	}(driver, ctx)

	defer func(session neo4j.SessionWithContext, ctx context.Context) {
		err := session.Close(ctx)
		if err != nil {
			panic(err)
		}
	}(session, ctx)

	result, err := session.Run(
		ctx,
		`
	MATCH (p:Person)-[:DIRECTED]->(:Movie {title: $title})
	RETURN p
	`,
		map[string]any{"title": "The Matrix"}, // (3)
		func(txConfig *neo4j.TransactionConfig) {
			txConfig.Timeout = 3 * time.Second // (4)
		},
	)
	PanicOnErr(err)

	people, err := neo4j.CollectTWithContext[neo4j.Node](ctx, result,
		func(record *neo4j.Record) (neo4j.Node, error) {
			person, _, err := neo4j.GetRecordValue[neo4j.Node](record, "p")
			return person, err
		})
	PanicOnErr(err)
	fmt.Println(people)

}

func PanicOnErr(err error) {
	if err != nil {
		log.Println(err)
		panic(err)
	}
}
