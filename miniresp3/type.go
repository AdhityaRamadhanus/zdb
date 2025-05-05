package miniresp3

type TypeRESP byte

var (
	RESPMap          TypeRESP = '%'
	RESPArray        TypeRESP = '*'
	RESPBulkString   TypeRESP = '$'
	RESPSimpleString TypeRESP = '+'
	RESPNumber       TypeRESP = ':'
	RESPDoubles      TypeRESP = ','
	RESPNull         TypeRESP = '_'
	RESPBulkError    TypeRESP = '!'
	RESPSimpleError  TypeRESP = '-'
)
