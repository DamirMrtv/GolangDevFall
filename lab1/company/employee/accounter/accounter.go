package accounter

type Accounter struct {
	position string
	salary   float64
	address  string
}

func (a *Accounter) GetPosition() string {
	return a.position
}

func (a *Accounter) SetPosition(pos string) {
	a.position = pos
}

func (a *Accounter) GetSalary() float64 {
	return a.salary
}

func (a *Accounter) SetSalary(sal float64) {
	a.salary = sal
}

func (a *Accounter) GetAddress() string {
	return a.address
}

func (a *Accounter) SetAddress(addr string) {
	a.address = addr
}
