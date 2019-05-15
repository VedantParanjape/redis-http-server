package main

import "fmt"
import "net/http"
import "strings"
import "log"
import "github.com/garyburd/redigo/redis"

type return_data struct {
	retstring     string
	retint64      int64
	retbulkstring []byte
	retarray      interface{}
	err           error
	rettype       string
}

func get_conn_redis(port string) redis.Conn {
	c, err := redis.Dial("tcp", port)
	if err != nil {
		panic(err)
	}
	return c
}

func do_redis(handle redis.Conn, command string, args ...string) return_data {
	s := make([]interface{}, len(args))
	for i, v := range args {
		s[i] = v
	}
	ret, err := handle.Do(command, s[0:]...)
	var rtrn return_data

	switch ret.(type) {
	case int64:
		rtrn.retint64, rtrn.err = redis.Int64(ret, err)
		rtrn.rettype = "int"
	case string:
		rtrn.retstring, rtrn.err = redis.String(ret, err)
		rtrn.rettype = "string"
	case []byte:
		rtrn.retbulkstring, rtrn.err = redis.Bytes(ret, err)
		rtrn.rettype = "byte"
	case interface{}:
		rtrn.retarray = ret
		rtrn.rettype = "interface"
	}
	return rtrn
}

func do_http(w http.ResponseWriter, r *http.Request, hndl redis.Conn) {
	s := strings.Split(r.URL.String(), "/")

	ret := do_redis(hndl, s[2], s[3:]...)
	fmt.Println(ret.err)

	if ret.err != nil {
		fmt.Fprintln(w, ">>", ret.err)
	} else {
		switch ret.rettype {
		case "string":
			log.Print("type: string, output: ", ret.retstring)
			fmt.Fprintln(w, ">>", ret.retstring)

		case "int":
			log.Print("type: int, output: ", ret.retint64)
			fmt.Fprintln(w, ">>", ret.retint64)

		case "byte":
			output, _ := redis.String(ret.retbulkstring, ret.err)
			log.Print("type: bulkstring, output: ", output)
			fmt.Fprintln(w, ">>", output)

		case "interface":
			log.Print("type: inteface, output: ", ret.retarray)
			fmt.Fprintln(w, ">>", ret.retarray)
		}
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
