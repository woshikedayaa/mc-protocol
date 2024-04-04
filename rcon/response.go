package rcon

type PackageType int32
type ResponseBody []byte

const (
	TypeAuthorize    PackageType = 3
	TypeCommand      PackageType = 2
	TypeMultiPackage PackageType = 0 // not implemented
)

type Response struct {
	Id   int32
	Size int32
	Typ  PackageType
	Body ResponseBody
}

func (resp *Response) String() string {
	return resp.Body.String()
}

func (b ResponseBody) String() string {
	return string(b)
}
