package  main

import "fmt"
import "github.com/garyburd/redigo/redis"

func get_conn_redis(port string) redis.Conn{
	c, err := redis.Dial("tcp",port)
	if err != nil{
		panic(err)
	}
	return c
}

func do_redis(handle redis.Conn, command string, args ...string) (string, error){
   s := make([]interface{}, len(args))
   for i, v := range args {
      s[i] = v
   }
   return redis.String(handle.Do(command, s[0:]...))   //// update the code such that only when output is string compatible this works
}

func main(){
	handle := get_conn_redis(":6379")
	defer handle.Close()

	output, err := do_redis(handle, "SET", "asa", "vedant")
	fmt.Println(output, err)
	
	input, err := do_redis(handle, "GET", "asa")
	fmt.Println(input, err)
	
	dele, err := do_redis(handle, "DEL", "asa")
    fmt.Println(dele, err)

}