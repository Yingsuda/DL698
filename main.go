package main

import (
	"dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/model"
	_ "dev.magustek.com/bigdata/dass/iotdriver/OP2_DL_698/sync"
	"gitee.com/iotdrive/leader.upload/loop"
	"gitee.com/iotdrive/leader.upload/model/upload"
	uploadConst "gitee.com/iotdrive/leader.upload/utils/const"
	"gitee.com/iotdrive/tools/logs"
	"time"
)

func main() {
	uploadKey := "DL698"
	upload.RegeditUpload(uploadKey, model.NewUpload698())
	l := &loop.Loop{
		DasName:       uploadConst.PROVIDE_DAS_KEY_SUB,
		UploadName:    uploadKey,
		AutoPushPoint: false,
	}

	err := l.LoopInit()
	if err != nil {
		logs.Error(err)
		return
	}
	l.Loop()
	//通知协程推出
	model.Cancel()
	time.Sleep(time.Second)
	logs.Error("upload all goroutine exit")
}
