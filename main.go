package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"reflect"
	"strings"
	"time"
)

/*
https://gorm.io/docs/index.html
*/

var AppLog logger.Interface

/*

Sample User Record:

JSON >

{
    "UserID": "6285557743a8bdeb2aa5dc07",
    "FirstName": "Sonia",
    "LastName": "Livingston",
    "Email": "sonialivingston@hinway.com",
    "Phone": "+1 (957) 570-2414",
    "Active": false,
    "Balance": "$1,174.11",
}

To the above JSON Record, we add the "string_rep" column : which is a '#' delimited values of each column

"string_rep" is Indexed and used for searching

Search Types:

- Exact Match

	- OR Search , Example : first_name = "Sonia" (OR) first_name = "Wendy" : Type-1
	- AND Search, Example : first_name = "Wendy" (AND) last_name = "Lawson" : Type-2

- Non Exact Match (Pattern Search)

	- OR Search , Example : user_id (CONTAINS) "62855" OR phone (CONTAINS) "570-2414" : Type-3
	- AND Search , Example : user_id (CONTAINS) "62855" AND phone (CONTAINS) "570-2414" : Type-4

--- SQL Query For : Type-1 ---

first_name ~* '#Sonia#' OR  first_name ~* '#Wendy#'

--- SQL Query For : Type-2 ---

first_name ~* '#Wendy#' AND  last_name ~* '#Lawson#'

--- SQL Query For : Type-3 ---

user_id ~* '62855' OR  phone ~* '570-2414'

--- SQL Query For : Type-4 ---

user_id ~* '62855' AND  phone ~* '570-2414'

Constructing "string_rep" before inserting the record back into the DB:

"StringRep": "#6285557743a8bdeb2aa5dc07#Sonia#Livingston#sonialivingston@hinway.com#+1 (957) 570-2414#false#$1,174.11#"

SQL Table Row (with strin_rep column) >

user_id    | 6285557743a8bdeb2aa5dc07
first_name | Sonia
last_name  | Livingston
email      | sonialivingston@hinway.com
phone      | +1 (957) 570-2414
active     | f
balance    | $1,174.11
string_rep | #6285557743a8bdeb2aa5dc07#Sonia#Livingston#sonialivingston@hinway.com#+1 (957) 570-2414#false#$1,174.11#

FYI : For each row, the value of string_rep column is a '#' delimited string/representation of all the
      values of other columns. This will help in searching for one or more strings which may be present
      in any of the other columns.

*/

type Search int

const (
	SearchAND Search = iota
	SearchOR
	SearchSingle
)

type ExactMatch bool

type User struct {
	UserBasic
	StringRep string `gorm:"index"`
}

type UserBasic struct {
	UserID    string `gorm:"primaryKey;column:user_id;"`
	FirstName string `gorm:"index:first_name;default:NA;column:first_name;"`
	LastName  string `gorm:"index:last_name;default:NA;column:last_name;"`
	Email     string `gorm:"index:email;default:no-reply@none.com;column:email;"`
	Phone     string `gorm:"index:phone;default:000-000-0000;column:phone;"`
	Active    bool   `gorm:"default:false;column:active;"`
	Balance   string `gorm:"default:0;column:balance;"`
}

var (
	PGSQLMETADATAHOST = ""
	PGSQLMETADATAPASS = ""
	PGSQLMETADATAUSER = ""
)

func Initialize() {
	PGSQLMETADATAHOST = os.Getenv("PGSQLMETADATAHOST")
	if PGSQLMETADATAHOST == "" {
		log.Fatal("environment variable PGSQLMETADATAHOST is not set, exiting...")
	}

	PGSQLMETADATAPASS = os.Getenv("PGSQLMETADATAPASS")
	if PGSQLMETADATAPASS == "" {
		log.Fatal("environment variable PGSQLMETADATAPASS is not set, exiting...")
	}

	PGSQLMETADATAUSER = os.Getenv("PGSQLMETADATAUSER")
	if PGSQLMETADATAUSER == "" {
		log.Fatal("environment variable PGSQLMETADATAUSER is not set, exiting...")
	}
}

