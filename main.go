package main

import (
	"context"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/block-vision/sui-go-sdk/constant"
	"github.com/block-vision/sui-go-sdk/models"
	"github.com/block-vision/sui-go-sdk/signer"
	"github.com/block-vision/sui-go-sdk/sui"
	"github.com/block-vision/sui-go-sdk/utils"
)

// package id:0xbda58f110ce755a63c007d68cc53f7ac68c780dc8fb1fb16ad52d797143b4799
//tx: 6S5b62crgfbikAEZKLwphkQqtf5wuJwQqkkRW4ywnxug
//obj:0x7dc80959fdd7b4c68ba0caa2f0f1182fb297817742caf170dc1f787d58317f3d

const (
	PackageId  = "0xbda58f110ce755a63c007d68cc53f7ac68c780dc8fb1fb16ad52d797143b4799"
	CountObjId = "0x7dc80959fdd7b4c68ba0caa2f0f1182fb297817742caf170dc1f787d58317f3d"
)

/*func init() {
	err := godotenv.Load(".env") // 默认读取 .env 文件
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}*/

func main() {
	a := app.New()
	w := a.NewWindow("counter")
	w.Resize(fyne.NewSize(600, 400))
	counter := GetCounter()
	counterLabel := widget.NewLabel(fmt.Sprintf("counter: %d", counter))
	txLabel := widget.NewLabel("transaction:")
	txLabel.Selectable = true
	txLabel.Wrapping = fyne.TextWrapWord
	addBtn := widget.NewButton("add", func() {
		tx := CallIncrement()
		counter := GetCounter()
		counterLabel.SetText(fmt.Sprintf("counter: %d", counter))
		txLabel.SetText(fmt.Sprintf("transaction: %s", tx))
	})

	ct := container.NewVBox(layout.NewSpacer(), counterLabel, txLabel, addBtn, layout.NewSpacer())
	w.SetContent(ct)
	w.ShowAndRun()
}

func CallIncrement() string {
	client := sui.NewSuiClient(constant.SuiTestnetEndpoint)
	ctx := context.Background()
	//secretKey := os.Getenv("SecretKey")
	mySigner, err := signer.NewSignerWithSecretKey(PriKey)
	if err != nil {
		panic(err)
	}
	rsp, err := client.MoveCall(ctx, models.MoveCallRequest{
		Signer:          mySigner.Address,
		PackageObjectId: "0xbda58f110ce755a63c007d68cc53f7ac68c780dc8fb1fb16ad52d797143b4799",
		Module:          "counter",
		Function:        "increment",
		TypeArguments:   []interface{}{},
		Arguments: []interface{}{
			"0x7dc80959fdd7b4c68ba0caa2f0f1182fb297817742caf170dc1f787d58317f3d",
		},
		GasBudget:     "100000000",
		ExecutionMode: "",
	})
	if err != nil {
		panic(err)
	}

	rsp2, err := client.SignAndExecuteTransactionBlock(ctx, models.SignAndExecuteTransactionBlockRequest{
		TxnMetaData: rsp,
		PriKey:      mySigner.PriKey,
		// only fetch the effects field
		Options: models.SuiTransactionBlockOptions{
			ShowInput:    true,
			ShowRawInput: true,
			ShowEffects:  true,
		},
		RequestType: "WaitForLocalExecution",
	})

	if err != nil {
		fmt.Println(err.Error())
		return ""
	}

	utils.PrettyPrint(rsp2)
	return rsp2.Digest
}

func GetCounter() uint8 {
	client := sui.NewSuiClient(constant.SuiTestnetEndpoint)
	object, err := client.SuiGetObject(context.Background(), models.SuiGetObjectRequest{
		ObjectId: CountObjId,
		Options: models.SuiObjectDataOptions{
			ShowContent: true,
		},
	})
	if err != nil {
		panic(err)
	}
	fields := object.Data.Content.Fields
	counter := fields["counter"].(float64)
	fmt.Printf("counter:%v", counter)
	return uint8(counter)
}
