# Introduction
- This package is currently developed based on `github.com/jmoiron/sqlx v1.3.5`.
- This package reads SQL created in XML and generates a `QueryClient`, which is a layer for executing SQL created in XML.
- This package handles only DML (`INSERT`, `SELECT`, `UPDATE`, `DELETE`) operations. The creator does not recommend handling DDL and other operations in this package.
- One XML file can generate one QueryClient.
- Creator: `Aivyss`
- License: `MIT LICENSE`

# Interface `ApplicationContext`
```go
type ApplicationContext interface {
  GetQueryClient(identifier string) (QueryClient, errors.Error)
  GetReadOnlyQueryClient(identifier string) (ReadOnlyQueryClient, errors.Error)
  RegisterQueryClient(client QueryClient) errors.Error
  GetDBs() *entity.DbSet
  GetDB(readDB bool) *sqlx.DB
}
```
- This package manages the overall application queries through the `ApplicationContext` interface.

# Creation: `ApplicationContext`
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

# Write DML Query
- All queries are created under the Body tag.
- Under Body, you can create multiple Select, Insert, Update, and Delete tags.
- DML tags should have unique names (name attribute).
- Creating queries:
  - Create queries directly under DML tags.
  - Use the `<Part>` tag to separate queries.
  - Use `<Case>` tags inside `<Part>` to create dynamic queries.
  - The name attribute of the `<Part>` tag is not mandatory, but it becomes mandatory for `<Part>` tags used with `<Case>` tags.
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

# Generating `QueryClient` Interface
- Each `QueryClient` is mapped to a unique identifier string. Therefore, duplicate identifiers are not allowed.
- To create a `QueryClient`, you can use the `context.NewQueryClient` function.
```go
func NewQueryClient(identifier string, filePath string) (QueryClient, errors.Error)
func NewReadOnlyQueryClient(identifier string, filePath string) (ReadOnlyQueryClient, errors.Error)
```
- While these methods can be used to create a `QueryClient`, it is preferable to use Context configuration.

# Methods of `QueryClient`
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
- The names of each method are self-explanatory, so explanations are kept brief.
- Methods with `Tx` deal with transactions.
- `tagName` is the unique name of the DML tag in XML.
- The key in `args` is the `:key` in the query, and the value in `args` is the value to be assigned.
- Conditions are used for dynamic queries when querying with `QueryClient` or `ReadOnlyQueryClient`.

# `Condition` and `PredicateConditions`
```go
type Condition struct {
	PartName string
	CaseName string
}
```
- Since the condition clauses are in the form `<Part><Case></Case></Part>` in this package, this structure is used to determine conditions.

```go
type PredicateConditions func() []*Condition
```
- Lambdas are used to determine dynamic query conditions when executing queries with `QueryClient` or `ReadOnlyQueryClient`.

# Error Handling
- In this package, instead of returning the basic error in Go, it returns a specific `errors.Error`.
- Only pre-defined `errors.Error` is returned.
- For the definition of `errors.Error`, refer to the `error_code` and `error_identifier` in the `errors` directory.

# Specific Examples
- You can understand more through the `test` test code.