func createRecord(user User, db *gorm.DB) (*gorm.DB, error) {
	var result *gorm.DB
	result = db.Create(&user)
	if result.Error != nil {
		log.Printf("error : could not create record : %v", result.Error.Error())
		return result, result.Error
	}
	log.Printf("result.RowsAffected : %v", result.RowsAffected)
	return result, nil
}

func InitializeLogger() {
	AppLog = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
}

func getUser() User {

	userBasic := UserBasic{
		UserID:    "628558706b92ac31676d779b",
		FirstName: "Mandy",
		LastName:  "Knowles",
		Email:     "mandyknowles@hinway.com",
		Phone:     "+1 (926) 579-2448",
		Active:    false,
		Balance:   "$3,682.63",
	}

	stringRep := getStringRep(userBasic)

	user := User{
		UserBasic: userBasic,
		StringRep: stringRep,
	}

	return user
}

func getStringRep(user UserBasic) string {
	v := reflect.ValueOf(user)
	values := make([]interface{}, v.NumField())
	stringRep := ""
	stringRep += "#"
	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
		stringRep += fmt.Sprintf("%v#", v.Field(i).Interface())
	}
	return stringRep
}

// this function will only update "string_rep" columns
// and return the user
func getUserFromBasic(user UserBasic) User {
	stringRep := getStringRep(user)

	myUser := User{
		UserBasic: user,
		StringRep: stringRep,
	}
	return myUser
}

func updateStringRepForUser(db *gorm.DB, userID string) {
	var userFromBackend User

	db.Where(map[string]interface{}{"user_id": userID}).Find(&userFromBackend)

	log.Printf("BACKEND_QUERY : userFromBackend : %v", userFromBackend)

	stringRep := getStringRep(userFromBackend.UserBasic)

	userFromBackend.StringRep = stringRep

	savedResult := db.Save(&userFromBackend)

	log.Printf("userID : (%v) , savedResult.RowsAffected : (%v)", userID, savedResult.RowsAffected)

}

/*
Automatically migrate your schema, to keep your schema up to date.

AutoMigrate will create tables, missing foreign keys, constraints, columns and indexes.
It will change existing column’s type if its size, precision, nullable changed.
It WON’T delete unused columns to protect your data.
*/

func InitializeTables(db *gorm.DB) error {
	err := db.AutoMigrate(&User{})
	if err != nil {
		return err
	}
	return nil
}

type Tabler interface {
	TableName() string
}

// TableName overrides the table name used by User to `users`
func (User) TableName() string {
	return "user_records"
}

