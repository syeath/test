"# Gui-Tools" 

收集的思路
1. 从txt中提取全部的url以及ip 
2. 针对域名(子域名爆破)
3. 针对ip(端口扫描)
4. 指纹识别

额外添加的功能
1. 目录扫描 
2. poc扫描

导入依赖
```go
go get github.com/chainreactors/fingers@master
go get -u github.com/spf13/cobra@latest
```


```go
go get -u github.com/gocarina/gocsv
```