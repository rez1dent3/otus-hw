package actions

import "net/http"

func Ping(writer http.ResponseWriter, _ *http.Request) {
	_, _ = writer.Write([]byte("ok"))
}
