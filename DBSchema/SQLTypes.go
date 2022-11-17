package DBSchema

var TypeIntegers = []string{
	"tinyint",
	"int",
	"smallint",
	"mediumint",
	"int8",
	"int4",
	"int2",
	"year",
	"SMALLINT",
	"INT",
	"INTEGER",
	"LONG",
}
var TypeBigIntegers = []string{
	"bigint",
	"NUMBER",
}
var TypeBool = []string{
	"bool",
}
var TypeStrings = []string{
	"char",
	"enum",
	"varchar",
	"nvarchar",
	"longtext",
	"mediumtext",
	"text",
	"ntext",
	"tinytext",
	"geometry",
	"uuid",
	"bpchar",
	"CHARACTER",
	"VARCHAR",
	"VARCHAR2",
	"NVARCHAR2",
	"CLOB",
	"CHAR",
	"NCHAR",
	"NCLOB",
	"JSON",
	"ROWID",
	"UROWID",
	"jsonb",
}
var TypeTimes = []string{
	"time",
	"timestamp",
	"datetimeoffset",
	"timestamptz",
	"TIMESTAMP(6) WITH TIME ZONE",
	"TIMESTAMP(6) WITH LOCAL TIME ZONE",
	"TIMESTAMP",
	"TIMESTAMP(6)",
}
var TypeDates = []string{
	"datetime",
	"date",
	"DATE",
}
var TypeFloat64 = []string{
	"decimal",
	"double",
	"numeric",
	"BINARY_FLOAT",
	"BINARY_DOUBLE",
	"DECIMAL",
}
var TypeFloat32 = []string{
	"float",
	"float8",
	"float4",
	"real",
	"FLOAT",
}
var TypeBinaries = []string{
	"binary",
	"blob",
	"longblob",
	"mediumblob",
	"varbinary",
	"BLOB",
	"BFILE",
}

func TypeContains(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}
