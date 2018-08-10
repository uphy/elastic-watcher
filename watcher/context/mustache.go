package context

type TemplateValue string

func (t TemplateValue) String(ctx ExecutionContext) (string, error) {
	return RenderTemplate(ctx, string(t))
}
