{{$svrType := .ServiceType}}
{{$svrName := .ServiceName}}

{{- range .MethodSets}}
const Operation{{$svrType}}{{.OriginalName}} = "/{{$svrName}}/{{.OriginalName}}"
{{- end}}

type {{.ServiceType}}HTTPServer interface {
{{- range .MethodSets}}
	{{- if ne .Comment ""}}
	{{.Comment}}
	{{- end}}
	{{.Name}}(context.Context, *{{.Request}}) (*{{.Reply}}, error)
{{- end}}
}

func Register{{.ServiceType}}HTTPRouter(s *gin.Engine, srv {{.ServiceType}}HTTPServer) {
	{{- range .Methods}}
        s.{{.Method}}("{{.Path}}", _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv))
	{{- end}}
}

{{range .Methods}}
func _{{$svrType}}_{{.Name}}{{.Num}}_HTTP_Handler(srv {{$svrType}}HTTPServer) func(ctx *gin.Context)  {
	return func(ctx *gin.Context)  {
		var in {{.Request}}
		{{- if .HasBody}}
		if err := ctx.ShouldBind(&in{{.Body}}); err != nil {
		    ctx.JSON(http.StatusBadRequest, nil)
			return
		}
		{{- if not (eq .Body "")}}
		if err := ctx.ShouldBindQuery(&in); err != nil {
		    ctx.JSON(http.StatusBadRequest, nil)
			return
		}
		{{- end}}
		{{- else}}
		if err := ctx.ShouldBind(&in{{.Body}}); err != nil {
		    ctx.JSON(http.StatusBadRequest, nil)
			return
		}
		{{- end}}

		reply, err := srv.{{.Name}}(ctx.Request.Context(), &in)
		if err != nil {
		    ctx.JSON(http.StatusBadRequest, nil)
			return
		}

		ctx.JSON(http.StatusOK, reply{{.ResponseBody}})
	}
}
{{end}}
{{/*
type {{.ServiceType}}HTTPClient interface {
{{- range .MethodSets}}
	{{.Name}}(ctx context.Context, req *{{.Request}}, opts ...http.CallOption) (rsp *{{.Reply}}, err error)
{{- end}}
}

type {{.ServiceType}}HTTPClientImpl struct{
	cc *http.Client
}

func New{{.ServiceType}}HTTPClient (client *http.Client) {{.ServiceType}}HTTPClient {
	return &{{.ServiceType}}HTTPClientImpl{client}
}

{{range .MethodSets}}
func (c *{{$svrType}}HTTPClientImpl) {{.Name}}(ctx context.Context, in *{{.Request}}, opts ...http.CallOption) (*{{.Reply}}, error) {
	var out {{.Reply}}
	pattern := "{{.Path}}"
	path := binding.EncodeURL(pattern, in, {{not .HasBody}})
	opts = append(opts, http.Operation(Operation{{$svrType}}{{.OriginalName}}))
	opts = append(opts, http.PathTemplate(pattern))
	{{if .HasBody -}}
	err := c.cc.Invoke(ctx, "{{.Method}}", path, in{{.Body}}, &out{{.ResponseBody}}, opts...)
	{{else -}}
	err := c.cc.Invoke(ctx, "{{.Method}}", path, nil, &out{{.ResponseBody}}, opts...)
	{{end -}}
	if err != nil {
		return nil, err
	}
	return &out, err
}
{{end}}
*/}}