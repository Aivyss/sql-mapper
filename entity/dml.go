package entity

type QueryEntity interface {
	Path() string
	Tag() string
	GetRawSql() string
	GetParts() []*Part
	IsSimpleSql() bool
}

type DMLBody struct {
	AbsFilePath string
	Selects     []*Select
	Inserts     []*Insert
	Deletes     []*Delete
	Updates     []*Update
}

type Case struct {
	CharData string
	Name     string
}

type Part struct {
	Name     string
	CharData string
	Cases    []*Case
}

type Select struct {
	*CommonFields
}

type CommonFields struct {
	FilePath  string
	Name      string
	RawSql    string
	Parts     []*Part
	SimpleSql bool
}

func (q *CommonFields) Path() string {
	return q.FilePath
}
func (q *CommonFields) Tag() string {
	return q.Name
}
func (q *CommonFields) GetRawSql() string {
	return q.RawSql
}
func (q *CommonFields) GetParts() []*Part {
	return q.Parts
}
func (q *CommonFields) IsSimpleSql() bool {
	return q.SimpleSql
}

type Insert struct {
	*CommonFields
}

type Update struct {
	*CommonFields
}

type Delete struct {
	*CommonFields
}
