package schema

type NativeTypeSpec struct {
	PrismaName string
	SQLType    string
	GoType     string
	Extension  string
}

var NativeTypes = []NativeTypeSpec{
	{PrismaName: "VarChar", SQLType: "VARCHAR", GoType: "string"},
	{PrismaName: "Char", SQLType: "CHAR", GoType: "string"},
	{PrismaName: "Text", SQLType: "TEXT", GoType: "string"},
	{PrismaName: "Decimal", SQLType: "NUMERIC", GoType: "string"},
	{PrismaName: "Numeric", SQLType: "NUMERIC", GoType: "string"},
	{PrismaName: "Uuid", SQLType: "UUID", GoType: "string"},
	{PrismaName: "Timestamptz", SQLType: "TIMESTAMPTZ", GoType: "time.Time"},
	{PrismaName: "Date", SQLType: "DATE", GoType: "time.Time"},
	{PrismaName: "SmallInt", SQLType: "SMALLINT", GoType: "int32"},
	{PrismaName: "Oid", SQLType: "OID", GoType: "int32"},
	{PrismaName: "Bit", SQLType: "BIT", GoType: "string"},
	{PrismaName: "VarBit", SQLType: "BIT VARYING", GoType: "string"},
	{PrismaName: "Inet", SQLType: "INET", GoType: "string"},
	{PrismaName: "Xml", SQLType: "XML", GoType: "string"},
	{PrismaName: "Real", SQLType: "REAL", GoType: "float64"},
	{PrismaName: "Money", SQLType: "MONEY", GoType: "string"},
	{PrismaName: "Json", SQLType: "JSON", GoType: "json.RawMessage"},
	{PrismaName: "Time", SQLType: "TIME", GoType: "time.Time"},
	{PrismaName: "Timetz", SQLType: "TIMETZ", GoType: "time.Time"},
	
	// Extension-backed types
	{PrismaName: "Citext", SQLType: "citext", GoType: "string", Extension: "citext"},
	{PrismaName: "Ltree", SQLType: "ltree", GoType: "string", Extension: "ltree"},
	{PrismaName: "Hstore", SQLType: "hstore", GoType: "map[string]*string", Extension: "hstore"},
}
