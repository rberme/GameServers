package gameProto

import "protos"

func init() {
	protos.SetMsg(IDConnectSuccessS2C, ConnectSuccessS2C{})
	protos.SetMsg(IDAgainConnectC2S, AgainConnectC2S{})
	protos.SetMsg(IDAgainConnectS2C, AgainConnectS2C{})

	protos.SetMsg(IDOtherLoginS2C, OtherLoginS2C{})
	protos.SetMsg(IDErrorMsgS2C, ErrorMsgS2C{})
	protos.SetMsg(IDUserLoginC2S, UserLoginC2S{})
	protos.SetMsg(IDUserLoginS2C, UserLoginS2C{})
	protos.SetMsg(IDGetUserInfoC2S, GetUserInfoC2S{})
	protos.SetMsg(IDGetUserInfoS2C, GetUserInfoS2C{})
}
