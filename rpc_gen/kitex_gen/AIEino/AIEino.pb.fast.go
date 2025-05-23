// Code generated by Fastpb v0.0.2. DO NOT EDIT.

package AIEino

import (
	fmt "fmt"
	fastpb "github.com/cloudwego/fastpb"
)

var (
	_ = fmt.Errorf
	_ = fastpb.Skip
)

func (x *AIWithOrdersReq) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_AIWithOrdersReq[number], err)
}

func (x *AIWithOrdersReq) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserInput, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *AIWithOrdersReq) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *AIWithOrdersResp) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_AIWithOrdersResp[number], err)
}

func (x *AIWithOrdersResp) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	var v string
	v, offset, err = fastpb.ReadString(buf, _type)
	if err != nil {
		return offset, err
	}
	x.Orders = append(x.Orders, v)
	return offset, err
}

func (x *AIWithPreCheckoutReq) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	case 2:
		offset, err = x.fastReadField2(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_AIWithPreCheckoutReq[number], err)
}

func (x *AIWithPreCheckoutReq) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	x.UserInput, offset, err = fastpb.ReadString(buf, _type)
	return offset, err
}

func (x *AIWithPreCheckoutReq) fastReadField2(buf []byte, _type int8) (offset int, err error) {
	x.UserId, offset, err = fastpb.ReadUint32(buf, _type)
	return offset, err
}

func (x *AIWithPreCheckoutResp) FastRead(buf []byte, _type int8, number int32) (offset int, err error) {
	switch number {
	case 1:
		offset, err = x.fastReadField1(buf, _type)
		if err != nil {
			goto ReadFieldError
		}
	default:
		offset, err = fastpb.Skip(buf, _type, number)
		if err != nil {
			goto SkipFieldError
		}
	}
	return offset, nil
SkipFieldError:
	return offset, fmt.Errorf("%T cannot parse invalid wire-format data, error: %s", x, err)
ReadFieldError:
	return offset, fmt.Errorf("%T read field %d '%s' error: %s", x, number, fieldIDToName_AIWithPreCheckoutResp[number], err)
}

func (x *AIWithPreCheckoutResp) fastReadField1(buf []byte, _type int8) (offset int, err error) {
	var v string
	v, offset, err = fastpb.ReadString(buf, _type)
	if err != nil {
		return offset, err
	}
	x.PreCheckoutRes = append(x.PreCheckoutRes, v)
	return offset, err
}

func (x *AIWithOrdersReq) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *AIWithOrdersReq) fastWriteField1(buf []byte) (offset int) {
	if x.UserInput == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 1, x.GetUserInput())
	return offset
}

func (x *AIWithOrdersReq) fastWriteField2(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 2, x.GetUserId())
	return offset
}

func (x *AIWithOrdersResp) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	return offset
}

func (x *AIWithOrdersResp) fastWriteField1(buf []byte) (offset int) {
	if len(x.Orders) == 0 {
		return offset
	}
	for i := range x.GetOrders() {
		offset += fastpb.WriteString(buf[offset:], 1, x.GetOrders()[i])
	}
	return offset
}

func (x *AIWithPreCheckoutReq) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	offset += x.fastWriteField2(buf[offset:])
	return offset
}

func (x *AIWithPreCheckoutReq) fastWriteField1(buf []byte) (offset int) {
	if x.UserInput == "" {
		return offset
	}
	offset += fastpb.WriteString(buf[offset:], 1, x.GetUserInput())
	return offset
}

func (x *AIWithPreCheckoutReq) fastWriteField2(buf []byte) (offset int) {
	if x.UserId == 0 {
		return offset
	}
	offset += fastpb.WriteUint32(buf[offset:], 2, x.GetUserId())
	return offset
}

func (x *AIWithPreCheckoutResp) FastWrite(buf []byte) (offset int) {
	if x == nil {
		return offset
	}
	offset += x.fastWriteField1(buf[offset:])
	return offset
}

func (x *AIWithPreCheckoutResp) fastWriteField1(buf []byte) (offset int) {
	if len(x.PreCheckoutRes) == 0 {
		return offset
	}
	for i := range x.GetPreCheckoutRes() {
		offset += fastpb.WriteString(buf[offset:], 1, x.GetPreCheckoutRes()[i])
	}
	return offset
}

func (x *AIWithOrdersReq) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *AIWithOrdersReq) sizeField1() (n int) {
	if x.UserInput == "" {
		return n
	}
	n += fastpb.SizeString(1, x.GetUserInput())
	return n
}

func (x *AIWithOrdersReq) sizeField2() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeUint32(2, x.GetUserId())
	return n
}

func (x *AIWithOrdersResp) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	return n
}

func (x *AIWithOrdersResp) sizeField1() (n int) {
	if len(x.Orders) == 0 {
		return n
	}
	for i := range x.GetOrders() {
		n += fastpb.SizeString(1, x.GetOrders()[i])
	}
	return n
}

func (x *AIWithPreCheckoutReq) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	n += x.sizeField2()
	return n
}

func (x *AIWithPreCheckoutReq) sizeField1() (n int) {
	if x.UserInput == "" {
		return n
	}
	n += fastpb.SizeString(1, x.GetUserInput())
	return n
}

func (x *AIWithPreCheckoutReq) sizeField2() (n int) {
	if x.UserId == 0 {
		return n
	}
	n += fastpb.SizeUint32(2, x.GetUserId())
	return n
}

func (x *AIWithPreCheckoutResp) Size() (n int) {
	if x == nil {
		return n
	}
	n += x.sizeField1()
	return n
}

func (x *AIWithPreCheckoutResp) sizeField1() (n int) {
	if len(x.PreCheckoutRes) == 0 {
		return n
	}
	for i := range x.GetPreCheckoutRes() {
		n += fastpb.SizeString(1, x.GetPreCheckoutRes()[i])
	}
	return n
}

var fieldIDToName_AIWithOrdersReq = map[int32]string{
	1: "UserInput",
	2: "UserId",
}

var fieldIDToName_AIWithOrdersResp = map[int32]string{
	1: "Orders",
}

var fieldIDToName_AIWithPreCheckoutReq = map[int32]string{
	1: "UserInput",
	2: "UserId",
}

var fieldIDToName_AIWithPreCheckoutResp = map[int32]string{
	1: "PreCheckoutRes",
}
