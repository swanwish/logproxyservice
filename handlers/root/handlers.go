package root

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/swanwish/go-common/logs"
	"github.com/swanwish/go-common/web"
	"github.com/swanwish/logproxyservice/executor"
)

type Handlers struct {
}

func (h Handlers) GetPathPrefix() string {
	return ""
}

func (h Handlers) InitRouter(r *mux.Router) {
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.HandleFunc("/crashlog", getCrashLog).Methods("GET")
}

func rootHandler(rw http.ResponseWriter, r *http.Request) {
	ctx := web.CreateHandlerContext(rw, r)
	ctx.ReplyOK()
}

func getCrashLog(rw http.ResponseWriter, r *http.Request) {
	ctx := web.CreateHandlerContext(rw, r)
	shExecutor := executor.NewShExecutor()
	result, err := shExecutor.RunScript("logcat -b crash -v long -d")
	if err != nil {
		logs.Errorf("Failed to run script, the error is %#v", err)
		ctx.ReplyError(err.Error(), http.StatusBadRequest)
		return
	}
	htmlContent := fmt.Sprintf(`<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta http-equiv="X-UA-Compatible" content="IE=edge">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>Crash Log</title>
	</head>
	<body>
		<h1>Crash Logs</h1>
		<pre>%s</pre>
	</body>
	</html>`, result)
	ctx.ReplyHtml(htmlContent)
}
