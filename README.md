# Intro
- このpackageは現時点で`github.com/jmoiron/sqlx v1.3.5`を元に開発されている。
- このpackageはXML上に作成したSQLを読み込み、`QueryClient`を作り出す。`QueryClient`はXMLで作成したSQLを実行するLayerである。
- このpackageはDML(`INSERT`、`SELECT`、`UPDATE`、`DELETE`)のみ扱う。DDLなどをプログラミング言語で扱う事を作成者はおすすめしない。
- 一つのXMLファイルは一つの`QueryClient`が生成できる。
- 作成者：`Aivyss`
- 権利: `MIT LICENSE`

# interface `ApplicationContext`

```go
type ApplicationContext interface {
  GetQueryClient(identifier string) (QueryClient, errors.Error)
  GetReadOnlyQueryClient(identifier string) (ReadOnlyQueryClient, errors.Error)
  RegisterQueryClient(client QueryClient) errors.Error
  GetDBs() *entity.DbSet
  GetDB(readDB bool) *sqlx.DB
}
```
- このpackageは`ApplicationContext`というinterfaceにより、アプリの全体のクエリを管理している。


# XML作成の仕方：`ApplicationContext`
```go
// DBの設定1(簡単)
context.Bootstrap(db) // initiator構造体の値をリターン
// DBの設定2(write, read)
context.BooststrapDual(writeDb, readDb) // initiator構造体の値をリターン
```
```go
// Contextを生成
context.Bootstrap(db).InitByXml("./setting/settings.xml")
```

```xml
<!-- settings.xml -->
<?xml version="1.0" encoding="UTF-8" ?>
<Context>
  <QueryClients>
    <QueryClient identifier="identifier4-1" readOnly="false" filePath="./mapper/sql1.xml"/>
    <QueryClient identifier="identifier4-2" readOnly="false" filePath="./mapper/sql2.xml"/>
    <QueryClient identifier="identifier4-3" readOnly="false" filePath="./mapper/sql3.xml"/>
    <QueryClient identifier="identifier1" readOnly="false" filePath="./mapper/sql1.xml"/>
    <QueryClient identifier="identifier2" readOnly="false" filePath="./mapper/sql2.xml"/>
    <QueryClient identifier="identifier3" readOnly="false" filePath="./mapper/sql3.xml"/>
  </QueryClients>
</Context>
```


# XML作成の仕方：DML Query Part
- 全てのクエリは`Body`というtagの直下に作成する。
- `Body`の直下には複数の`Select`、`Insert`、`Update`、`Delete`tagを作成できる。
- DML tagは名前(`name` attribute)を重複しないように作成する。
- クエリ作成方法
  - DML tagの直下にクエリを作成
  - `<Part>`tagを利用してクエリを分けて作成する
  - `<Part>`tagの内部に`<Case>`tagを入れて動的クエリを作成する
  - `<Part>`tagの`name`attributeは必須ではないが、`<Case>`tagが使われる`<Part>`tagは`name`が必須になる。
```xml
<?xml version="1.0" encoding="UTF-8" ?>
<Body>
    <Select name="specificUser" list="false">
        <Part>
            SELECT
        </Part>
        <Part name="condition1">
            <Case name="case1">
                USER_NAME,
            </Case>
            <Case name="case2">
      
            </Case>
        </Part>
        <Part>
            USER_ID,
        </Part>
        <Part name="condition2">
            <Case name="case3">
                PASSWORD,
            </Case>
            <Case name="case4">
      
            </Case>
        </Part>
        <Part>
            ACCOUNT_ID
            FROM
            ACCOUNT
            WHERE
            USER_ID = :user_id
        </Part>
    </Select>

    <Select name="allUsers">
        SELECT
            ACCOUNT_ID,
            USER_NAME,
            USER_ID,
            PASSWORD
        FROM
            ACCOUNT
    </Select>

    <Insert name="saveOneUser">
        INSERT INTO ACCOUNT (
            USER_NAME,
            USER_ID,
            PASSWORD
        ) VALUES (
            :user_name,
            :user_id,
            :password
        )
    </Insert>

    <Update name="updateUserNameForOneUser">
        UPDATE ACCOUNT SET
            USER_NAME = :user_name
        WHERE
            USER_ID = :user_id
    </Update>

    <Delete name="deleteOneUser">
        DELETE FROM ACCOUNT
        WHERE
            USER_ID = :user_id
    </Delete>

    <Delete name="fullDelete">
        DELETE FROM ACCOUNT
        WHERE 1=1
    </Delete>
</Body>
```

