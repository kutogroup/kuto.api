package pkg

import "encoding/base64"

const (
	defaultTable = "vE2GNUsDWJIOY4c1wuQprPnqfK_zdCLiFH9texZ5-7SkM8yBmohXAg3RTljaV60b"
)

//KutoBase64 自定义base64
type KutoBase64 struct {
	Coder *base64.Encoding
}

//NewBase64 新建base64
func NewBase64(rand int) *KutoBase64 {
	nt := make([]byte, len(defaultTable))
	copy(nt, defaultTable)

	for i := 0; i < len(nt); i++ {
		p := ((rand + i) % (len(nt) - i)) + i

		if p != i {
			tmp := nt[i]
			nt[i] = nt[p]
			nt[p] = tmp
		}
	}

	newTable := string(nt)
	t := &KutoBase64{
		Coder: base64.NewEncoding(newTable),
	}

	return t
}

//Encode base64加密
func (c *KutoBase64) Encode(src string) string {
	return c.EncodeBytes([]byte(src))
}

//EncodeBytes base64加密
func (c *KutoBase64) EncodeBytes(src []byte) string {
	return c.Coder.EncodeToString(src)
}

//Decode base64解密
func (c *KutoBase64) Decode(src string) (string, error) {
	b, err := c.DecodeToBytes(src)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

//DecodeToBytes base64解密
func (c *KutoBase64) DecodeToBytes(src string) ([]byte, error) {
	return c.Coder.DecodeString(src)
}
