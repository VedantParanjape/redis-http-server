package main

import "fmt"
import "net/http"
import "strings"
import "github.com/garyburd/redigo/redis"

func get_conn_redis(port string) redis.Conn {
	c, err := redis.Dial("tcp", port)
	if err != nil {
		panic(err)
	}
	return c
}

func do_redis(handle redis.Conn, command string, args ...string) (string, error) {
	s := make([]interface{}, len(args))
	for i, v := range args {
		s[i] = v
	}
	return  redis.String(handle.Do(command, s[0:]...))

}

func do_http(w http.ResponseWriter, r *http.Request, hndl redis.Conn) {
	s := strings.Split(r.URL.String(), "/")

	ret, err := do_redis(hndl, s[2], s[3:]...)
	fmt.Println(err)

	if err != nil {
		fmt.Fprintln(w, ">>", err)
	} else {
		fmt.Fprintln(w, ">>", ret)
	}
}

func main() {
	handle := get_conn_redis(":6379")
	defer handle.Close()

	http.HandleFunc("/do/", func(w http.ResponseWriter, r *http.Request) {
		do_http(w, r, handle)
	})
	http.ListenAndServe(":8001", nil)
}
