package Controller

import (
	"XGQ/1_DistributedLedgerDemo/Model"
	"crypto/ecdsa"
	"crypto/elliptic"
	"os"
	"strconv"

	"github.com/Unknwon/goconfig"
)

// All funtion in this file are test context

func fixAccount(amgr *AccountMgr) {
	CfgPath := ".\\data\\AccountCfg.ini"

	_, err := os.Stat(CfgPath)

	if err != nil {
		amgr.AddAccount("张三", &amgr.pksf)
		amgr.AddAccount("李四", &amgr.pksf)
		amgr.AddAccount("王五", &amgr.pksf)

		os.Create(CfgPath)
		cfg, _ := goconfig.LoadConfigFile(CfgPath)

		for key, value := range amgr.AccountArr {
			cfg.SetValue(key, "PublicKeyHash", value.PublicKeyHash)
			cfg.SetValue(key, "PublicKey.X", value.PublicKey.X.String())
			cfg.SetValue(key, "PublicKey.Y", value.PublicKey.Y.String())

			privateKey := value.GetPrivateKey()
			cfg.SetValue(key, "privateKey.X", privateKey.X.String())
			cfg.SetValue(key, "privateKey.Y", privateKey.Y.String())
			cfg.SetValue(key, "privateKey.D", privateKey.D.String())

			cfg.SetValue(key, "Nickname", value.Nickname)
			cfg.SetValue(key, "Balance", strconv.FormatFloat(value.Balance, 'E', -1, 64))
			cfg.SetValue(key, "UtxoHash", value.UtxoHash)
		}

		goconfig.SaveConfigFile(cfg, CfgPath)

		return
	}

	var privateKey1, privateKey2, privateKey3 ecdsa.PrivateKey
	cfg, _ := goconfig.LoadConfigFile(CfgPath)

	// Address1
	xgqA := Model.XgqAccount{}
	Addr := "e0fc0a12c2cfd7ee4ae3cc718fb84d152ff2d53057334b2ab9939b5282e3b544"
	xgqA.PublicKeyHash, _ = cfg.GetValue(Addr, "PublicKeyHash")
	csv, _ := cfg.GetValue(Addr, "PublicKey.X")
	xgqA.PublicKey.X = str2BigInt(csv)
	csv, _ = cfg.GetValue(Addr, "PublicKey.Y")
	xgqA.PublicKey.Y = str2BigInt(csv)
	xgqA.PublicKey.Curve = elliptic.P521()

	csv, _ = cfg.GetValue(Addr, "privateKey.X")
	privateKey1.X = str2BigInt(csv)
	csv, _ = cfg.GetValue(Addr, "privateKey.Y")
	privateKey1.Y = str2BigInt(csv)
	csv, _ = cfg.GetValue(Addr, "privateKey.D")

	privateKey1.Curve = elliptic.P521()
	privateKey1.D = str2BigInt(csv)
	xgqA.SetPrivateKey(privateKey1)

	xgqA.Nickname, _ = cfg.GetValue(Addr, "Nickname")
	csv, _ = cfg.GetValue(Addr, "Balance")
	xgqA.Balance, _ = strconv.ParseFloat(csv, 64)
	xgqA.UtxoHash, _ = cfg.GetValue(Addr, "UtxoHash")
	amgr.AccountArr[xgqA.PublicKeyHash] = xgqA

	// Address2
	xgqA = Model.XgqAccount{}
	Addr = "5d3fd2cc8a3c491c87b5e834d5970c99921ee177b2f387d4fe70fffeafa9d16f"
	xgqA.PublicKeyHash, _ = cfg.GetValue(Addr, "PublicKeyHash")
	csv, _ = cfg.GetValue(Addr, "PublicKey.X")
	xgqA.PublicKey.X = str2BigInt(csv)
	csv, _ = cfg.GetValue(Addr, "PublicKey.Y")
	xgqA.PublicKey.Y = str2BigInt(csv)
	xgqA.PublicKey.Curve = elliptic.P521()

	csv, _ = cfg.GetValue(Addr, "privateKey.X")
	privateKey2.X = str2BigInt(csv)
	csv, _ = cfg.GetValue(Addr, "privateKey.Y")
	privateKey2.Y = str2BigInt(csv)
	csv, _ = cfg.GetValue(Addr, "privateKey.D")

	privateKey2.Curve = elliptic.P521()
	privateKey2.D = str2BigInt(csv)
	xgqA.SetPrivateKey(privateKey2)

	xgqA.Nickname, _ = cfg.GetValue(Addr, "Nickname")
	csv, _ = cfg.GetValue(Addr, "Balance")
	xgqA.Balance, _ = strconv.ParseFloat(csv, 64)
	xgqA.UtxoHash, _ = cfg.GetValue(Addr, "UtxoHash")
	amgr.AccountArr[xgqA.PublicKeyHash] = xgqA

	// Address3
	xgqA = Model.XgqAccount{}
	Addr = "4d959a2efa604d7aa8fa56cd1e499ff42c527ea64fa6a65429aca454d156dadd"
	xgqA.PublicKeyHash, _ = cfg.GetValue(Addr, "PublicKeyHash")
	csv, _ = cfg.GetValue(Addr, "PublicKey.X")
	xgqA.PublicKey.X = str2BigInt(csv)
	csv, _ = cfg.GetValue(Addr, "PublicKey.Y")
	xgqA.PublicKey.Y = str2BigInt(csv)
	xgqA.PublicKey.Curve = elliptic.P521()

	csv, _ = cfg.GetValue(Addr, "privateKey.X")
	privateKey3.X = str2BigInt(csv)
	csv, _ = cfg.GetValue(Addr, "privateKey.Y")
	privateKey3.Y = str2BigInt(csv)
	csv, _ = cfg.GetValue(Addr, "privateKey.D")
	privateKey3.D = str2BigInt(csv)

	privateKey3.Curve = elliptic.P521()
	xgqA.SetPrivateKey(privateKey3)

	xgqA.Nickname, _ = cfg.GetValue(Addr, "Nickname")
	csv, _ = cfg.GetValue(Addr, "Balance")
	xgqA.Balance, _ = strconv.ParseFloat(csv, 64)
	xgqA.UtxoHash, _ = cfg.GetValue(Addr, "UtxoHash")
	amgr.AccountArr[xgqA.PublicKeyHash] = xgqA
}

