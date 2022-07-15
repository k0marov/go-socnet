package main

import (
	"fmt"
	"github.com/k0marov/go-socnet/core"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

const Port = 4242

func main() {
	http.ListenAndServe(fmt.Sprintf(":%v", Port), core.Setup())
}
