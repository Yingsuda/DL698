package model

type ElectricityMeter struct {
	//电压
	voltage float64
	ch      chan struct{}
}

func NewElectricityMeter() *ElectricityMeter {
	em := &ElectricityMeter{
		ch: make(chan struct{}, 1),
	}
	em.ch <- struct{}{}
	return em
}

func (e *ElectricityMeter) UpdateElectricityInfo() {
	<-e.ch
	defer func() {
		e.ch <- struct{}{}
	}()

	//更新数据

}

func (e *ElectricityMeter) GetElectricityInfo() {
	<-e.ch
	defer func() {
		e.ch <- struct{}{}
	}()

	//获取数据

}
