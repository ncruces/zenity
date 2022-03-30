package strftime

type parser struct {
	fmt      string
	specs    map[byte]string
	unpadded map[byte]string
	writeLit func(byte) error
	writeFmt func(string) error
	fallback func(spec byte, pad bool) error
}

func (p *parser) parse() error {
	const (
		initial = iota
		specifier
		nopadding
	)

	var err error
	state := initial
	for _, b := range []byte(p.fmt) {
		switch state {
		default:
			if b == '%' {
				state = specifier
			} else {
				err = p.writeLit(b)
			}

		case specifier:
			if b == '-' {
				state = nopadding
				continue
			}
			if s, ok := p.specs[b]; ok {
				err = p.writeFmt(s)
			} else {
				err = p.fallback(b, true)
			}
			state = initial

		case nopadding:
			if s, ok := p.unpadded[b]; ok {
				err = p.writeFmt(s)
			} else if s, ok := p.specs[b]; ok {
				err = p.writeFmt(s)
			} else {
				err = p.fallback(b, false)
			}
			state = initial
		}

		if err != nil {
			return err
		}
	}

	switch state {
	case specifier:
		return p.writeLit('%')
	case nopadding:
		return p.writeLit('-')
	}
	return nil
}
