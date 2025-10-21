package main

import (
	"context"
	"fmt"

	"github.com/adamwreuben/twiggasdk/twigga"
)

func main() {
	twiggaClient, err := twigga.NewTwiggaClient("./twigga/bongo.json")
	if err != nil {
		fmt.Println("error**: ", err.Error())
		return
	}

	// dataToAdd := map[string]interface{}{
	// 	"database": "0cDQnjdbLwGPZsCAfh87",
	// 	"table":    "test_collection",
	// }

	// docId := fmt.Sprintf(`%s_%s`, "0cDQnjdbLwGPZsCAfh87", "test_collection")

	// resp, err := twiggaClient.CreateDocumentWithID(context.Background(), "ChangeFeedsData", docId, dataToAdd)
	// if err != nil {
	// 	fmt.Println("CreateDocumentAuto Error: ", err.Error())
	// }

	// fmt.Println("resp: ", string(resp))

	// found, err := twiggaClient.DocumentExists(context.Background(), "Applications", dataToCheck)
	// if err != nil {
	// 	fmt.Println("documentFound Error: ", err.Error())
	// }

	// log.Println("Found: ", found)

	// qresp, _ := twiggaClient.QueryDocuments(context.Background(), "MeterRegistration", map[string]interface{}{
	// 	"imei": "862273041539147",
	// })

	// docsRaw, ok := qresp["documents"].([]interface{})
	// if !ok || len(docsRaw) == 0 {
	// 	log.Printf("No documents found for IMEI: %s", "862273041539147")
	// 	return
	// }

	// if len(docsRaw) > 0 {
	// 	doc, _ := docsRaw[0].(map[string]interface{})
	// 	meterNumber := doc["meterNumber"].(string)
	// 	fmt.Println("meterNumber: ", meterNumber)
	// }

	// res, err := twiggaClient.DocumentExists(context.Background(), "ChangeFeedsData", dataToAdd)
	// if err != nil {
	// 	fmt.Println("Error*: ", err.Error())
	// 	return
	// }

	// fmt.Println(res)

	// res, err := twiggaClient.Login(context.Background(), "jack@jack.com", "123")
	// if err != nil {
	// 	fmt.Println("Error*: ", err.Error())
	// 	return
	// }

	// twiggaResp, _ := twiggaClient.GetDocument(context.Background(), "Inbox", "adamreuben@bongocloud.co.tz_Inbox")

	// var data map[string]interface{}
	// json.Unmarshal(twiggaResp, &data)

	// messages := data["messages"].([]interface{})
	// fmt.Println("messages: ", messages[0])

	dataToFilter := map[string]interface{}{
		"name": "Willy William",
	}
	res, _ := twiggaClient.QueryDocuments(context.Background(), "Students", dataToFilter)
	fmt.Println("Response")
	fmt.Println(res)

}
