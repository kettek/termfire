package messages

type handler struct {
	kind        string
	failureKind string
	once        bool
	fn          func(Message, *MessageFailure)
}

type MessageHandler struct {
	handlers []*handler
}

func (b *MessageHandler) OnMessage(m Message) {
	var calledFuncs []func(Message, *MessageFailure)

	var failureMessage *MessageFailure

	handlers := make([]*handler, 0)
	// If it's a failure message, see if the failure kind exists in our handlers.
	if fail, ok := m.(*MessageFailure); ok {
		failureMessage = fail
		for _, h := range b.handlers {
			if h.failureKind == fail.Command || h.kind == fail.Command {
				calledFuncs = append(calledFuncs, h.fn)
				if !h.once {
					handlers = append(handlers, h)
				}
			} else {
				handlers = append(handlers, h)
			}
		}
	} else {
		for _, h := range b.handlers {
			if h.kind == m.Kind() {
				calledFuncs = append(calledFuncs, h.fn)
				if !h.once {
					handlers = append(handlers, h)
				}
			} else {
				handlers = append(handlers, h)
			}
		}
	}
	b.handlers = handlers

	for _, fn := range calledFuncs {
		fn(m, failureMessage)
	}
}

func (b *MessageHandler) Clear() {
	b.handlers = nil
}

func (b *MessageHandler) On(m Message, f Message, fn func(Message, *MessageFailure)) *handler {
	fKind := ""
	if f != nil {
		fKind = f.Kind()
	}
	handler := handler{m.Kind(), fKind, false, fn}
	b.handlers = append(b.handlers, &handler)
	return &handler
}

func (b *MessageHandler) Once(m Message, f Message, fn func(Message, *MessageFailure)) *handler {
	fKind := ""
	if f != nil {
		fKind = f.Kind()
	}
	handler := handler{m.Kind(), fKind, true, fn}
	b.handlers = append(b.handlers, &handler)
	return &handler
}

func (b *MessageHandler) Off(h *handler) {
	for i, handler := range b.handlers {
		if handler == h {
			b.handlers = append(b.handlers[:i], b.handlers[i+1:]...)
			return
		}
	}
}

func (b *MessageHandler) HasHandlerFor(m Message) bool {
	for _, handler := range b.handlers {
		if handler.kind == m.Kind() {
			return true
		}
	}
	return false
}
