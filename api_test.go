package datacenter

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/suifengpiao14/datacenter/logger"
	"github.com/suifengpiao14/datacenter/source"
)

func TestAPI(t *testing.T) {
	route := "/api/1/hello"
	method := "POST"
	api := &API{
		Methods: "post,get",
		Route:   route,
		InputLineSchema: `version=http://json-schema.org/draft-07/schema,id=input,direction=in
		fullname=pageIndex,dst=pageIndex,format=number,required
		fullname=pageSize,dst=pageSize,format=number,required`,
		OutputLineSchema: `version=http://json-schema.org/draft-07/schema,id=output,direction=out
		fullname=items[].content,src=PaginateOut.#.content,required
		fullname=items[].createdAt,src=PaginateOut.#.created_at,required
		fullname=items[].deletedAt,src=PaginateOut.#.deleted_at,required
		fullname=items[].description,src=PaginateOut.#.description,required
		fullname=items[].icon,src=PaginateOut.#.icon,required
		fullname=items[].id,src=PaginateOut.#.id,required
		fullname=items[].key,src=PaginateOut.#.key,required
		fullname=items[].label,src=PaginateOut.#.label,required
		fullname=items[].thumb,src=PaginateOut.#.thumb,required
		fullname=items[].title,src=PaginateOut.#.title,required
		fullname=items[].updatedAt,src=PaginateOut.#.updated_at,required
		fullname=pageInfo.pageIndex,src=input.pageIndex,required
		fullname=pageInfo.pageSize,src=input.pageSize,required
		fullname=pageInfo.total,src=PaginateTotalOut,required`,
		//{pageInfo:{pageIndex:input.pageIndex,pageSize:input.pageSize,total:PaginateTotalOut},items:{content:PaginateOut.#.content,createdAt:PaginateOut.#.created_at,deletedAt:PaginateOut.#.deleted_at}|@group}
		PreScript: `
		input["Offset"]=int(input["pageIndex"]) * int(input["pageSize"])
		input["Limit"]=int(input["pageSize"])
		`,
		MainScript: `
		
		PaginateTotalOut:=execSQLTPL("PaginateTotal",input)
		PaginateOut :=execSQLTPL("Paginate",input)
		out:=""
		out=gsjson.Set(out,"input",input)
		out=gsjson.SetRaw(out,"PaginateTotalOut",PaginateTotalOut)
		out=gsjson.SetRaw(out,"PaginateOut",PaginateOut)
		return  out
		`,
		PostScript: ``,
		AfterEvent: ``,
	}
	capi, err := NewApiCompiled(api)
	if err != nil {
		panic(err)
	}
	tplStr := `
	{{define "PaginateWhere"}} {{end}}
	{{define "PaginateTotal"}} select count(*) as count from component where 1=1 {{template "PaginateWhere" .}} and deleted_at is null; {{end}}
	{{define "Paginate"}} select * from component where 1=1 {{template "PaginateWhere" .}} and deleted_at is null order by updated_at desc limit :Offset,:Limit ; {{end}}
	`
	sourceConfig := `
	   	{"logLevel":"debug","dsn":"root:123456@tcp(10.0.11.125:3306)/office_web_site?charset=utf8&timeout=1s&readTimeout=5s&writeTimeout=5s&parseTime=False&loc=Local&multiStatements=true","timeout":30}
	   	`
	var dbConfig source.DBExecProviderConfig
	err = json.Unmarshal([]byte(sourceConfig), &dbConfig)
	if err != nil {
		panic(err)
	}
	dbSource := source.DBExecProvider{Config: dbConfig}
	capi.AddTpl("", tplStr, &dbSource)

	container := NewContainer()
	container.RegisterAPI(capi)

	routeCapi, ok := container.GetCApi(route, method)
	if !ok {
		panic(ok)
	}

	inputJson := `{"pageIndex":"","pageSize":"20"}`
	out, err := routeCapi.Run(context.TODO(), inputJson)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)
}

