package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
	"time"
)

func (setup *FabricSetup) Invoke(key, value string) (string, error) {

	// Prepare arguments
	var args []string
	args = append(args, "invoke")
	args = append(args, "invoke")
	args = append(args, key)
	args = append(args, value)

	eventID := "eventInvoke"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte(fmt.Sprintf("Transient data in %v invoke", key))

	// Register a notification handler on the client
	notifier := make(chan *chclient.CCEvent)
	rce, err := setup.Client.RegisterChaincodeEvent(notifier, setup.ChannelID, eventID)
	if err != nil {
		fmt.Errorf("注册链码事件异常：%v", err)
	}

	response, err := setup.Client.Execute(chclient.Request{ChaincodeID: setup.ChannelID, Fcn: args[0], Args: [][]byte{[]byte(args[1]), []byte(args[2]), []byte(args[3])}, TransientMap: transientDataMap})
	if err != nil {
		return "", fmt.Errorf("执行交易失败: %v", err)
	}

	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %s\n", ccEvent)
	case <-time.After(time.Second * 20):
		return "", fmt.Errorf("没有收到CC事件(%s)", eventID)
	}

	// Unregister the notification handler previously created on the client
	err = setup.Client.UnregisterChaincodeEvent(rce)

	// 成功并返回交易ID
	return response.TransactionID.ID, nil
}
