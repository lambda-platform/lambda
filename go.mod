module github.com/lambda-platform/lambda

go 1.15



require (
	github.com/99designs/gqlgen v0.13.0 // indirect
	github.com/PaesslerAG/gval v1.1.0 // indirect
	github.com/araddon/dateparse v0.0.0-20210429162001-6b43995a97de // indirect
	github.com/foolin/goview v0.3.0
	github.com/go-redis/redis v6.15.9+incompatible // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/grokify/html-strip-tags-go v0.0.1 // indirect
	github.com/iancoleman/strcase v0.1.3 // indirect
	github.com/jinzhu/gorm v1.9.16 // indirect
	github.com/joho/godotenv v1.3.0 // indirect
	github.com/kelseyhightower/envconfig v1.4.0 // indirect
	github.com/labstack/echo/v4 v4.3.0
	github.com/lambda-platform/agent v0.0.1
	github.com/lambda-platform/chart v0.0.1
	github.com/lambda-platform/crudlogger v0.0.1
	github.com/lambda-platform/dataform v0.0.1 // indirect
	github.com/lambda-platform/datagrid v0.0.1 // indirect
	//	github.com/lambda-platform/dataanalytic v0.0.1
	github.com/lambda-platform/datasource v0.0.1 // indirect
	//	github.com/lambda-platform/arcGIS v0.0.1
	github.com/lambda-platform/graphql v0.0.1
	github.com/lambda-platform/krud v0.0.1

	github.com/lambda-platform/moqup v0.0.1
	github.com/lambda-platform/notify v0.0.1
	github.com/lambda-platform/template v0.0.1
	github.com/otiai10/copy v1.6.0 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tealeg/xlsx v1.0.5 // indirect
	github.com/thedevsaddam/govalidator v1.9.10
	github.com/vektah/gqlparser/v2 v2.1.0 // indirect
)



replace github.com/lambda-platform/template v0.0.1 => ../template

replace github.com/lambda-platform/agent v0.0.1 => ../agent


replace github.com/lambda-platform/krud v0.0.1 => ../krud

replace github.com/lambda-platform/dataform v0.0.1 => ../dataform

replace github.com/lambda-platform/datagrid v0.0.1 => ../datagrid

replace github.com/lambda-platform/datasource v0.0.1 => ../datasource

//PRO MODULES
replace github.com/lambda-platform/crudlogger v0.0.1 => ../crudlogger

replace github.com/lambda-platform/notify v0.0.1 => ../notify

//replace github.com/lambda-platform/arcGIS v0.0.1 => ../arcGIS
replace github.com/lambda-platform/graphql v0.0.1 => ../graphql

//replace github.com/lambda-platform/dataanalytic v0.0.1 => ../dataanalytic
replace github.com/lambda-platform/chart v0.0.1 => ../chart

replace github.com/lambda-platform/moqup v0.0.1 => ../moqup
