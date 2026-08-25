package runner

import "context"

type Dispatcher struct{ Attempts int }

func (d Dispatcher) Execute(ctx context.Context, work func(context.Context) error) <-chan error {
	result := make(chan error, 1)
	go func() {
		defer close(result)
		var err error
		for i := 0; i < d.Attempts; i++ {
			err = work(context.Background())
			if err == nil {
				break
			}
		}
		result <- err
	}()
	return result
}
