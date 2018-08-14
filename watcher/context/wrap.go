package context

type scopedExecutionContext struct {
	ExecutionContext
	id      string
	payload JSONObject
	vars    JSONObject
}

func newScopedContext(ctx ExecutionContext) (ExecutionContext, error) {
	payload, err := ctx.Payload().Clone()
	if err != nil {
		return nil, err
	}
	vars, err := ctx.Vars().Clone()
	return &scopedExecutionContext{
		ExecutionContext: ctx,
		id:               generateID(),
		payload:          payload,
		vars:             vars,
	}, nil
}
func (s *scopedExecutionContext) ID() string {
	return s.id
}
func (s *scopedExecutionContext) Vars() JSONObject {
	return s.vars
}
func (s *scopedExecutionContext) SetVars(vars JSONObject) {
	s.vars = vars
}
func (s *scopedExecutionContext) Payload() JSONObject {
	return s.payload
}
func (s *scopedExecutionContext) SetPayload(payload JSONObject) {
	s.payload = payload
}
