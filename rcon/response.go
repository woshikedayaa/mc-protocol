package rcon

type PackageType int32
type ResponseBody []byte

const (
	TypeAuthorize    PackageType = 3
	TypeCommand      PackageType = 2
	TypeMultiPackage PackageType = 0
)

type Response struct {
	id   int32
	size int32
	typ  PackageType
	body ResponseBody
}

func (b ResponseBody) String() string {
	return string(b)
}
