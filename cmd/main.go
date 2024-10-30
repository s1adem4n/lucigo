package main

import (
	"context"
	"fmt"
	"log/slog"
	"lucigo/pkg/db"
)

func main() {
	conn, err := db.OpenDB()
	if err != nil {
		slog.Error("failed to open database", "error", err)
		return
	}
	defer conn.Close()

	queries := db.New(conn)

	ctx := context.Background()
	user, err := queries.GetUser(ctx, "admin")
	fmt.Println(user, err)
}
