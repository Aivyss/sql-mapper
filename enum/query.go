package enum

type QueryEnum int

const (
	SELECT QueryEnum = iota
	INSERT
	UPDATE
	DELETE
	CREATE
	DROP
)
