package main

import (
	"fmt"
	"github.com/qiushengxp/myblockcode/blockchain"
	"os"
)

func main() {
	fsetup := blockchain.FabricSetup{
		ConfigFile:      "config.yaml",
		ChannelID:       "qiushengchannel",
		ChannelConfig:   os.Getenv("GOPATH") + "/src/github.com/qiushengxp/myblockcode/fixtures/channel-artifacts/channel.tx",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath:   "github.com/qiushengxp/myblockcode/chaincode/",
		OrgAdmin:        "Admin",
		OrgName:         "Org1",
		UserName:        "User1",
	}

	err := fsetup.Initialize()
	if err != nil {
		_ = fmt.Errorf("Fabric SDK 初始化失败：%v", err)
		fmt.Println(err.Error())
		return
	}
	err = fsetup.InstallAndInstantiateCC()
	if err != nil {
		fmt.Errorf("Fabric SDK 初始化客户端失败：%v", err)
		fmt.Println(err.Error())
		return
	}
	txid, err := fsetup.Query("hello")
	if err != nil {
		fmt.Errorf("交易异常：%v", err)
		return
	}
	fmt.Println("交易成功，ID：", txid)
}
