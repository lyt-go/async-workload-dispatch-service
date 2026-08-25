package payload

type Parser struct{ buffer []byte }

func (p *Parser) Parse(value string) []byte {
	p.buffer = append(p.buffer[:0], value...)
	return p.buffer
}
