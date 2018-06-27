package dispatch

// ISession .
type ISession interface {
}

// type ReceiveMsg struct {
// 	Msg []byte
// }

// IDispatch .
type IDispatch interface {
	Process(session ISession, msg []byte)
}
