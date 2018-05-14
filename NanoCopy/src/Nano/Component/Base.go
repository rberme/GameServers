package Component

// Base 组件接口的默认实现
type Base struct{}

// Init 初始化组件
func (c *Base) Init() {}

// AfterInit 在初始化完成之后调用
func (c *Base) AfterInit() {}

// BeforeShutdown 组件关闭前调用
func (c *Base) BeforeShutdown() {}

// Shutdown 关闭组件
func (c *Base) Shutdown() {}
