package Chat

const ( //私有数据处理
	// PrivateMsgStart 私有数据开始
	PrivateMsgStart = iota + 1000000
	// PrivateTestMsg 私有消息测试
	PrivateTestMsg
	// PrivateTestMsgRET 私有消息测试 返回
	PrivateTestMsgRET
)

const ( //mainloop处理
	// PublicMsgStart 共享数据开始
	PublicMsgStart = iota + 1500000
	// PublicTestMsg 共享数据测试消息 1
	PublicTestMsg
	// PublicTestMsgRET 共享数据测试消息 返回 2
	PublicTestMsgRET
	// PublicLogin 登录 3
	PublicLogin
	// PublicLoginRET 登录返回 4
	PublicLoginRET
	// PublicLogout 登出 5
	PublicLogout
	// PublicLogoutRET 登出返回 6
	PublicLogoutRET
	// PublicChating 普通聊天 7
	PublicChating
	// PublicChatingRET 普通聊天返回 8
	PublicChatingRET
)
