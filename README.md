go-zero两种方式使用gRPC网关
protoDescriptor 和 grpcReflection

环境:etcd \go1.18

go mod init demo-gateway

study.pb生成
```bash
protoc --descriptor_set_out=etc/study.pb rpc/study/study.proto
```

### 根据proto文件生成相关代码

```bash
goctl rpc protoc rpc/study/study.proto --go_out=rpc/study --go-grpc_out=rpc/study --zrpc_out=rpc/study --style=goZero -m
```

### 根据api文件生成相关代码

```bash
goctl api go -api api/study/study.api -dir api/study --style=goZero
```
###go.mod 限制go-zero版本
require (
github.com/zeromicro/go-zero v1.5.3
)