func main() {
	var err error

	Initialize()
	InitializeLogger()

	dsn := fmt.Sprintf(
		"host=%v user=%v password=%v dbname=testdb port=5432 sslmode=disable TimeZone=America/Los_Angeles",
		PGSQLMETADATAHOST,
		PGSQLMETADATAUSER,
		PGSQLMETADATAPASS,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: AppLog,
	})
	if err != nil {
		log.Fatalf("error : %v", err.Error())
	}
	log.Printf("%v", db)

	// ----------------------------------------------------------------------------------------------------

	log.Printf("---[Create/Initialize Table]---")

	err = InitializeTables(db)
	if err != nil {
		log.Printf("error : could not create tables : %v", err.Error())
		return
	}

	// Create a single record
	sampleUser := getUser()
	_, _ = createRecord(sampleUser, db)

	// ----------------------------------------------------------------------------------------------------

	// Bulk Insert

	log.Printf("---[Bulk Insert]---")
	userList := GetUserRecords()
	db.Create(userList)

	// ----------------------------------------------------------------------------------------------------

	// Delete all the rows

	log.Printf("---[Deleting All Rows]---")

	db.Exec("DELETE FROM users")

	// ----------------------------------------------------------------------------------------------------

	log.Printf("---[Dropping Table]---")

	err = db.Migrator().DropTable(&User{})
	if err != nil {
		log.Printf("error : could not drop table : %v", err.Error())
	}

	// ----------------------------------------------------------------------------------------------------

	// Again, re-create the table

	log.Printf("---[Creating Table]---")

	err = InitializeTables(db)
	if err != nil {
		log.Printf("error : could not create tables : %v", err.Error())
		return
	}

	// ----------------------------------------------------------------------------------------------------

	// Insert Using Batch Pool Size
	log.Printf("---[Insert In Batches]---")
	db.CreateInBatches(userList, 4)

	// ----------------------------------------------------------------------------------------------------

	// Get all records

	log.Printf("---[Get All Records]---")

	var users []User

	_ = db.Find(&users)

	//for _, user := range users {
	//
	//	log.Printf("user.UserID    : %v", user.UserID)
	//	log.Printf("user.FirstName : %v", user.FirstName)
	//	log.Printf("user.LastName  : %v", user.LastName)
	//
	//	prettyPrintData(user)
	//}

	log.Printf("Total Number of Records : %v", len(users))

	// ----------------------------------------------------------------------------------------------------

	// Upsert / On Conflict

	log.Printf("---[Upsert / On Conflict]---")

	user1basic := UserBasic{
		UserID:    "628555772a8b7b9926ffb917",
		FirstName: "Wendy-1",
		LastName:  "Lawson-1",
		Email:     "wendylawson@hinway.com",
		Phone:     "+1 (907) 523-2723",
		Active:    false,
		Balance:   "$200,000.00",
	}

	log.Printf("user1basic >")
	prettyPrintData(user1basic)

	user1 := getUserFromBasic(user1basic)

	log.Printf("user1 >")
	prettyPrintData(user1)

	// Update all columns, except primary keys, to new value on conflict

	result1 := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&user1)

	log.Printf("Upsert / On Conflict : result.RowsAffected : %v", result1.RowsAffected)

	// ----------------------------------------------------------------------------------------------------

	// Upsert / On Conflict

	log.Printf("---[Upsert / On Conflict]---")

	user2basic := UserBasic{
		UserID:  "628555772a8b7b9926ffb917",
		Email:   "wendylawson@hinway.com",
		Active:  false,
		Balance: "$200,000.00",
	}

	log.Printf("user2basic >")
	prettyPrintData(user2basic)

	user2 := getUserFromBasic(user2basic)

	log.Printf("user2 >")
	prettyPrintData(user2)

	// Update all columns, except primary keys, to new value on conflict

	result2 := db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(&user2)

	log.Printf("Upsert / On Conflict : result.RowsAffected : %v", result2.RowsAffected)

	//updateStringRepForUser(db, user2basic.UserID)

	// ----------------------------------------------------------------------------------------------------

	user3basic := UserBasic{
		UserID:    "628555772a8b7b9926ffb917",
		FirstName: "Wendy--2",
		LastName:  "Lawson--2",
	}

	log.Printf("user3basic >")
	prettyPrintData(user3basic)

	user3 := getUserFromBasic(user3basic)

	log.Printf("user3 >")
	prettyPrintData(user3)

	// Update specific fields
	result3 := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "user_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"first_name", "last_name"}),
	}).Create(&user3)

	//updateStringRepForUser(db, user3.UserID)

	log.Printf("Upsert / On Conflict : result.RowsAffected : %v", result3.RowsAffected)

	// ----------------------------------------------------------------------------------------------------

	// Limit and Offset

	log.Printf("---[Limit / Offset]---")

	var partialUsers []User

	// SELECT * FROM users OFFSET 5 LIMIT 3;

	db.Limit(3).Offset(5).Find(&partialUsers)

	for _, user := range partialUsers {
		prettyPrintData(user)
	}

	// ----------------------------------------------------------------------------------------------------

	// query with primary key

	log.Printf("---[Query Users With List of Primary Keys]---")
	users = make([]User, 0)
	db.Where("user_id IN ?", []string{"628555772a8b7b9926ffb917", "6285557743a8bdeb2aa5dc07"}).Find(&users)
	for _, user := range users {
		prettyPrintData(user)
	}

	// ----------------------------------------------------------------------------------------------------

	// query using structs

	// FYI : When querying with struct, GORM will only query with non-zero fields, that means if your field’s
	// value is 0, '', false or other zero values, it won’t be used to build query conditions,

	// for example:
	// 		db.Where(&User{Name: "jinzhu", Age: 0}).Find(&users)
	//		translates to >
	// 		SELECT * FROM users WHERE name = "jinzhu";

	// To include zero values in the query conditions, you can use a map, which will
	// include all key-values as query conditions, for example:

	// 		db.Where(map[string]interface{}{"Name": "jinzhu", "Age": 0}).Find(&users)
	//		translates to >
	// 		SELECT * FROM users WHERE name = "jinzhu" AND age = 0;

	log.Printf("---[Query Using structs]---")
	var searchData User
	var user User
	searchData.Email = "sonialivingston@hinway.com"
	db.Where(&searchData).First(&user)
	prettyPrintData(user)

	// ----------------------------------------------------------------------------------------------------

	// query using maps

	log.Printf("---[Query Using maps]---")

	users = make([]User, 0)
	db.Where(map[string]interface{}{"first_name": "Stacy", "last_name": "Mason"}).Find(&users)
	for _, user := range users {
		prettyPrintData(user)
	}

	// ----------------------------------------------------------------------------------------------------

	// get columne names for model

	log.Printf("---[Column names for 'User']---")

	columnNames := getColumnNamesForModel(db, &User{})

	prettyPrintData(columnNames)

	// ----------------------------------------------------------------------------------------------------

	log.Printf("---[Non Exact Search | Query Using 'string_rep' column | AND search]---")
	searchStrings := []string{"772", "none.com"}
	users = make([]User, 0)
	searchType := SearchAND
	sqlQuery, err := getSQLQueryForNonExactPatternSearch(searchStrings, searchType)
	if err != nil {
		log.Printf("error : %v", err.Error())
		return
	}
	db.Where(sqlQuery).Find(&users)
	for _, user := range users {
		prettyPrintData(user)
	}
	if len(users) == 0 {
		log.Printf("no record found for search type ( %v ) and search string ( %v )", searchType, searchStrings)
	}

	// ----------------------------------------------------------------------------------------------------

	log.Printf("---[Non Exact Search | Query Using 'string_rep' column | OR search]---")
	searchStrings = []string{"Marisol", "Davidson"}
	users = make([]User, 0)
	searchType = SearchOR
	sqlQuery, err = getSQLQueryForNonExactPatternSearch(searchStrings, searchType)
	if err != nil {
		log.Printf("error : %v", err.Error())
		return
	}
	db.Where(sqlQuery).Find(&users)
	for _, user := range users {
		prettyPrintData(user)
	}
	if len(users) == 0 {
		log.Printf("no record found for search type ( %v ) and search string ( %v )", searchType, searchStrings)
	}

	// ----------------------------------------------------------------------------------------------------

	log.Printf("---[Exact Search | Query For Search | OR]---")
	searchStrings = []string{"Marisol", "Davidson", "466-3255", "62855577fc3729572a693d79", "62855577fc3", "DONAcampos@hinway.COM"}
	users = make([]User, 0)

	//users, err = getRecordsForExactSearchOR(db, searchStrings)
	//if err != nil {
	//	log.Printf("error : %v", err.Error())
	//	return
	//}

	searchType = SearchOR
	sqlQuery, err = getSQLQueryForExactSearch(searchStrings, searchType)
	if err != nil {
		log.Printf("error : %v", err.Error())
		return
	}
	db.Where(sqlQuery).Find(&users)

	for _, user := range users {
		prettyPrintData(user)
	}
	if len(users) == 0 {
		log.Printf("no record found for exact search for search strings ( %v )", searchStrings)
	}

	// ----------------------------------------------------------------------------------------------------

	log.Printf("---[Exact Search | Query For Search | AND]---")
	searchStrings = []string{"Wendy", "wendylawson@hinway2.com", "628555772a8b7b9926ffb919"}
	users = make([]User, 0)

	//users, err = getRecordsForExactSearchAND(db, searchStrings)
	//if err != nil {
	//	log.Printf("error : %v", err.Error())
	//	return
	//}

	searchType = SearchAND
	sqlQuery, err = getSQLQueryForExactSearch(searchStrings, searchType)
	if err != nil {
		log.Printf("error : %v", err.Error())
		return
	}
	db.Where(sqlQuery).Find(&users)
	for _, user := range users {
		prettyPrintData(user)
	}
	if len(users) == 0 {
		log.Printf("no record found for exact search for search strings ( %v )", searchStrings)
	}

	// ----------------------------------------------------------------------------------------------------

}

