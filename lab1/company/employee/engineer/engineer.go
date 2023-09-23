package engineer

type Engineer struct {
	position string
	salary   float64
	address  string
}

func (e *Engineer) GetPosition() string {
	return e.position
}

func (e *Engineer) SetPosition(pos string) {
	e.position = pos
}

func (e *Engineer) GetSalary() float64 {
	return e.salary
}

func (e *Engineer) SetSalary(sal float64) {
	e.salary = sal
}

func (e *Engineer) GetAddress() string {
	return e.address
}

func (e *Engineer) SetAddress(addr string) {
	e.address = addr
}
