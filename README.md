## gPRCを使ったAPIサーバの作成（TODOアプリ）

### 機能
* タスクの追加
* タスクのステータス変更
* タスクの削除

### 非機能要件
* 単体テスト
* シナリオテスト（```go build -buildmode=plugin```がWindows非対応）
* DBを使用しない（sync.Mapを使う）
* エラーハンドリングを行う


### 使用技術など
* connect
gRPC互換のHTTP APIを構築するためのフレームワーク
[github→link](https://github.com/connectrpc)

* evans
gRPCをコマンドラインから確認できるツール
[github→link](https://github.com/ktr0731/evans)

* senarigo
シナリオをテストを行うためのライブラリ
[github→link](https://github.com/zoncoen/scenarigo)