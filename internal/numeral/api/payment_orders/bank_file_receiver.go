package payment_orders

import (
	"context"
	"encoding/csv"
	"fmt"
	"log/slog"
	"os"
	"strings"

	"github.com/fsnotify/fsnotify"
)

type BankFileReceiver struct {
	inputLocation string
	service       *Service
	watcher       *fsnotify.Watcher
}

func NewBankFileReceiver(inputLocation string, service *Service) (*BankFileReceiver, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = watcher.Add(inputLocation)
	if err != nil {
		return nil, err
	}
	return &BankFileReceiver{
		inputLocation: inputLocation,
		service:       service,
		watcher:       watcher,
	}, nil
}

func (g BankFileReceiver) ReceiveFilesFromBank() error {
	done := make(chan bool)
	go func() {
		for {
			select {
			case <-done:
				return
			// watch for events
			case event := <-g.watcher.Events:
				if event.Op != fsnotify.Create || !strings.HasSuffix(event.Name, ".response.csv") {
					continue
				}
				fmt.Printf("EVENT! %#v\n", event)
				fileID := strings.TrimPrefix(strings.TrimSuffix(event.Name, ".response.csv"), g.inputLocation)
				file, err := os.Open(event.Name)
				if err != nil {
					fmt.Println(err)
				}
				reader := csv.NewReader(file)
				records, _ := reader.ReadAll()
				if len(records) != 2 {
					// for simplicity we'll just consider files with single payment
				}
				slog.Info("something")
				paymentDetails := records[1]
				status := paymentDetails[1]
				err = g.service.UpdateStatus(context.Background(), fileID, status)
				if err != nil {
					fmt.Println(err)
				}
				// watch for errors
			case err := <-g.watcher.Errors:
				fmt.Println("ERROR", err)
			}
		}
	}()

	// out of the box fsnotify can watch a single file, or a single directory
	if err := g.watcher.Add(g.inputLocation); err != nil {
		return err
	}

	return nil
}

func (g *BankFileReceiver) Stop() {
	g.watcher.Close()
}
