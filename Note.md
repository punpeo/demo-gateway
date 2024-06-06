//plugin.go
type Plugin interface {
RpcHandler
}

//rpc_handler.go
type RpcHandler interface {
OnReceiveResponse(string, metadata.MD, http.ResponseWriter) string
OnReceiveTrailers(*status.Status, metadata.MD) metadata.MD
OnResolveMethod(*desc.MethodDescriptor)
OnSendHeaders(*http.Request, metadata.MD) metadata.MD
OnReceiveHeaders(metadata.MD) metadata.MD
}
func (h *BasicRpcHandler) OnReceiveResponse(respJson string, _ metadata.MD, _ http.ResponseWriter) string {
return respJson
}

func (h *BasicRpcHandler) OnReceiveTrailers(_ *status.Status, md metadata.MD) metadata.MD {
return md
}

func (h *BasicRpcHandler) OnResolveMethod(_ *desc.MethodDescriptor) {
}

func (h *BasicRpcHandler) OnSendHeaders(_ *http.Request, md metadata.MD) metadata.MD {
return md
}

func (h *BasicRpcHandler) OnReceiveHeaders(md metadata.MD) metadata.MD {
return md
}

//插件实现
//empty.go
type PluginEmpty struct {
gateway.BasicRpcHandler
}
func (p *PluginEmpty) OnReceiveResponse(respJson string, md metadata.MD, _ http.ResponseWriter) string {
}
