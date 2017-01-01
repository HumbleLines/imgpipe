// Package imageops pkg/imageops/imageops.go
package imageops

// Handler defines a function that transforms image bytes and returns result/error.
type Handler func([]byte) ([]byte, error)

// Pipeline enables functional-style chained transformations on image bytes.
type Pipeline struct {
	steps []Handler
}

// NewPipeline constructs a pipeline for sequential image processing operations.
func NewPipeline() *Pipeline {
	return &Pipeline{}
}

// Add appends a processing handler to the chain.
func (p *Pipeline) Add(fn Handler) *Pipeline {
	p.steps = append(p.steps, fn)
	return p
}

// Run executes the handler pipeline on the provided data.
func (p *Pipeline) Run(data []byte) ([]byte, error) {
	var err error
	for _, step := range p.steps {
		data, err = step(data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}
