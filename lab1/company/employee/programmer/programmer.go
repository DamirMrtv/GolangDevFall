package programmer

type Programmer struct {
	position string
	salary   float64
	address  string
}

func (p *Programmer) GetPosition() string {
	return p.position
}

func (p *Programmer) SetPosition(pos string) {
	p.position = pos
}

func (p *Programmer) GetSalary() float64 {
	return p.salary
}

func (p *Programmer) SetSalary(sal float64) {
	p.salary = sal
}

func (p *Programmer) GetAddress() string {
	return p.address
}

func (p *Programmer) SetAddress(addr string) {
	p.address = addr
}
