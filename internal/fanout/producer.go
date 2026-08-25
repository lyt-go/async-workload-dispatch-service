package fanout

import "errors"

type Item struct {
	ID      string
	Invalid bool
}

func Produce(items []Item) (<-chan Item, <-chan error) {
	output := make(chan Item)
	errorsOut := make(chan error)
	go func() {
		for _, item := range items {
			if item.Invalid {
				errorsOut <- errors.New("invalid dispatch item: " + item.ID)
				return
			}
			output <- item
		}
		close(output)
	}()
	return output, errorsOut
}
