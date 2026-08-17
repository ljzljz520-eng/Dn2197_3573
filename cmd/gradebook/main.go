package main

import (
	"fmt"
	"log"
	"os"

	"gradebook/fixtures"
	"gradebook/service"
	"gradebook/storage"
)

func main() {
	path := "gradebook.db"
	if len(os.Args) > 1 {
		path = os.Args[1]
	}
	store, err := storage.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer store.Close()
	app := service.NewGradebookService(store)
	if records, e := app.List(); e != nil {
		log.Fatal(e)
	} else if len(records) == 0 {
		for _, record := range fixtures.Batch() {
			if e := app.Add(record, "cli"); e != nil {
				log.Fatal(e)
			}
		}
	}
	report, err := app.BuildReport("Graduate Gradebook")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(service.RenderReport(report))
}
