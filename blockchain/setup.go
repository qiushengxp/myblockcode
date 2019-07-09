package blockchain

import (
	"fmt"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chclient"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/chmgmtclient"
	"github.com/hyperledger/fabric-sdk-go/api/apitxn/resmgmtclient"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/ccpackager/gopackager"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabsdk"
	"github.com/hyperledger/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/cauthdsl"
	"time"
)

//定义结构体
type FabricSetup struct {
	ConfigFile      string //sdk配置文件所在路径
	ChannelID       string //应用通道名称
	ChannelConfig   string //应用通道交易配置文件所在路径
	ChaincodePath   string // 链码地址
	ChaincodeGoPath string // 链码所在根目录
	OrgAdmin        string // 组织管理员名称
	OrgName         string //组织名称
	UserName        string	// 用户名
	Initialized     bool   //是否初始化
	Client          chclient.ChannelClient
	Admin           resmgmtclient.ResourceMgmtClient //fabric环境中资源管理者
	SDK             *fabsdk.FabricSDK                //SDK实例
}

//1. 创建SDK实例并使用SDK实例创建应用通道，将Peer节点加入到创建的应用通道中
func (f *FabricSetup) Initialize() error {
	if f.Initialized {
		return fmt.Errorf("SDK 已实例化")
	}

	// 创建 sdfk
	sdk, err := fabsdk.New(config.FromFile(f.ConfigFile))
	if err != nil {
		return fmt.Errorf("SDK实例化失败:v%", err)
	}
	f.SDK = sdk
	//创建一个具有管理权限的应用通道客户端管理对象
	chmClient, err := f.SDK.NewClient(fabsdk.WithUser(f.OrgAdmin), fabsdk.WithOrg(f.OrgName)).ChannelMgmt()
	if err != nil {
		return fmt.Errorf("创建应用通道管理客户端管理对象失败:%v", err)
	}

	session, err := f.SDK.NewClient(fabsdk.WithUser(f.OrgAdmin), fabsdk.WithOrg(f.OrgName)).Session()
	if err != nil {
		return fmt.Errorf("获取当前会话用户对象失败:%v", err)
	}

	orgAdminUser := session

	//指定创建应用通道所需要的所有参数
	/*
	   $ peer channel create -o orderer.example.com:7050 -c $CHANNEL_NAME -f ./channel-artifacts/channel.tx --tls --cafile \
	   /opt/gopath/src/github.com/hyperledger/fabric/peer/crypto/ordererOrganizations/example.com/orderers/orderer.example.com/msp/tlscacerts/tlsca.example.com-cert.pem
	*/
	chReq := chmgmtclient.SaveChannelRequest{ChannelID: f.ChannelID, ChannelConfig: f.ChannelConfig, SigningIdentity: orgAdminUser}

	if err := chmClient.SaveChannel(chReq); err != nil {
		return fmt.Errorf("创建应用通道失败:%v", err)
	}
	time.Sleep(time.Second * 5)

	f.Admin, err = f.SDK.NewClient(fabsdk.WithUser(f.OrgAdmin)).ResourceMgmt()
	if err != nil {
		return fmt.Errorf("peer加入节点失败:%v", err)
	}
	f.Initialized = true
	fmt.Println("SDK 实例化成功")
	return nil
}

func (setup *FabricSetup) InstallAndInstantiateCC() error {

	// Create a new go lang chaincode package and initializing it with our chaincode
	ccPkg, err := gopackager.NewCCPackage(setup.ChaincodePath, setup.ChaincodeGoPath)
	if err != nil {
		return fmt.Errorf("failed to create chaincode package: %v", err)
	}

	// Install our chaincode on org peers
	// The resource management client send the chaincode to all peers in its channel in order for them to store it and interact with it later
	installCCReq := resmgmtclient.InstallCCRequest{Name: setup.ChannelID, Path: setup.ChaincodePath, Version: "1.0", Package: ccPkg}
	_, err = setup.Admin.InstallCC(installCCReq)
	if err != nil {
		return fmt.Errorf("failed to install cc to org peers %v", err)
	}

	// Set up chaincode policy
	// The chaincode policy is required if your transactions must follow some specific rules
	// If you don't provide any policy every transaction will be endorsed, and it's probably not what you want
	// In this case, we set the rule to : Endorse the transaction if the transaction have been signed by a member from the org "org1.hf.chainhero.io"
	ccPolicy := cauthdsl.SignedByAnyMember([]string{"org1.example.com"})

	// Instantiate our chaincode on org peers
	// The resource management client tells to all peers in its channel to instantiate the chaincode previously installed
	err = setup.Admin.InstantiateCC(setup.ChannelID, resmgmtclient.InstantiateCCRequest{Name: setup.ChannelID, Path: setup.ChaincodePath, Version: "1.0", Args: [][]byte{[]byte("init")}, Policy: ccPolicy})
	if err != nil {
		return fmt.Errorf("failed to instantiate the chaincode: %v", err)
	}

	// Channel client is used to query and execute transactions
	setup.Client, err = setup.SDK.NewClient(fabsdk.WithUser(setup.UserName)).Channel(setup.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to create new channel client: %v", err)
	}

	fmt.Println("Chaincode Installation & Instantiation Successful")
	return nil
}
