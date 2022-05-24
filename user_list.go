package main

import (
	"encoding/json"
	"fmt"
	"log"
)

/*

https://json-generator.com/#generate

  [
  '{{repeat(20)}}',
  {
    _id: '{{objectId()}}',
    first_name: '{{firstName()}}',
    last_name: '{{surname()}}',
    email: '{{email()}}',
    phone: '+1 {{phone()}}',
    isActive: '{{bool()}}',
    balance: '{{floating(1000, 4000, 2, "$0,0.00")}}'
  }
]

*/

func GetUserList() []map[string]interface{} {
	userList := fmt.Sprintf(`[
	  {
		"_id": "62855577dd46c9dd7e7e8fab",
		"first_name": "Swanson",
		"last_name": "Baird",
		"email": "swansonbaird@hinway.com",
		"phone": "+1 (802) 428-3754",
		"isActive": false,
		"balance": "$2,014.62"
	  },
	  {
		"_id": "628555774195f8fd24b31c6f",
		"first_name": "Bonnie",
		"last_name": "Little",
		"email": "bonnielittle@hinway.com",
		"phone": "+1 (981) 502-3955",
		"isActive": true,
		"balance": "$1,596.06"
	  },
	  {
		"_id": "62855577c75a2bffe8099c63",
		"first_name": "Browning",
		"last_name": "Travis",
		"email": "browningtravis@hinway.com",
		"phone": "+1 (877) 415-2490",
		"isActive": false,
		"balance": "$3,306.06"
	  },
	  {
		"_id": "62855577b919c5002fe856a0",
		"first_name": "Miles",
		"last_name": "Bond",
		"email": "milesbond@hinway.com",
		"phone": "+1 (822) 480-2450",
		"isActive": true,
		"balance": "$2,028.84"
	  },
	  {
		"_id": "6285557721d7e6f8d5343a9d",
		"first_name": "Collier",
		"last_name": "Cardenas",
		"email": "colliercardenas@hinway.com",
		"phone": "+1 (885) 525-3000",
		"isActive": false,
		"balance": "$1,381.63"
	  },
	  {
		"_id": "6285557711c7ed18cf923d17",
		"first_name": "Roslyn",
		"last_name": "Owen",
		"email": "roslynowen@hinway.com",
		"phone": "+1 (841) 468-3975",
		"isActive": true,
		"balance": "$2,071.96"
	  },
	  {
		"_id": "6285557704966a5df4dfe23d",
		"first_name": "Owens",
		"last_name": "Burns",
		"email": "owensburns@hinway.com",
		"phone": "+1 (955) 567-3607",
		"isActive": false,
		"balance": "$3,503.04"
	  },
	  {
		"_id": "62855577ad64df82b2c481a0",
		"first_name": "Candy",
		"last_name": "Meyers",
		"email": "candymeyers@hinway.com",
		"phone": "+1 (875) 445-2347",
		"isActive": false,
		"balance": "$1,715.14"
	  },
	  {
		"_id": "62855577dbf6b90435ff2107",
		"first_name": "Evans",
		"last_name": "Peterson",
		"email": "evanspeterson@hinway.com",
		"phone": "+1 (816) 484-3296",
		"isActive": true,
		"balance": "$1,852.75"
	  },
	  {
		"_id": "628555771bdb1ab40ffe5d9f",
		"first_name": "Marisol",
		"last_name": "Griffin",
		"email": "marisolgriffin@hinway.com",
		"phone": "+1 (817) 466-3255",
		"isActive": false,
		"balance": "$3,627.11"
	  },
	  {
		"_id": "62855577af1aee8a772c90f0",
		"first_name": "Guthrie",
		"last_name": "Sanders",
		"email": "guthriesanders@hinway.com",
		"phone": "+1 (878) 567-3847",
		"isActive": true,
		"balance": "$1,889.68"
	  },
	  {
		"_id": "6285557774bcfb65bb51c001",
		"first_name": "Pamela",
		"last_name": "Davidson",
		"email": "pameladavidson@hinway.com",
		"phone": "+1 (968) 406-2645",
		"isActive": true,
		"balance": "$2,271.65"
	  },
	  {
		"_id": "62855577dbab16fb2e03dab6",
		"first_name": "Lindsey",
		"last_name": "Whitfield",
		"email": "lindseywhitfield@hinway.com",
		"phone": "+1 (810) 429-3432",
		"isActive": true,
		"balance": "$2,953.70"
	  },
	  {
		"_id": "62855577d812e7b34c7cf6e2",
		"first_name": "Leola",
		"last_name": "Ramsey",
		"email": "leolaramsey@hinway.com",
		"phone": "+1 (842) 524-2799",
		"isActive": true,
		"balance": "$1,033.86"
	  },
	  {
		"_id": "62855577fc3729572a693d79",
		"first_name": "Stacy",
		"last_name": "Mason",
		"email": "stacymason@hinway.com",
		"phone": "+1 (887) 465-2768",
		"isActive": false,
		"balance": "$2,611.62"
	  },
	  {
		"_id": "62855577bba9fd23f6878e63",
		"first_name": "Dona",
		"last_name": "Campos",
		"email": "donacampos@hinway.com",
		"phone": "+1 (894) 598-2963",
		"isActive": false,
		"balance": "$1,917.88"
	  },
	  {
		"_id": "62855577477bdbf14b55b5ad",
		"first_name": "Anna",
		"last_name": "Frederick",
		"email": "annafrederick@hinway.com",
		"phone": "+1 (823) 409-3858",
		"isActive": true,
		"balance": "$3,929.84"
	  },
	  {
		"_id": "62855577cc46f4b32485137f",
		"first_name": "Maynard",
		"last_name": "Howard",
		"email": "maynardhoward@hinway.com",
		"phone": "+1 (951) 533-2249",
		"isActive": true,
		"balance": "$1,323.58"
	  },
	  {
		"_id": "6285557743a8bdeb2aa5dc07",
		"first_name": "Sonia",
		"last_name": "Livingston",
		"email": "sonialivingston@hinway.com",
		"phone": "+1 (957) 570-2414",
		"isActive": false,
		"balance": "$1,174.11"
	  },
	  {
		"_id": "628555772a8b7b9926ffb917",
		"first_name": "Wendy",
		"last_name": "Lawson",
		"email": "wendylawson@hinway.com",
		"phone": "+1 (907) 523-2723",
		"isActive": false,
		"balance": "$1,582.33"
	  },
      {
		"_id": "628555772a8b7b9926ffb918",
		"first_name": "Wendy",
		"last_name": "Lawson000",
		"email": "wendylawson@hinway.com",
		"phone": "+1 (907) 523-2723",
		"isActive": false,
		"balance": "$1,582.33"
	  },
	  {
		"_id": "628555772a8b7b9926ffb919",
		"first_name": "Wendy",
		"last_name": "Lawson000",
		"email": "wendylawson@hinway2.com",
		"phone": "+1 (907) 523-2723",
		"isActive": false,
		"balance": "$1,582.33"
	  }
	]`)

	jsonMap := make([]map[string]interface{}, 0)
	_ = json.Unmarshal([]byte(userList), &jsonMap)
	return jsonMap
}

func GetUserRecords() []User {
	users := make([]User, 0)

	userList := GetUserList()
	for _, user := range userList {
		myUserBasic := UserBasic{
			UserID:    user["_id"].(string),
			FirstName: user["first_name"].(string),
			LastName:  user["last_name"].(string),
			Email:     user["email"].(string),
			Phone:     user["phone"].(string),
			Active:    user["isActive"].(bool),
			Balance:   user["balance"].(string),
		}

		stringRep := getStringRep(myUserBasic)

		log.Printf("stringRep : %v", stringRep)

		myUser := getUserFromBasic(myUserBasic)

		users = append(users, myUser)
	}
	return users
}
