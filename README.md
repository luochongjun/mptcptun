<div>
  <img width="190" height="210" align="left" src="https://raw.githubusercontent.com/v2fly/v2fly-github-io/master/docs/.vuepress/public/readme-logo.png" alt="V2Ray"/>
  <br>
  <h1>Project V</h1>
  <p>Project V is a set of network tools that helps you to build your own computer network. It secures your network connections and thus protects your privacy.</p>
</div>

[![GitHub Test Badge](https://github.com/v2fly/v2ray-core/workflows/Test/badge.svg)](https://github.com/v2fly/v2ray-core/actions)
[![codecov.io](https://codecov.io/gh/v2fly/v2ray-core/branch/master/graph/badge.svg?branch=master)](https://codecov.io/gh/v2fly/v2ray-core?branch=master)
[![codebeat](https://goreportcard.com/badge/github.com/v2fly/v2ray-core)](https://goreportcard.com/report/github.com/v2fly/v2ray-core)
[![Codacy Badge](https://app.codacy.com/project/badge/Grade/e150b7ede2114388921943bf23d95161)](https://www.codacy.com/gh/v2fly/v2ray-core/dashboard?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=v2fly/v2ray-core&amp;utm_campaign=Badge_Grade)
[![Downloads](https://img.shields.io/github/downloads/v2fly/v2ray-core/total.svg)](https://github.com/v2fly/v2ray-core/releases/latest)

## Related Links

- [Documentation](https://www.v2fly.org) and [Newcomer's Instructions](https://www.v2fly.org/guide/start.html)
- Welcome to translate V2Ray documents via [Transifex](https://www.transifex.com/v2fly/public/)

## Packaging Status

> If you are willing to package V2Ray for other distros/platforms, please let us know or seek for help via [GitHub issues](https://github.com/v2fly/v2ray-core/issues).

[![Packaging status](https://repology.org/badge/vertical-allrepos/v2ray.svg)](https://repology.org/project/v2ray/versions)

## License

[The MIT License (MIT)](https://raw.githubusercontent.com/v2fly/v2ray-core/master/LICENSE)

## Credits

This repo relies on the following third-party projects:

- In production:
  - [gorilla/websocket](https://github.com/gorilla/websocket)
  - [lucas-clemente/quic-go](https://github.com/lucas-clemente/quic-go)
  - [pires/go-proxyproto](https://github.com/pires/go-proxyproto)
  - [seiflotfy/cuckoofilter](https://github.com/seiflotfy/cuckoofilter)
  - [google/starlark-go](https://github.com/google/starlark-go)
  - [jhump/protoreflect](https://github.com/jhump/protoreflect)
  - [inetaf/netaddr](https://github.com/inetaf/netaddr)

- For testing only:
  - [miekg/dns](https://github.com/miekg/dns)
  - [h12w/socks](https://github.com/h12w/socks)


## 配置文件转换
配置文件通过protobuf转换，如果新增配置需要安装protoc来进行转换
在下面的链接中下载对应平台的包
https://github.com/protocolbuffers/protobuf/releases/tag/v3.17.1

解压到系统路径中，例如
```
 sudo unzip -d  /usr/local/  protoc-3.17.3-linux-x86_64.zip
 #记得给相关文件赋予权限
 sudo chmod 777 /usr/local/include/google/protobuf/*.proto 
 ```
 安装go相关插件
 ```
go install -v google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
go install -v google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
 ```
 生成pb.go
 ```
 go run ./infra/vprotogen/
 ```