/*
Leverage '#' -> which is the delimiter

Case Insensitive , Non-Exact Search (Pattern Matching)

testdb=> select * from user_records where string_rep ~* '6285557711c';
         user_id          | first_name | last_name |         email         |       phone       | active |  balance  |                                          string_rep
--------------------------+------------+-----------+-----------------------+-------------------+--------+-----------+-----------------------------------------------------------------------------------------------
 6285557711c7ed18cf923d17 | Roslyn     | Owen      | roslynowen@hinway.com | +1 (841) 468-3975 | t      | $2,071.96 | #6285557711c7ed18cf923d17#Roslyn#Owen#roslynowen@hinway.com#+1 (841) 468-3975#true#$2,071.96#
(1 row)

Case Insensitive , Exact String Search

testdb=> select * from user_records where string_rep ~* '#6285557711c#';
 user_id | first_name | last_name | email | phone | active | balance | string_rep
---------+------------+-----------+-------+-------+--------+---------+------------
(0 rows)

Case Insensitive , Exact String Search (AND)

testdb=> select * from user_records where string_rep ~* '#miles#' AND string_rep ~* '#bond#';
         user_id          | first_name | last_name |        email         |       phone       | active |  balance  |                                         string_rep
--------------------------+------------+-----------+----------------------+-------------------+--------+-----------+---------------------------------------------------------------------------------------------
 62855577b919c5002fe856a0 | Miles      | Bond      | milesbond@hinway.com | +1 (822) 480-2450 | t      | $2,028.84 | #62855577b919c5002fe856a0#Miles#Bond#milesbond@hinway.com#+1 (822) 480-2450#true#$2,028.84#
(1 row)

Case Insensitive , Exact String Search (AND)

testdb=> select * from user_records where string_rep ~* '#miles#' AND string_rep ~* '#bond1#';
 user_id | first_name | last_name | email | phone | active | balance | string_rep
---------+------------+-----------+-------+-------+--------+---------+------------
(0 rows)

*/

