package main

import (
	"database/sql"
	"fmt"
	//_ "github.com/Go-SQL-Driver/MySQL"
	"reflect"
	"strings"
	"errors"
)

func main() {
	//db, err := sql.Open("mysql", "root:@/studygolang?charset=utf8")
	//if err != nil {
		//panic(err)
	//}
	//defer db.Close()
	//// insert(db)
	////update(db)
	//var user = struct {
	    //Username string
	    //Email string
	//}{
		//"22",
		//"fwef@163.com",
	//}
	//FindAll(db, &user)
	arr := []string{}
	fmt.Println(strings.Join(arr, " AND "))

}

func insert(db *sql.DB) {
	stmt, err := db.Prepare("INSERT INTO user_login(uid, email, username, passcode, passwd) VALUES(?, ?, ?, ?, ?)")
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec("1", "studygolang@gmail.com", "polaris", "2few3", "hwifwe23fweifwef")
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
}

func FindAll(db *sql.DB, model interface{}) {
	columns, err := extractFields(model)
	if err != nil {
		panic(err)
	}
	columnStr := strings.Join(columns, ",")
	stmt, err := db.Prepare("SELECT " + columnStr + " FROM user_login WHERE uid=?")
	if err != nil {
		panic(err)
	}
	
	stmt.QueryRow(1)
}

func extractFields(obj interface{}) ([]string, error) {
	dataStructValue := reflect.Indirect(reflect.ValueOf(obj))
	if dataStructValue.Kind() != reflect.Struct {
		return nil, errors.New("expected a struct")
	}

	dataStructType := dataStructValue.Type()
	count := dataStructType.NumField()
	fieldSlice := make([]string, count)
	for i := 0; i < count; i++ {
		field := dataStructType.Field(i)
		fieldSlice[i] = strings.ToLower(field.Name)
	}
	return fieldSlice, nil
}

// convertStruct2Slice 将struct转为map，其中struct的field为map的key，但是key的首字母小写
func convertStruct2Slice(obj interface{}) ([]string, []interface{}, error) {
	dataStructValue := reflect.Indirect(reflect.ValueOf(obj))
	if dataStructValue.Kind() != reflect.Struct {
		return nil, nil, errors.New("expected a struct")
	}

	dataStructType := dataStructValue.Type()
	count := dataStructType.NumField()
	fieldSlice := make([]string, count)
	valueSlice := make([]interface{}, count)
	for i := 0; i < count; i++ {
		field := dataStructType.Field(i)
		fieldName := field.Name

		fieldSlice[i] = strings.ToLower(fieldName)
		val := dataStructValue.FieldByName(fieldName).Interface()
		fmt.Printf("%v\n", val)
		valueSlice[i] = &val
	}
	return fieldSlice, valueSlice, nil
}

func find(db *sql.DB) {
	stmt, err := db.Prepare("SELECT username, email FROM user_login WHERE uid=?")
	if err != nil {
		panic(err)
	}

	/*
		var user = struct {
		    username string
		    email string
		}{}
	*/
	var username string
	var email string
	err = stmt.QueryRow(1).Scan(&email, &username)
	//err = stmt.QueryRow(1).Scan(&user.email, &user.username)
	if err != nil {
		panic(err)
	}
	//fmt.Println(user.email, user.username)
	fmt.Println(email, username)
}

func update(db *sql.DB) {
	stmt, err := db.Prepare("UPDATE user_login SET username=? WHERE uid=?")
	if err != nil {
		panic(err)
	}

	result, err := stmt.Exec("xuxinhua", 1)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.RowsAffected())
}
