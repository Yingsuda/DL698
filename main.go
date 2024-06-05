package main

import (
	"dl698/model"
	"dl698/utils"
	"fmt"
)

func main() {
	e := model.NewElectricityMeter("192.168.30.118:6001", "12345678")

	err := e.Login()
	if err != nil {
		fmt.Println("Login error:", err)
		return
	}
	oad := "20000200"

	for i := 0; i < 3; i++ {
		up := model.UploadPoint{
			Oad:     oad,
			Uid:     i,
			Dt:      utils.DT_Uint16,
			Value:   100,
			IsArray: true,
		}
		//
		e.AddPoint(&up) //分项电压
		e.UpdateElectricityInfo(i, 2413)
	}

	e.Start()
}