/*
getRecordsForExactSearchOR  => Is replaced with getSQLQueryForExactSearch

getRecordsForExactSearchOR : Is deprecated, but still can be used
*/
func getRecordsForExactSearchOR(db *gorm.DB, searchStrings []string) ([]User, error) {
	users := make([]User, 0)
	for _, searchString := range searchStrings {
		userList := make([]User, 0)
		lowerCaseSearchString := strings.ToLower(searchString)
		db.Where("LOWER(user_id) = LOWER(?)", lowerCaseSearchString).
			Or("LOWER(first_name) = LOWER(?)", lowerCaseSearchString).
			Or("LOWER(last_name) = LOWER(?)", lowerCaseSearchString).
			Or("LOWER(phone) = LOWER(?)", lowerCaseSearchString).
			Or("LOWER(email) = LOWER(?)", lowerCaseSearchString).Find(&userList)
		for _, user := range userList {
			users = append(users, user)
		}
	}
	return users, nil
}

/*
getRecordsForExactSearchAND  => Is replaced with getSQLQueryForExactSearch

getRecordsForExactSearchAND : Is deprecated, but still can be used
*/
func getRecordsForExactSearchAND(db *gorm.DB, searchStrings []string) ([]User, error) {
	users := make([]User, 0)

	for _, searchString := range searchStrings {
		lowerCaseSearchString := strings.ToLower(searchString)
		db = db.Debug()
		db = db.Where(db.Where("LOWER(user_id) = LOWER(?)", lowerCaseSearchString).
			Or("LOWER(first_name) = LOWER(?)", lowerCaseSearchString).
			Or("LOWER(last_name) = LOWER(?)", lowerCaseSearchString).
			Or("LOWER(phone) = LOWER(?)", lowerCaseSearchString).
			Or("LOWER(email) = LOWER(?)", lowerCaseSearchString))

	}

	err := db.Find(&users).Error
	if err != nil {
		log.Printf("error : %v", err.Error())
		return users, err
	}
	return users, nil
}

