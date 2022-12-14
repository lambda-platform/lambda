package generator

import (
	"encoding/json"
	"fmt"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/lambda-platform/lambda/generator/models"
	"github.com/lambda-platform/lambda/generator/utils"
	"github.com/lambda-platform/lambda/graphql/plugin/resolvergen"
	lambdaModels "github.com/lambda-platform/lambda/models"
	"github.com/otiai10/copy"
	"os"
	"strings"
)

func Generate() {
	cfg, err := config.LoadConfig("lambda/graph/gqlgen.yml")

	err = api.Generate(cfg,
		api.AddPlugin(resolvergen.New("lambda/lambda/graph/resolvers")),
	)

	if err != nil {

		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(3)

		/*
			Replace Project Path to client path
		*/

		//gqlGeneratedFile, _ := os.ReadFile( "lambda/graph/generated/generated.go")
		////gqlGeneratedFileContent := strings.ReplaceAll(string(gqlGeneratedFile), "/", "")
		//utils.WriteFile(gqlGeneratedFileContent, projectPath + "/lambda/graph/generated/generated.go")
		//
		//modelGeneratedFile, _ := os.ReadFile(projectPath + "/lambda/graph/model/models_gen.go")
		//modelGeneratedFileContent := strings.ReplaceAll(string(modelGeneratedFile), projectPath+"/", "")
		//utils.WriteFile(modelGeneratedFileContent, projectPath + "/lambda/graph/model/models_gen.go")
		//
		//resolverGeneratedFile, _ := os.ReadFile(projectPath + "/lambda/graph/schemas.resolvers.go")
		//resolverGeneratedFileContent := strings.ReplaceAll(string(resolverGeneratedFile), projectPath+"/", "")
		//utils.WriteFile(resolverGeneratedFileContent, projectPath + "/lambda/graph/schemas.resolvers.go")
		//
		//mutationsGeneratedFile, fileError := os.ReadFile(projectPath + "/lambda/graph/mutations.resolvers.go")
		//if(fileError == nil){
		//	mutationsFileContent := strings.ReplaceAll(string(mutationsGeneratedFile), projectPath+"/", "")
		//	utils.WriteFile(mutationsFileContent, projectPath + "/lambda/graph/mutations.resolvers.go")
		//}

	}

}
func GenerateSchema(graphqlchemas []models.ProjectSchemas, dbSchema lambdaModels.DBSCHEMA) {

	GqlTables := []models.GqlTable{}

	for _, preTable := range graphqlchemas {
		GqlTable := models.GqlTable{}
		json.Unmarshal([]byte(preTable.Schema), &GqlTable)
		GqlTables = append(GqlTables, GqlTable)
	}

	resolverTmplate := `package resolvers

import (
	"context"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/graphql/gql"
	"lambda/lambda/graph/model"
	"lambda/lambda/graph/models"
	%s
)

func %s(ctx context.Context, sorts []*model.Sort, groupFilters []*model.GroupFilter, filters []*model.Filter%s, limit *int, offset *int) ([]*models.%s, error) {
	%s
	result := []*models.%s{}
	requestColumns, %s:= gql.GetColumns(ctx, "%s")
	requestColumns = append(requestColumns, "%s")
	requestColumns = append(requestColumns, []string{%s}...)
	requestColumns = gql.RemoveDuplicate(requestColumns)
	query := DB.DB.Select(requestColumns)
	columns := %sColumns()
	query, errorFilter := gql.Filter(filters, query,columns)
	query, errorGroupFilter := gql.GroupFilter(groupFilters, query,columns)
	if(errorFilter != nil){
		return result, errorFilter 
	}
	if(errorGroupFilter != nil){
		return result, errorGroupFilter 
	}
	query, errorOrder := gql.Order(sorts, query, columns)
	if(errorOrder != nil){
		return result, errorOrder
	}
	if(limit != nil){
		if(*limit >= 1){
			query = query.Limit(*limit)
		}
	}
	if(offset != nil){
		if(*offset >= 1){
			query = query.Offset(*offset)
		}
	}

	err := query.Find(&result).Error

	%s
}

func %sColumns() []string {
	return []string{%s}
}
`
	subTemp := `var %sSubs = map[string]model.SubTable{
	%s
}
func %sSub(table string) model.SubTable {
	return %sSubs[table]
}
`
	setSubTemplate := `
func Set%sSubs(ctx context.Context, parents []*models.%s, subs[]gql.Sub, subSorts []*model.SubSort, subFilters []*model.SubFilter) ([]*models.%s, error) {
	parentIds := ""
	for _, parent := range parents{
		if(parentIds == ""){
			parentIds = %s
		} else {
			parentIds = parentIds + ","+%s
		}
	}
	for _, sub := range subs {
		%s
	}

	return  parents, nil
}`
	subSetTemp := `
    if (sub.Table == "%s"){
			subItem := %sSub("%s")
			sorts := []*model.Sort{}
			filters := []*model.Filter{}
			for _, sort := range subSorts {
				if sort.Table == "%s" {
					newSort := model.Sort{
						Column: sort.Column,
						Order: sort.Order,
					}
					sorts = append(sorts, &newSort)
				}
			}
			for _, filter := range subFilters {
				if filter.Table == "%s" {
					newFilter := model.Filter{
						Column: filter.Column,
						Condition: filter.Condition,
						Value: filter.Value,
					}
					filters = append(filters, &newFilter)
				}
			}
			parentFilter := model.Filter{}

			parentFilter.Condition = "whereIn"
			parentFilter.Column = subItem.ConnectionField
			parentFilter.Value = parentIds
			filters = append(filters, &parentFilter)

			sub.Columns = append(sub.Columns, subItem.ConnectionField)
			SubItems, err  := %s
			if err != nil {
				return parents, err
			}
			for _, SubItemData := range SubItems{
				for i, _ := range parents{
					%s{
						parents[i].%s = append(parents[i].%s, SubItemData)
					}
				}
			}
		}`

	paginationTmplate := `package resolvers

import (
	"context"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/utils"
	"github.com/lambda-platform/lambda/graphql/gql"
	"lambda/lambda/graph/model"
	"lambda/lambda/graph/models"
)

func Paginate(ctx context.Context, sorts []*model.Sort, groupFilters []*model.GroupFilter, filters []*model.Filter, subSorts []*model.SubSort, subFilters []*model.SubFilter, page int, size int) (*model.Paginate, error) {

	target, _, err := gql.GetPaginationTargetAndColumns(ctx)
	requestColumns, %s := gql.GetColumns(ctx, target)

	Paginate := model.Paginate{
		Page: 0,
		Total:0,
		LastPage:0,
	}
	if(err != nil){
		return &Paginate, err
	}
	query := DB.DB
`
	QueryContent := "type Query {\n"
	Pagination := "\ntype paginate  {\n    page: Int!\n    total: Int!\n    last_page: Int!\n"

	paginationSub := "_"

	for _, table := range GqlTables {
		modelAlias := GetModelAlias(table.Table)

		Identity := GetModelAlias(table.Identity)

		subTables := []string{}
		subTablesMap := ""
		subSetTemps := ""
		colunms := TableMetaColumnsWithMeta(dbSchema.TableMeta[table.Table], table.Table, table.HiddenColumns)
		columnsString := ""

		IdentityColumn := map[string]string{
			"column": "",
		}
		for _, column := range colunms {
			if columnsString == "" {
				columnsString = columnsString + "\"" + column["column"] + "\""
			} else {
				columnsString = columnsString + ", " + "\"" + column["column"] + "\""
			}
			if column["column"] == table.Identity {
				IdentityColumn = column
			}
		}
		IdentityForSub := "parent." + Identity
		ParemtIDNullCheck := ""

		if IdentityColumn["column"] != "" {

			if IdentityColumn["nullable"] != "YES" {

			} else {
				IdentityForSub = "*" + IdentityForSub
				ParemtIDNullCheck = "*"
			}
		}
		if len(table.Subs) >= 1 {
			paginationSub = "subs"
		}
		for _, sub := range table.Subs {
			SubConnectionNullCheck := ""
			subAlias := sub.Table
			subTables = append(subTables, sub.Table)
			subTablesMap = subTablesMap + fmt.Sprintf(`"%s": model.SubTable{
	Table:"%s",
	ParentIdentity:"%s",
	ConnectionField:"%s",
},
`,
				sub.Table,
				sub.Table,
				sub.ParentIdentity,
				sub.ConnectionField,
			)
			subCaller := GetModelAlias(subAlias) + "(ctx, sorts, []*model.GroupFilter{}, filters, nil, nil)"
			subHasSub := false
			for _, tableCheck := range GqlTables {
				if tableCheck.Table == sub.Table {
					if len(tableCheck.Subs) >= 1 {
						subHasSub = true
					}
				}
			}
			if subHasSub {
				subCaller = GetModelAlias(subAlias) + "(ctx, sorts, []*model.GroupFilter{}, filters, subSorts, subFilters, nil, nil)"
			}

			subColunms := TableMetaColumnsWithMeta(dbSchema.TableMeta[sub.Table], sub.Table, []string{})

			SubConnectionColumn := map[string]string{
				"column": "",
			}
			for _, subColunm := range subColunms {

				if subColunm["column"] == sub.ConnectionField {
					SubConnectionColumn = subColunm
				}
			}
			if SubConnectionColumn["column"] != "" {
				if SubConnectionColumn["column"] == "surgalt_id" && sub.Table == "surgalt_elselt" {

				}
				if SubConnectionColumn["nullable"] != "YES" {

				} else {

					SubConnectionNullCheck = "*"
				}
			}

			subSetTemps = subSetTemps + fmt.Sprintf(subSetTemp,
				sub.Table,
				modelAlias,
				sub.Table,
				sub.Table,
				sub.Table,
				subCaller,
				`if fmt.Sprintf("%v", `+ParemtIDNullCheck+`parents[i].`+Identity+`) == fmt.Sprintf("%v", `+SubConnectionNullCheck+`SubItemData.`+GetModelAlias(sub.ConnectionField)+`)`,
				//ParemtIDNullCheck,
				//Identity,
				//SubConnectionNullCheck,
				//dbToStruct.FmtFieldName(dbToStruct.StringifyFirstChar(sub.ConnectionField)),
				GetModelAlias(subAlias),
				GetModelAlias(subAlias),
			)

		}
		parentConnectsions := ""
		for _, tableCheck := range GqlTables {
			for _, sub := range tableCheck.Subs {
				if table.Table == sub.Table {
					if parentConnectsions == "" {
						parentConnectsions = "\"" + sub.ConnectionField + "\""
					} else {
						parentConnectsions = parentConnectsions + ",\"" + sub.ConnectionField + "\""
					}
				}
			}
		}

		/////////dbSchema lambdaModels.DBSCHEMA
		structStr := TableMetaToStruct(dbSchema.TableMeta[table.Table], table.Table, table.HiddenColumns, "models", subTables)

		schema := TableMetaToGraphql(dbSchema.TableMeta[table.Table], table.Table, table.HiddenColumns, subTables, false)

		//schemaOrderBy := dbToStruct.TableToGraphqlOrderBy(table.Table, table.HiddenColumns)

		if len(table.Subs) >= 1 {
			QueryContent = QueryContent + "    " + strings.ToLower(table.Table) + "(sorts:[sort], groupFilters:[groupFilter], filters:[filter], subSorts:[subSort], subFilters:[subFilter], limit: Int, offset: Int): [" + modelAlias + "!]\n"
		} else {
			QueryContent = QueryContent + "    " + strings.ToLower(table.Table) + "(sorts:[sort], groupFilters:[groupFilter], filters:[filter], limit: Int, offset: Int): [" + modelAlias + "!]\n"
		}
		Pagination = Pagination + "    " + strings.ToLower(table.Table) + ":[" + modelAlias + "!]\n"

		utils.WriteFileFormat(structStr, "lambda/graph/models/"+modelAlias+".go")

		authCheck := ""
		if table.CheckAuth.IsLoggedIn {
			authCheck = `_, authErr := gql.CheckAuth(ctx, []int{` + strings.Trim(strings.Replace(fmt.Sprint(table.CheckAuth.Roles), " ", ",", -1), "[]") + `})
	if authErr != nil {
		return nil, authErr
	}`
		}

		subFilterOrders := ""
		subFromCtx := "_"
		resolverReturn := `return result, err`
		importStrconv := ``
		if len(table.Subs) >= 1 {
			subFilterOrders = ", subSorts []*model.SubSort, subFilters []*model.SubFilter"
			subFromCtx = "subs"
			resolverReturn = fmt.Sprintf(`if(len(subs) >= 1){
		resultWithSubs, errorsub := Set%sSubs(ctx, result, subs, subSorts, subFilters)
		return resultWithSubs, errorsub
	}else{
		return result, err
	}`, modelAlias)
			importStrconv = "\"fmt\""
		}

		resolver := fmt.Sprintf(resolverTmplate,
			importStrconv,
			modelAlias,
			subFilterOrders,
			modelAlias,
			authCheck,
			modelAlias,
			subFromCtx,
			table.Table,
			table.Identity,
			parentConnectsions,
			modelAlias,
			resolverReturn,
			modelAlias,
			columnsString,
		)

		if len(table.Subs) >= 1 {
			resolver = resolver + fmt.Sprintf(subTemp,
				modelAlias,
				subTablesMap,
				modelAlias,
				modelAlias)

			resolver = resolver + fmt.Sprintf(setSubTemplate,
				modelAlias,
				modelAlias,
				modelAlias,
				`fmt.Sprintf("%v", `+IdentityForSub+`)`,
				`fmt.Sprintf("%v", `+IdentityForSub+`)`,
				subSetTemps,
			)

		}

		actions := createActions(table, modelAlias, colunms)

		resolver = resolver + actions

		utils.WriteFileFormat(resolver, "lambda/graph/resolvers/"+modelAlias+".go")

		utils.WriteFile(schema, "lambda/graph/schemas/"+GetModelAlias(modelAlias)+".graphql")

		paginationReturn := "return &Paginate, nil"

		if len(table.Subs) >= 1 {
			paginationReturn = fmt.Sprintf(`if len(subs) >= 1 {
				resultWithSubs, errorsub := Set%sSubs(ctx, Paginate.%s, subs, subSorts, subFilters)
				Paginate.%s = resultWithSubs
				return &Paginate, errorsub
			} else {
				return &Paginate, nil
			}`, modelAlias, modelAlias, modelAlias)
		}

		paginationTmplate = paginationTmplate + fmt.Sprintf(`if(target == "%s"){
		%s
		requestColumns = append(requestColumns, "%s")
		requestColumns = append(requestColumns, []string{%s}...)
 		requestColumns = utils.RemoveDuplicateStr(requestColumns)
		query = query.Select(requestColumns)
		data := []*models.%s{}
		
		TabeColumns := %sColumns()
		query, errorFilter := gql.Filter(filters, query,TabeColumns)
		if(errorFilter != nil){
			return &Paginate, errorFilter
		}
		query, errorGroupFilter := gql.GroupFilter(groupFilters, query,TabeColumns)
		if(errorGroupFilter != nil){
			return &Paginate, errorGroupFilter 
		}		
		query, errorOrder := gql.Order(sorts, query, TabeColumns)
		if(errorOrder != nil){
			return &Paginate, errorOrder
		}
		
		pagination := utils.Paging(&utils.Param{
			DB:    query,
			Page:  page,
			Limit: size,
		}, &data)
		Paginate.%s = data
		Paginate.LastPage = pagination.LastPage
		Paginate.Total = int(pagination.Total)
		%s
	}`, strings.ToLower(table.Table), authCheck, table.Identity, parentConnectsions, modelAlias, modelAlias, GetModelAlias(table.Table), paginationReturn) + "\n"

	}

	Pagination = Pagination + "}\n"

	QueryContent = QueryContent + "    paginate(sorts: [sort], groupFilters:[groupFilter], filters:[filter], subSorts:[subSort], subFilters:[subFilter], page:Int!, size:Int!): paginate!\n}\n" + Pagination + "\n"

	paginationTmplate = fmt.Sprintf(paginationTmplate,
		paginationSub,
	)
	paginationTmplate = paginationTmplate + "return &Paginate, nil\n}"

	subscriptionSchema, subscriptions := createSubscription(GqlTables)

	QueryContent = QueryContent + "\n" + subscriptionSchema

	utils.WriteFile(QueryContent, "lambda/graph/schemas/schemas.graphql")

	utils.WriteFileFormat(paginationTmplate, "lambda/graph/resolvers/Paginate.go")

	createGraphqlFile(subscriptions)

	createActionUpdateActions(GqlTables, dbSchema)

}
func createActionUpdateActions(GqlTables []models.GqlTable, dbSchema lambdaModels.DBSCHEMA) {

	mutations := `type Mutation {
`
	mutationTemp := `	%s
	%s
    %s`

	mutationFound := false

	for _, table := range GqlTables {
		if table.Actions.Create || table.Actions.Update {
			modelAlias := GetModelAlias(table.Table)

			schema := TableMetaToGraphql(dbSchema.TableMeta[table.Table], table.Table, []string{"created_at", "updated_at", "deleted_at", "CREATED_AT", "UPDATED_AT", "DELETED_AT", table.Identity}, []string{}, true)

			createMutation := ""
			if table.Actions.Create {
				mutationFound = true
				descrition := "mutation-create"
				if table.Subscription {
					descrition = descrition + ":subscription-" + modelAlias
				}
				createMutation = fmt.Sprintf("\n    \"%s\"\n    create%s(input: %sInput!):%s!", descrition, modelAlias, modelAlias, modelAlias)
			}
			updateMutation := ""
			if table.Actions.Update {
				mutationFound = true
				descrition := "mutation-update"
				if table.Subscription {
					descrition = descrition + ":subscription-" + modelAlias
				}
				updateMutation = fmt.Sprintf("\n    \"%s\"\n    update%s(id: ID!, input:%sInput!):%s!", descrition, modelAlias, modelAlias, modelAlias)
			}
			deleteMutation := ""
			if table.Actions.Delete {
				mutationFound = true
				descrition := "mutation-delete"
				if table.Subscription {
					descrition = descrition + ":subscription-" + modelAlias
				}
				deleteMutation = fmt.Sprintf("\n    \"%s\"\n    delete%s(id: ID!):%s!", descrition, modelAlias, modelAlias)
			}
			mutations = mutations + fmt.Sprintf(mutationTemp, createMutation, updateMutation, deleteMutation)

			utils.WriteFile(schema, "lambda/graph/schemas/"+modelAlias+"Input.graphql")
		}
	}
	mutations = mutations + "\n}"
	if mutationFound {
		utils.WriteFile(mutations, "lambda/graph/schemas/mutations.graphql")
	}

}
func createGraphqlFile(subscriptions []map[string]string) {
	temp := `package graph

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-redis/redis"
	"github.com/gofiber/fiber/v2"
	"github.com/gorilla/websocket"
	"github.com/lambda-platform/lambda/config"
	lambdaPlayground "github.com/lambda-platform/lambda/graphql/playground"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"lambda/lambda/graph/generated"
	"net/http"
	"sync"
	"time"
	%s
)

type Cache struct {
	client redis.UniversalClient
	ttl    time.Duration
}

func Set(e *fiber.App) {

	graphqlConfig := generated.Config{Resolvers: &Resolver{
		%s
		mutex: sync.Mutex{},
	}}

	graphqlHandler := handler.New(generated.NewExecutableSchema(graphqlConfig))
	playgroundHandler := playground.Handler("GraphQL playground", "/query")
	lambdaPlaygroundHandler := lambdaPlayground.Handler("GraphQL playground", "/query")

	graphqlHandler.AddTransport(transport.Options{})
	graphqlHandler.AddTransport(transport.GET{})
	graphqlHandler.AddTransport(transport.POST{})
	graphqlHandler.AddTransport(transport.MultipartForm{})
	graphqlHandler.SetQueryCache(lru.New(1000))

	

	graphqlHandler.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New(100),
	})
	graphqlHandler.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
		Upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	})

	e.Post("/query", func(c *fiber.Ctx) error {
		fasthttpadaptor.NewFastHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			graphqlHandler.ServeHTTP(writer, request)
		})(c.Context())
		return nil
	})

	if config.Config.Graphql.Debug == "true" {
		graphqlHandler.Use(extension.Introspection{})

		e.Get("/play", func(c *fiber.Ctx) error {
			fasthttpadaptor.NewFastHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				playgroundHandler.ServeHTTP(writer, request)
			})(c.Context())
			return nil
		})
	
		e.Get("/play-full", func(c *fiber.Ctx) error {
			fasthttpadaptor.NewFastHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				lambdaPlaygroundHandler.ServeHTTP(writer, request)
			})(c.Context())
			return nil
		})
	}
	
}`
	tempResolver := `package graph

import (
	%s
	"sync"
)

type Resolver struct {
	%s
	mutex          sync.Mutex
}

`

	subs := ""
	subsR := ""
	for _, subscription := range subscriptions {
		subs = subs + fmt.Sprintf("%s%sChannel: map[string]chan *models.%s{},\n", subscription["model"], subscription["action"], subscription["model"])
		subsR = subsR + fmt.Sprintf("%s%sChannel map[string]chan *models.%s\n", subscription["model"], subscription["action"], subscription["model"])

	}
	modelsImport := ""
	if len(subscriptions) >= 1 {
		modelsImport = fmt.Sprintf("\"lambda/graph/models\"")
	}
	graphql := fmt.Sprintf(temp,
		modelsImport,
		subs,
	)
	resolver := fmt.Sprintf(tempResolver,
		modelsImport,
		subsR,
	)

	utils.WriteFileFormat(graphql, "lambda/graph/graphql.go")
	utils.WriteFileFormat(resolver, "lambda/graph/resolver.go")

}
func createSubscription(GqlTables []models.GqlTable) (string, []map[string]string) {

	Subscription := `type Subscription {
`
	SubscriptionTmp := `
    "%s"
    %sCreated: %s!
`
	SubscriptionTmpUpdated := `
    "%s"
    %sUpdated: %s!
`
	SubscriptionTmpDeleted := `
    "%s"
    %sDeleted: %s!
`
	Subscriptions := []map[string]string{}
	SubscriptionFound := false
	for _, table := range GqlTables {

		if table.Subscription && table.Actions.Create {
			SubscriptionFound = true
			modelAlias := GetModelAlias(table.Table)
			Subscription = Subscription + fmt.Sprintf(SubscriptionTmp, "subscription-created:"+modelAlias, modelAlias, modelAlias)

			Subscriptions = append(Subscriptions, map[string]string{
				"model":  modelAlias,
				"action": "Created",
			})
		}

		if table.Subscription && table.Actions.Update {
			SubscriptionFound = true
			modelAlias := GetModelAlias(table.Table)
			Subscription = Subscription + fmt.Sprintf(SubscriptionTmpUpdated, "subscription-updated:"+modelAlias, modelAlias, modelAlias)

			Subscriptions = append(Subscriptions, map[string]string{
				"model":  modelAlias,
				"action": "Updated",
			})
		}

		if table.Subscription && table.Actions.Delete {
			SubscriptionFound = true
			modelAlias := GetModelAlias(table.Table)
			Subscription = Subscription + fmt.Sprintf(SubscriptionTmpDeleted, "subscription-deleted:"+modelAlias, modelAlias, modelAlias)

			Subscriptions = append(Subscriptions, map[string]string{
				"model":  modelAlias,
				"action": "Deleted",
			})
		}

	}
	Subscription = Subscription + "\n}"
	if SubscriptionFound {
		return Subscription, Subscriptions
	} else {
		return "", Subscriptions
	}

}
func createActions(table models.GqlTable, modelAlias string, colunms []map[string]string) string {

	createTemp := `
func Create%s(ctx context.Context, input model.%sInput) (*models.%s, error) {
			row := models.%s{}
			%s
			
			err := DB.DB.Create(&row).Error
			return &row, err
		}`
	updateTemp := `
func Update%s(ctx context.Context, id string, input model.%sInput) (*models.%s, error) {
			row := models.%s{}
			DB.DB.Where("%s = ?", id).Find(&row)
			%s
			err := DB.DB.Save(&row).Error

			return &row, err
		}`
	deleteTemp := `
func Delete%s(ctx context.Context, id string) (*models.%s, error) {
			row := models.%s{}
			DB.DB.Where("%s = ?", id).Find(&row)
			err := DB.DB.Where("%s = ?", id).Delete(&models.%s{}).Error
			return &row, err
		}`

	actions := ""

	columnsWithInput := ""
	//colunms = strings.ReplaceAll(colunms, "\"", "")
	//colunms = strings.ReplaceAll(colunms, " ", "")
	for _, column := range colunms {
		if column["column"] != table.Identity && column["column"] != "created_at" && column["column"] != "updated_at" && column["column"] != "deleted_at" && column["column"] != "CREATED_AT" && column["column"] != "UPDATED_AT" && column["column"] != "DELETED_AT" {

			columnReady := GetModelAlias(column["column"])

			//	fmt.Println(column["column"])
			//	fmt.Println(column["dataType"])
			if column["nullable"] != "YES" || column["dataType"] == "datetime" {
				columnsWithInput = columnsWithInput + fmt.Sprintf("row.%s = input.%s\n",
					columnReady,
					columnReady,
				)
			} else {
				columnsWithInput = columnsWithInput + fmt.Sprintf("row.%s = input.%s\n",
					columnReady,
					columnReady,
				)
			}

		}
	}

	if table.Actions.Create {

		actions = actions + fmt.Sprintf(createTemp,
			modelAlias,
			modelAlias,
			modelAlias,
			modelAlias,
			columnsWithInput,
		)
	}

	if table.Actions.Update {

		actions = actions + fmt.Sprintf(updateTemp,
			modelAlias,
			modelAlias,
			modelAlias,
			modelAlias,
			table.Identity,
			columnsWithInput,
		)
	}

	if table.Actions.Delete {

		actions = actions + fmt.Sprintf(deleteTemp,
			modelAlias,
			modelAlias,
			modelAlias,
			table.Identity,
			table.Identity,
			modelAlias,
		)
	}

	return actions

}
func GQLInit(dbSchema lambdaModels.DBSCHEMA, graphqlchemas []models.ProjectSchemas) {

	if len(graphqlchemas) >= 1 {

		AbsolutePath := utils.AbsolutePath()

		modelsPatch := "lambda/graph/models"
		generatedPath := "lambda/graph/generated"
		schemaPatch := "lambda/graph/schemas"
		resolversPatch := "lambda/graph/resolvers"
		schemaCommonPatch := "lambda/graph/schemas-common"
		if _, err := os.Stat(modelsPatch); os.IsNotExist(err) {

			os.MkdirAll(modelsPatch, 0777)
			os.MkdirAll(schemaPatch, 0777)
			os.MkdirAll(resolversPatch, 0777)
			os.MkdirAll(schemaCommonPatch, 0777)
			os.MkdirAll(generatedPath, 0777)

		} else {

			os.RemoveAll("lambda/graph")
			os.RemoveAll(modelsPatch)
			os.RemoveAll(schemaPatch)
			os.RemoveAll(resolversPatch)
			os.RemoveAll(schemaCommonPatch)
			os.MkdirAll(modelsPatch, 0777)
			os.MkdirAll(schemaPatch, 0777)
			os.MkdirAll(resolversPatch, 0777)
			os.MkdirAll(schemaCommonPatch, 0777)
			os.MkdirAll(generatedPath, 0777)
		}

		copy.Copy(AbsolutePath+"graphql/schemas-common/", "lambda/graph/schemas-common/")

		gqlgenFile, _ := os.ReadFile(AbsolutePath + "graphql/gqlgen.yml.example")
		//gqlgenFileContent := strings.ReplaceAll(string(gqlgenFile), "PROJECTNAME", projectName)
		//gqlgenFileContent = strings.ReplaceAll(gqlgenFileContent, "PROJECTAPTH", dir)
		utils.WriteFile(string(gqlgenFile), "lambda/graph/gqlgen.yml")

		//graphqlFile, _ := os.ReadFile(AbsolutePath + "/graph/graphql.go.exmaple")
		//graphqlFileContent := strings.ReplaceAll(string(graphqlFile), "PROJECTNAME", projectName)
		//WriteFile(graphqlFileContent, dir+"/graph/graphql.go")
		GenerateSchema(graphqlchemas, dbSchema)
		Generate()
		fmt.Println("GRAPHQL INIT DONE")
	} else {
		graphqlPatch := "lambda/graph"
		if _, err := os.Stat(graphqlPatch); os.IsNotExist(err) {

			os.MkdirAll(graphqlPatch, 0777)

		} else {

			os.RemoveAll(graphqlPatch)
			os.MkdirAll(graphqlPatch, 0777)

		}
		utils.WriteFile(`package graph
import "github.com/gofiber/fiber/v2"
func Set(e *fiber.App) {}
`, "lambda/graph/graphql.go")
	}

}
