package services

import (
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/LilLebowski/loyalty-system/internal/clients"
	"github.com/LilLebowski/loyalty-system/internal/utils"
)

const timeoutErrTooManyRequests = 2 * time.Minute

type PoolWorker struct {
	client       *clients.Client
	orderService *Order
	orderIn      chan string
	Err          chan error
}

func NewPoolWorker(client *clients.Client, serviceUser *Order) *PoolWorker {
	ordersIn := make(chan string, 10)
	err := make(chan error)
	return &PoolWorker{client: client, orderService: serviceUser, orderIn: ordersIn, Err: err}
}

func (p *PoolWorker) StarIntegration(countWorker int, requestTime *time.Ticker) {
	pauses := make([]chan struct{}, 0)
	for i := 0; i < countWorker; i++ {
		name := i
		pause := make(chan struct{})
		p.worker(name, pause)
		pauses = append(pauses, pause)
	}

	go func() {
		for range requestTime.C {
			numbers, err := p.orderService.GetOrdersNotProcessed()
			if err != nil {
				log.Println(err)
				break
			}
			for _, n := range numbers {
				number := n
				p.orderIn <- number
			}
		}
	}()

	for err := range p.Err {
		log.Printf("error %s", err.Error())
		if errors.Is(err, utils.ErrTooManyRequests) {
			go func() {
				for _, pause := range pauses {
					ch := pause
					ch <- struct{}{}
				}
			}()
		}
	}
}

func (p *PoolWorker) worker(nameWorker int, pause chan struct{}) {
	go func() {
		defer close(pause)
		for {
			select {
			case order := <-p.orderIn:
				log.Printf("worker %d, order %s send request to accrual services", nameWorker, order)
				accrual, err := p.client.CheckAccrual(order)
				if err != nil {
					p.Err <- fmt.Errorf("error worker %d %w", nameWorker, err)
					break
				}
				log.Printf("worker %d, save %v in order", nameWorker, accrual)
				err = p.orderService.UpdateOrder(accrual)
				if err != nil {
					p.Err <- fmt.Errorf("error worker %d %w", nameWorker, err)
				}
			case <-pause:
				log.Printf("worker %d do pause", nameWorker)
				time.Sleep(timeoutErrTooManyRequests)
			}
		}
	}()
}
