# ToDoアプリ - クリーンアーキテクチャ採用

## 概要
本プロジェクトは、クリーンアーキテクチャを採用して作成したToDoアプリです。  
Go言語と `echo` フレームワークを使用して開発しましたが、アーキテクチャ設計により `gin` など他のフレームワークへの切り替えが容易に行える構造になっています。  
また、OpenAPI を利用して API ドキュメントを生成し、Swagger UI でAPIを確認できます。

---

## 機能
- ユーザー認証 (JWT + CSRF)
- ToDoの作成、取得、更新、削除
- Swagger UI による API ドキュメントの確認

---

## 使用技術
- **言語**: Go
- **フレームワーク**: Echo
- **データベース**: MySQL
- **APIドキュメント**: OpenAPI (Swagger UI)
- **ツール**:
  - Docker (開発環境)
  - oapi-codegen (OpenAPIからコード生成)

---


以下のコマンドを実行してopenapi.yamlからGoのコードを作成
```sh
oapi-codegen --config=./adapter/controller/echo/config.yaml ./api/openapi.yaml
```

## セットアップ

### 環境構築手順
1. リポジトリをクローンします。

```sh
git clone [git@github.com:kazukisasajima/go-todo-app-clean-arch.git](git@github.com:kazukisasajima/go-todo-app-clean-arch.git)
cd go-todo-app-clean-arch
```

2. MySQLの起動
```sh
pushd ./build/docker && docker-compose up -d mysql && popd
```

3. Swagger UIの起動
```sh
pushd ./build/docker && docker-compose up -d swagger-ui && popd
```

4.SQLクライアントの使用 (任意)
MySQLコンテナに接続して直接クエリを実行できます。
```sh
pushd ./build/docker && docker-compose run mysql-cli && popd
```

5. OpenAPIコード生成
以下のコマンドで openapi.yaml からGoのコードを生成します。
```sh
oapi-codegen --config=./adapter/controller/echo/config.yaml ./api/openapi.yaml
```

## サーバーの起動方法
```sh
go run ./cmd/server/main.go
```

フロントエンド([react-todo-v2](https://github.com/kazukisasajima/react-todo-v2))を起動してから以下URLにアクセス  
[http://localhost:3000](http://localhost:3000/)

## API ドキュメント
Swagger UIを使用してAPI仕様を確認できます。
Swagger UI URL: [http://localhost:8080/swagger](http://localhost:8080/swagger)