func TestAPIMemory(t *testing.T) {
	route := "/api/1/hello"
	method := "POST"
	api := &API{
		Methods: "post,get",
		Route:   route,
		InputLineSchema: `version=http://json-schema.org/draft-07/schema,id=input,direction=in
		fullname=pageIndex,dst=pageIndex,format=number,required
		fullname=pageSize,dst=pageSize,format=number,required`,
		OutputLineSchema: `version=http://json-schema.org/draft-07/schema,id=output,direction=out
		fullname=items[].content,src=PaginateOut.#.content,required
		fullname=items[].createdAt,src=PaginateOut.#.created_at,required
		fullname=items[].deletedAt,src=PaginateOut.#.deleted_at,required
		fullname=items[].description,src=PaginateOut.#.description,required
		fullname=items[].icon,src=PaginateOut.#.icon,required
		fullname=items[].id,src=PaginateOut.#.id,required
		fullname=items[].key,src=PaginateOut.#.key,required
		fullname=items[].label,src=PaginateOut.#.label,required
		fullname=items[].thumb,src=PaginateOut.#.thumb,required
		fullname=items[].title,src=PaginateOut.#.title,required
		fullname=items[].updatedAt,src=PaginateOut.#.updated_at,required
		fullname=pageInfo.pageIndex,src=input.pageIndex,required
		fullname=pageInfo.pageSize,src=input.pageSize,required
		fullname=pageInfo.total,src=PaginateTotalOut,required`,
		//{pageInfo:{pageIndex:input.pageIndex,pageSize:input.pageSize,total:PaginateTotalOut},items:{content:PaginateOut.#.content,createdAt:PaginateOut.#.created_at,deletedAt:PaginateOut.#.deleted_at}|@group}
		PreScript: `
		input["Offset"]=int(input["pageIndex"]) * int(input["pageSize"])
		input["Limit"]=int(input["pageSize"])
		`,
		MainScript: `
		
		PaginateTotalOut:=execSQLTPL(ctx,"PaginateTotal",input)
		PaginateOut :=execSQLTPL(ctx,"Paginate",input)
		out:=""
		out=gsjson.Set(out,"input",input)
		out=gsjson.SetRaw(out,"PaginateTotalOut",PaginateTotalOut)
		out=gsjson.SetRaw(out,"PaginateOut",PaginateOut)
		return  out
		`,
		PostScript: ``,
		AfterEvent: ``,
	}
	capi, err := NewApiCompiled(api)
	if err != nil {
		panic(err)
	}
	tplStr := `
	{{define "PaginateWhere"}} {{end}}
	{{define "PaginateTotal"}} select count(*) as count from component where 1=1 {{template "PaginateWhere" .}} and deleted_at is null; {{end}}
	{{define "Paginate"}} select * from component where 1=1 {{template "PaginateWhere" .}} and deleted_at is null order by updated_at desc limit :Offset,:Limit ; {{end}}
	`
	/* 	sourceConfig := `
	   	{"logLevel":"debug","dsn":"root:123456@tcp(10.0.11.125:3306)/office_web_site?charset=utf8&timeout=1s&readTimeout=5s&writeTimeout=5s&parseTime=False&loc=Local&multiStatements=true","timeout":30}
	   	`
	*/
	/* 	var dbConfig db.DBExecProviderConfig
	   	err = json.Unmarshal([]byte(sourceConfig), &dbConfig)
	   	if err != nil {
	   		panic(err)
	   	}
	   	source := db.DBExecProvider{Config: dbConfig}
	*/
	selectSQL := "select * from component where 1=1   and deleted_at is null order by updated_at desc limit 0,20 ;"
	countSQL := "select count(*) as count from component where 1=1   and deleted_at is null;"
	source := MemorySource{InOutMap: make(map[string]string)}
	source.InOutMap[selectSQL] = "[{\"content\":\"\",\"created_at\":\"2022-11-17 22:03:11\",\"deleted_at\":\"\",\"description\":\"??????????????????????????????,???????????????????????????????????????,??????????????????\",\"icon\":\"flaticon-research-1\",\"id\":\"6\",\"key\":\"pay\",\"label\":\"\",\"thumb\":\"../assets/images/case-studies/case-studies-another6.jpg\",\"title\":\"????????????\",\"updated_at\":\"2022-11-17 22:03:23\"},{\"content\":\"\",\"created_at\":\"2022-11-17 22:01:43\",\"deleted_at\":\"\",\"description\":\"??????????????????????????????????????????,??????????????????????????????????????????????????????,?????????????????????????????????\",\"icon\":\"flaticon-branding\",\"id\":\"5\",\"key\":\"fsm\",\"label\":\"\",\"thumb\":\"../assets/images/case-studies/case-studies-another5.jpg\",\"title\":\"???????????????\",\"updated_at\":\"2022-11-17 22:02:05\"},{\"content\":\"\",\"created_at\":\"2022-11-17 21:59:47\",\"deleted_at\":\"\",\"description\":\"??????????????????,??????????????????????????????????????????,?????????????????????,?????????????????????,??????????????????????????????,??????????????????????????????????????????,???????????????????????????,?????????????????????????????????\",\"icon\":\"flaticon-education\",\"id\":\"4\",\"key\":\"validator\",\"label\":\"\",\"thumb\":\"../assets/images/case-studies/case-studies-another4.jpg\",\"title\":\"??????????????????\",\"updated_at\":\"2022-11-17 22:00:44\"},{\"content\":\"\",\"created_at\":\"2022-11-17 21:58:58\",\"deleted_at\":\"\",\"description\":\"excel????????????????????????,????????????????????????,?????????????????????????????????,?????????????????????????????????\",\"icon\":\"flaticon-web-programming\",\"id\":\"3\",\"key\":\"excel\",\"label\":\"\",\"thumb\":\"../assets/images/case-studies/case-studies-another3.jpg\",\"title\":\"excel????????????\",\"updated_at\":\"2022-11-17 21:59:15\"},{\"content\":\"\",\"created_at\":\"2022-11-17 21:53:42\",\"deleted_at\":\"\",\"description\":\"?????????????????????????????????????????????????????????????????????????????????????????????????????????.??????auth2.0,openID???????????????\",\"icon\":\"flaticon-research-1\",\"id\":\"2\",\"key\":\"account\",\"label\":\"\",\"thumb\":\"../assets/images/case-studies/case-studies-another2.jpg\",\"title\":\"????????????\",\"updated_at\":\"2022-11-17 21:54:42\"},{\"content\":\"\",\"created_at\":\"2022-11-17 21:51:27\",\"deleted_at\":\"\",\"description\":\"?????????????????????????????????????????????,????????????????????????,???????????????????????????????????????\",\"icon\":\"flaticon-branding\",\"id\":\"1\",\"key\":\"address\",\"label\":\"\",\"thumb\":\"../assets/images/case-studies/case-studies-another1.jpg\",\"title\":\"?????????\",\"updated_at\":\"2022-11-17 21:53:08\"}]"
	source.InOutMap[countSQL] = "6"
	capi.AddTpl("", tplStr, &source)

	container := NewContainer()
	container.RegisterAPI(capi)

	routeCapi, ok := container.GetCApi(route, method)
	if !ok {
		panic(ok)
	}

	inputJson := `{"pageIndex":"","pageSize":"20"}`
	ctx := context.Background()
	ctx = context.WithValue(ctx, "traceID", "12345")

	readChain := logger.GetLoggerChain()
	go func() {
		for {
			data := <-readChain
			fmt.Println(data)
		}
	}()
	out, err := routeCapi.Run(ctx, inputJson)
	if err != nil {
		panic(err)
	}
	fmt.Println(out)

}