// FYI : ~* makes it case-insensitive search

func getSQLQueryForNonExactPatternSearch(searchStrings []string, searchType Search) (string, error) {
	sqlQuery := ""
	counter := 1

	if len(searchStrings) == 0 {
		return sqlQuery, errors.New("length of searchStrings is 0 , please provide valid list")
	}

	if len(searchStrings) == 1 {
		searchType = SearchSingle
	}

	for _, searchString := range searchStrings {
		if searchType == SearchAND {
			if counter == len(searchStrings) {
				sqlQuery += fmt.Sprintf(" string_rep ~* ('%v') ", searchString)
			} else {
				sqlQuery += fmt.Sprintf(" string_rep ~* ('%v') AND ", searchString)
			}
			counter++
		} else if searchType == SearchOR {
			if counter == len(searchStrings) {
				sqlQuery += fmt.Sprintf(" string_rep ~* ('%v') ", searchString)
			} else {
				sqlQuery += fmt.Sprintf(" string_rep ~* ('%v') OR ", searchString)
			}
			counter++
		} else if searchType == SearchSingle {
			sqlQuery = fmt.Sprintf(" string_rep ~* ('%v') ", searchStrings[0])
		} else {
			return sqlQuery, errors.New("please provide valid search type")
		}
	}

	log.Printf("getSQLQueryForNonExactPatternSearch : sqlQuery >>")
	log.Printf(sqlQuery)

	return sqlQuery, nil
}

// FYI : ~* makes it case-insensitive search
//     : string_rep ~* ('#%v#') -> makes it exact search (case insensitive)

func getSQLQueryForExactSearch(searchStrings []string, searchType Search) (string, error) {
	sqlQuery := ""
	counter := 1

	if len(searchStrings) == 0 {
		return sqlQuery, errors.New("length of searchStrings is 0 , please provide valid list")
	}

	if len(searchStrings) == 1 {
		searchType = SearchSingle
	}

	for _, searchString := range searchStrings {
		if searchType == SearchAND {
			if counter == len(searchStrings) {
				sqlQuery += fmt.Sprintf(" string_rep ~* ('#%v#') ", searchString)
			} else {
				sqlQuery += fmt.Sprintf(" string_rep ~* ('#%v#') AND ", searchString)
			}
			counter++
		} else if searchType == SearchOR {
			if counter == len(searchStrings) {
				sqlQuery += fmt.Sprintf(" string_rep ~* ('#%v#') ", searchString)
			} else {
				sqlQuery += fmt.Sprintf(" string_rep ~* ('#%v#') OR ", searchString)
			}
			counter++
		} else if searchType == SearchSingle {
			sqlQuery = fmt.Sprintf(" string_rep ~* ('#%v#') ", searchStrings[0])
		} else {
			return sqlQuery, errors.New("please provide valid search type")
		}
	}

	log.Printf("getSQLQueryForExactSearch : sqlQuery >>")
	log.Printf(sqlQuery)

	return sqlQuery, nil
}

func (u User) AfterCreate(db *gorm.DB) (err error) {
	log.Printf(">>>> [ AfterCreate() ] <<<<")
	stringRep := getStringRep(u.UserBasic)
	log.Printf("AfterCreate : u.UserID : %v", u.UserID)
	u.StringRep = stringRep
	db.Model(u).Save(u)
	return nil
}

func prettyPrintData(data interface{}) {
	dataBytes, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Printf("error : could not MarshalIndent json : %v", err.Error())
		return
	}
	fmt.Printf("\n%v\n\n", string(dataBytes))
}

func getColumnNamesForModel(db *gorm.DB, myModel interface{}) []string {
	columnNames := make([]string, 0)
	result, _ := db.Debug().Migrator().ColumnTypes(&myModel)
	for _, v := range result {
		columnNames = append(columnNames, v.Name())
	}
	return columnNames
}
