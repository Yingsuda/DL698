package DLContorl

import "fmt"

type ControlHandle func(oad string, value interface{}) error

var Control_Map = make(map[string]ControlHandle)

func init() {

}

func RegisterControlHandle(oad string, f func(oad string, value interface{}) error) error {
	fmt.Println("Register OAD is ", oad)
	Control_Map[oad] = f
	return nil
}

func DoControl(oad string, value interface{}) error {
	if f, ok := Control_Map[oad]; ok {
		return f(oad, value)
	} else {
		fmt.Println("DoControl OAD Not Exist")
		return fmt.Errorf("oad %s not support control", oad)
	}
}
