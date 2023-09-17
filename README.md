# Intro
- このpackageは現時点で`github.com/jmoiron/sqlx v1.3.5`を元に開発されている。
- このpackageはXML上に作成したSQLを読み込み、`QueryClient`を作り出す。`QueryClient`はXMLで作成したSQLを実行するLayerである。
- このpackageはDML(`INSERT`、`SELECT`、`UPDATE`、`DELETE`)のみ扱う。DDLなどをプログラミング言語で扱う事を作成者はおすすめしない。
- 一つのXMLファイルは一つの`QueryClient`が生成できる。
- 作成者：`Aivyss`
- 権利: `MIT LICENSE`

# XML作成の仕方
- 全てのクエリは`Body`というtagの直下に作成する。
- `Body`の直下には複数の`Select`、`Insert`、`Update`、`Delete`tagを作成できる。
- DML tagは名前(`name` attribute)を重複しないように作成する。
```xml
<?xml version="1.0" encoding="UTF-8" ?>
<Body>
    <Select name="specificUser">
        SELECT
            ACCOUNT_ID,
            USER_NAME,
            USER_ID,
            PASSWORD
        FROM
            ACCOUNT
        WHERE
            USER_ID = :user_id
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
func NewQueryClient(db *sqlx.DB, identifier string, filePath string) (QueryClient, errors.Error)
```

# `endpoint.NewQueryClient`
- XMLを読み込み、クエリを登録してから`QueryClient`を生成する関数
- parameters
  - `db` : DB接続の元になる。
  - `identifier`: `QueryClient`の職別文字。重複はできない。
  - `filePath`: 読み込む対象になるファイルのpath。
- returns
  - `QueryClient`:クエリ操作のinterface
  - `errors.Error`: `QueryClient`の生成が失敗したら、このエラーを通じてエラーの原因が把握できる。

# `QueryClient`のメソッド
```go
type QueryClient interface {
	BeginTx(ctx context.Context) (*sqlx.Tx, errors.Error)
	RollbackTx(ctx context.Context, tx *sqlx.Tx) errors.Error
	CommitTx(ctx context.Context, tx *sqlx.Tx) errors.Error

	GetOne(ctx context.Context, tagName string, dest any, args map[string]any) errors.Error
	GetOneTx(ctx context.Context, tx *sqlx.Tx, tagName string, dest any, args map[string]any) errors.Error
	Get(ctx context.Context, tagName string, dest any, args map[string]any) errors.Error
	GetTx(ctx context.Context, tx *sqlx.Tx, tagName string, dest any, args map[string]any) errors.Error

	InsertOne(ctx context.Context, tagName string, args map[string]any) errors.Error
	InsertOneTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any) errors.Error

	Delete(ctx context.Context, tagName string, args map[string]any) (int64, errors.Error)
	DeleteTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any) (int64, errors.Error)

	Update(ctx context.Context, tagName string, args map[string]any) (int64, errors.Error)
	UpdateTx(ctx context.Context, tx *sqlx.Tx, tagName string, args map[string]any) (int64, errors.Error)

	GetRawQuery(tagName string, enum enum.QueryEnum) (*string, errors.Error)
}
```
- 各メソッドの名前は直ぐ分かるように名付けられている為、説明は簡単にしたい。
- TxがついているメソッドはTransactionと関している。
- `tagName`はXMLのDML tagの固有の名前である。
- `args`のkeyはクエリの`:key`であり`args`のvalueは代入される値である。

# エラーの管理
- このpackageではgolangの基本的な`error`をリターンするわけではなく、固有の`errors.Error`を返す。
- `errors.Error`は既に定義されているものだけリターンするようになっている。
- `errors.Error`の定義に関しては`errors`ディレクトリの`error_code`と`error_identifier`を参照すること。

# 具体的な例
`example`のテストコードで把握できる。

