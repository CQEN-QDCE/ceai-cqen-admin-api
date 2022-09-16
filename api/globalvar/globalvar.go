package globalvar

var ProdEnv bool

func SetProdEnv(b bool) {
	ProdEnv = b
}

func IsProdEnv() bool {
	return ProdEnv == true
}
