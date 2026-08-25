package batch

type Resource interface{ Close() error }
type Factory interface {
	Open(index int) (Resource, error)
}

type Processor struct{ Factory Factory }

func (p Processor) Process(count int, visit func(int, Resource) error) error {
	for i := 0; i < count; i++ {
		resource, err := p.Factory.Open(i)
		if err != nil {
			return err
		}
		defer resource.Close()
		if err := visit(i, resource); err != nil {
			return err
		}
	}
	return nil
}
