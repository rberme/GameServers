package Nano

import "Nano/Component"

var (
	comps = make([]regComp, 0)
)

type regComp struct {
	comp Component.IComponent
	opts []Component.Option
}

func startupComponents() {
	// component initialize hooks
	for _, c := range comps {
		c.comp.Init()
	}

	// component after initialize hooks
	for _, c := range comps {
		c.comp.AfterInit()
	}

	// register all components
	for _, c := range comps {
		if err := handler.register(c.comp, c.opts); err != nil {
			logger.Println(err.Error())
		}
	}
	handler.DumpServices()
}

func shutdownComponents() {
	// reverse call `BeforeShutdown` hooks
	length := len(comps)
	for i := length - 1; i >= 0; i-- {
		comps[i].comp.BeforeShutdown()
	}

	// reverse call `Shutdown` hooks
	for i := length - 1; i >= 0; i-- {
		comps[i].comp.Shutdown()
	}
}
