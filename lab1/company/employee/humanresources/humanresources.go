package humanresources

type HumanResources struct {
	position string
	salary   float64
	address  string
}

func (hr *HumanResources) GetPosition() string {
	return hr.position
}

func (hr *HumanResources) SetPosition(pos string) {
	hr.position = pos
}

func (hr *HumanResources) GetSalary() float64 {
	return hr.salary
}

func (hr *HumanResources) SetSalary(sal float64) {
	hr.salary = sal
}

func (hr *HumanResources) GetAddress() string {
	return hr.address
}

func (hr *HumanResources) SetAddress(addr string) {
	hr.address = addr
}
