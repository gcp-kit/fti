# fti

## 使い方

1. ビルド `go build ./cmd/fti/*.go`
2. config.yml を作成
3. 実行 `./main -c "config.yml""`

## config.yml

### targets

コレクション名のディレクトリのリスト

### firestore_project_on_emulator

謎

### firestore_emulator_host

firestore emulator のホスト

ex) `localhost:20048`

## テストデータの作り方

1. FirestoreのCollection名と同じディレクトリを作成する
    1. 大文字小文字の区別があるので注意
2. `config.yaml` の `targets` へディレクトリパスを追加する
    1. この際、順序に気をつけること(ツールは記述されてる順序のとおりに実行する)
3. 1で作成したディレクトリ内に `.js` または `.json` のファイルを作る
    1. ディレクトリ内にはファイルをいくつおいても問題はない。ただし、その際の実行順序は保証されない。

## テストデータ

### 概念

#### version
ファイル形式のバージョン。ツールの互換性のために存在するため、任意のバージョンをつけてはダメ。  
現状は 1.0 のみ。(2021/11/10)

#### ref
参照されるときのID。対象のすべてのデータを通してユニークである必要がある。
データ投入時に自動採番されたIDがこのrefで参照できる。  
refを参照する場合は `$ref_id` のように参照する。

### 形式

#### json

```json
{
  "version": "1.0",
  "items": [
    {
      "ref": "参照されるときのID(重複禁止)",
      "payload": {
        "name": "hoge"
      },
      "SubCollections": {
        "Collection1": [
          {
            "key": "value1"
          },
          {
            "key": "value2"
          }]
      }
    }
  ]
}
```

#### js

最終的に、下記の形式(サンプル)の配列で認識されるものであれば何をしても良い。
内部的にはv8エンジンを搭載しているため、かなり自由なjsが使えると思うが、どこまでの構文に対応しているかは不明。

```js
[
   {
      ref: `参照されるときのID(重複禁止)`,
      payload: {
         parent_id: '$parent_id__1',
         created_at: new Date(), 
         deleted_at: null,
         // 実際に投入されるデータ
      }
   }
]
```

## Q & A

### 日時を入れたい？

#### json

RFC3339(ISO8601)の形式で文字列として入れる

#### js

jsのDateObjectを入れる。
