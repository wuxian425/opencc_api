# OpenCC API
將 OpenCC 暴露成 Web API。

## 運行
> [!TIP]
> 以源代碼運行本項目需要有 Golang 環境。

Clone 本項目至本地，並在項目根目錄下運行：
```bash
go run main.go
```
即可以預設 `localhost:1145`、引入 `s2tw` 詞典的狀態啟動 Web 服務。

### Command-line Arguments
```bash
-port   # 設定 Web 服務運行的端口號
-dicts  # 設定服務需要引入的詞典（預設引入 s2tw、導入多個詞典時以 + 號分隔）
```

## 使用
訪問 `localhost(:port)/convert`，提供兩個參數

- `text`    要轉換的文字（必填，否則回傳 400 錯誤）
- `dict`    要使用的詞典（參考 [longbridge/opencc](https://github.com/longbridge/opencc)），必須在運行前用命令行引數定義。

### Example  
運行 `go run main.go -port=11451 -dicts=s2twp+s2tw` 後，   
訪問 `http://localhost:11451/convert?dict=s2twp&text=变量` 將回傳 `變數`。

## 開源
本項目使用 MIT License 開放原始碼。

## 致謝
- [longbridge/opencc](https://github.com/longbridge/opencc) - OpenCC 的 Golang 實現。(Apache 2.0)
- [BYVoid/OpenCC](https://github.com/BYVoid/OpenCC/) - Library for conversion between Traditional and Simplified Chinese。(Apache 2.0)