func setCustomBalance(amgr *AccountMgr) {

	Addr := "5d3fd2cc8a3c491c87b5e834d5970c99921ee177b2f387d4fe70fffeafa9d16f"
	xgqA := amgr.AccountArr[Addr]

	xgqA.Balance = 50

	amgr.AccountArr[Addr] = xgqA
}

func addOneBlock(ctrl *Controller) {
	amgr := ctrl.amgr
	pow := ctrl.pow

	Addr1 := "e0fc0a12c2cfd7ee4ae3cc718fb84d152ff2d53057334b2ab9939b5282e3b544"
	Addr2 := "5d3fd2cc8a3c491c87b5e834d5970c99921ee177b2f387d4fe70fffeafa9d16f"
	Addr3 := "4d959a2efa604d7aa8fa56cd1e499ff42c527ea64fa6a65429aca454d156dadd"

	xgqA1 := amgr.AccountArr[Addr1]
	xgqA2 := amgr.AccountArr[Addr2]
	xgqA3 := amgr.AccountArr[Addr3]

	AddrArr := make([]string, 3, 3)
	LastValueArr := make([]float64, 3, 3)
	PayValueArr := make([]float64, 3, 3)
	UtxoHashArr := make([]string, 3, 3)

	AddrArr[0] = xgqA2.PublicKeyHash
	LastValueArr[0] = xgqA2.Balance
	PayValueArr[0] = 0
	UtxoHashArr[0] = xgqA2.UtxoHash

	AddrArr[1] = xgqA1.PublicKeyHash
	LastValueArr[1] = xgqA1.Balance
	PayValueArr[1] = 1
	UtxoHashArr[1] = xgqA1.UtxoHash

	AddrArr[2] = xgqA3.PublicKeyHash
	LastValueArr[2] = xgqA3.Balance
	PayValueArr[2] = 4
	UtxoHashArr[2] = xgqA3.UtxoHash

	Tx, _ := xgqA2.Pay(
		AddrArr,
		LastValueArr,
		PayValueArr,
		UtxoHashArr)

	pow.Add2TrasactionPool(Tx, &amgr.AccountArr)

	ctrl.accountVerifytest()
}

func (ctrl *Controller) accountVerifytest() {
	amgr := ctrl.amgr
	//Addr1 := "e0fc0a12c2cfd7ee4ae3cc718fb84d152ff2d53057334b2ab9939b5282e3b544"
	Addr2 := "5d3fd2cc8a3c491c87b5e834d5970c99921ee177b2f387d4fe70fffeafa9d16f"
	//Addr3 := "4d959a2efa604d7aa8fa56cd1e499ff42c527ea64fa6a65429aca454d156dadd"

	// xgqA1 := amgr.AccountArr[Addr1]
	xgqA2 := amgr.AccountArr[Addr2]
	//xgqA3 := amgr.AccountArr[Addr3]

	TestStr_1 := "aaa01"
	TestStr_2 := "aaa03"

	sg, _ := xgqA2.Sign(TestStr_1)
	res, _ := xgqA2.Verify(TestStr_2, sg)

	if res {
		return
	}
}
