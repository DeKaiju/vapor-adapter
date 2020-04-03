# vapor-adapter
An adapter to interact with vapor implemented in GO

## Usage
```
import "github.com/DeKaiju/vapor-adapter"
```

## Warning
server.go中大部分函数需要传入的是用户地址在全节点中的对应的accountId,后端需要处理好用户地址与accountId的对应关系  
CreateAccount中传入的accountAlias不能重名，建议对用户传入的名称加上时间戳或者其他随机字符串作为accountAlias  
BuildTransaction构建新的交易会找零到新的地址，用户地址在链上的金额并不等于用户整个钱包的余额  
