package DBSchema

var TypeIntegers = []string{
	"tinyint",
	"int",
	"smallint",
	"mediumint",
	"int4",
	"int2",
	"year",
	"SMALLINT",
	"INT",
	"INTEGER",
}
var TypeBigIntegers = []string{
	"bigint",
	"NUMBER",
	"int8",
	"LONG",
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
	"JSONB",
	"json",
	"ROWID",
	"UROWID",
	"jsonb",
	"time",
	"geometry",
}
var TypeTimes = []string{
	"datetimeoffset",
	"timestamptz",
	"TIMESTAMP(6) WITH TIME ZONE",
	"TIMESTAMP(6) WITH LOCAL TIME ZONE",
	"TIMESTAMP",
	"TIMESTAMP(6)",
	"datetime",
	"date",
	"DATE",
	"timestamp",
	"TIMESTAMP",
}

var TypeDates = []string{}
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
	"BFILE",
	"BLOB",
}

var TypeGeo = []string{}

func TypeContains(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}