# interface: `QueryClient`の生成
- 各`QueryClient`は固有の`identifier`文字列にmappingされる。すなわち、重複した`identifier`はできない。
- `QueryClient`を生成する為には`endpoint.NewQueryClient`関数を使う。
```go
func NewQueryClient(identifier string, filePath string) (QueryClient, errors.Error)
func NewReadOnlyQueryClient(identifier string, filePath string) (ReadOnlyQueryClient, errors.Error)
```
- このメソッドを利用することでQueryClientは生成できるが、Context設定を利用するのが好ましい。



# `QueryClient`のメソッド
```go
type QueryClient interface {
  ReadOnlyQueryClient
  InsertOne(ctx context.Context, tagName string, args map[string]any, conditions ...entity.PredicateConditions) errors.Error
  InsertOneTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any, conditions ...entity.PredicateConditions) errors.Error
  
  Delete(ctx context.Context, tagName string, args map[string]any, conditions ...entity.PredicateConditions) (int64, errors.Error)
  DeleteTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any, conditions ...entity.PredicateConditions) (int64, errors.Error)
  
  Update(ctx context.Context, tagName string, args map[string]any, conditions ...entity.PredicateConditions) (int64, errors.Error)
  UpdateTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any, conditions ...entity.PredicateConditions) (int64, errors.Error)
}

type ReadOnlyQueryClient interface {
  BeginTx(ctx context.Context) (*sqlx.Tx, errors.Error)
  RollbackTx(ctx context.Context, tx *sqlx.Tx) errors.Error
  CommitTx(ctx context.Context, tx *sqlx.Tx) errors.Error
  
  GetOne(ctx context.Context, tagName string, dest any, args map[string]any, conditions ...entity.PredicateConditions) errors.Error
  GetOneTx(ctx context.Context, tx *sqlx.Tx, tagName string, dest any, args map[string]any, conditions ...entity.PredicateConditions) errors.Error
  Get(ctx context.Context, tagName string, dest any, args map[string]any, conditions ...entity.PredicateConditions) errors.Error
  GetTx(ctx context.Context, tx *sqlx.Tx, tagName string, dest any, args map[string]any, conditions ...entity.PredicateConditions) errors.Error
  
  Id() string
  ReadOnly() bool
}
```
- 各メソッドの名前は直ぐ分かるように名付けられている為、説明は簡単にしたい。
- TxがついているメソッドはTransactionと関している。
- `tagName`はXMLのDML tagの固有の名前である。
- `args`のkeyはクエリの`:key`であり`args`のvalueは代入される値である。
- conditionsは動的クエリに使われる条件関数

# `Condition`と`PredicateConditions`
```go
type Condition struct {
	PartName string
	CaseName string
}
```
- このpackageは条件文が`<Part><Case></Case></Part>`の形になっているため、この構造体で条件を決める。

```go
type PredicateConditions func() []*Condition
```
- `QueryClient`や`ReadOnlyQueryClient`でクエリを実施する際、動的クエリの条件を決めるlambda

# エラーの管理
- このpackageではgolangの基本的な`error`をリターンするわけではなく、固有の`errors.Error`を返す。
- `errors.Error`は既に定義されているものだけリターンするようになっている。
- `errors.Error`の定義に関しては`errors`ディレクトリの`error_code`と`error_identifier`を参照すること。

# 具体的な例
`test`のテストコードで把握できる。

