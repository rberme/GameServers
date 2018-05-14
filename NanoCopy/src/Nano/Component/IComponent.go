package Component

// IComponent 组件接口
type IComponent interface {
	Init()
	AfterInit()
	BeforeShutdown()
	Shutdown()
}
