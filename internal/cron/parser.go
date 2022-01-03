package cron

// parser a simple recursive descent parser for cron
// time elements
type parser struct {
	min    int
	max    int
	pos    int
	vals   []int
	tokens []token
}

// processList
// range comma range
// range
func (p *parser) processList() {
	for p.pos < len(p.tokens) {
		p.commaNext()
		numbers := p.processRangeWithStep()
		p.vals = append(p.vals, numbers...)
	}
}

// processRangeWithStep
func (p *parser) processRangeWithStep() []int {
	step := 1
	numbers := p.processRange()
	if len(numbers) == 0 {
		return numbers
	}
	if p.pos < len(p.tokens) {
		if p.slashNext() {
			if p.pos < len(p.tokens) && p.tokens[p.pos].ttype == number {
				step = p.tokens[p.pos].num
				p.pos++
			}
		}
	}
	min := numbers[0]
	max := min
	if len(numbers) > 1 {
		max = numbers[1]
	}

	// generate the list
	vals := []int{}
	for i := numbers[0]; i <= max; i += step {
		vals = append(vals, i)
	}
	return vals
}

// processRange
// num dash num
// num
// star
func (p *parser) processRange() []int {
	number, ok := p.processVal()
	if ok {
		// hack to deal with *
		if number == -1 {
			return []int{p.min, p.max}
		}
		if p.dashNext() {
			nextNumber, ok := p.processVal()
			if ok {
				return []int{number, nextNumber}
			}
		}
		return []int{number}
	}
	return []int{}
}

// processVal returns the 'number' or -1 if we have a 0
// and or false if not a numeric value
func (p *parser) processVal() (int, bool) {
	if p.pos >= len(p.tokens) {
		return 0, false
	}
	t := p.tokens[p.pos]
	switch t.ttype {
	case number:
		p.pos++
		return t.num, true
	case star:
		p.pos++
		return -1, true
	}
	return 0, false
}

// commaNext return true if there is a comma next (and moves
// on pointer to next token if so)
func (p *parser) commaNext() bool {
	return p.typeNext(comma)
}

// slashNext return true if there is a slash next (and moves
// on pointer to next token if so)
func (p *parser) slashNext() bool {
	return p.typeNext(slash)
}

// dashNext return true if there is a dash next (and moves
// on pointer to next token if so)
func (p *parser) dashNext() bool {
	return p.typeNext(dash)
}

// typeNext return true if there is the type specified next
// (and moves on pointer to next token if so)
func (p *parser) typeNext(t tokenType) bool {
	if p.pos >= len(p.tokens) {
		return false
	}
	tok := p.tokens[p.pos]
	if tok.ttype == t {
		p.pos++
		return true
	}
	return false
}
