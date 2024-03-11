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

func (resp *Response) ID() int32 {
	return resp.id
}

func (resp *Response) Size() int32 {
	return resp.size
}

func (resp *Response) Type() PackageType {
	return resp.typ
}

func (resp *Response) Body() ResponseBody {
	return resp.body
}

func (resp *Response) String() string {
	return resp.body.String()
}

func (b ResponseBody) String() string {
	return string(b)
}
