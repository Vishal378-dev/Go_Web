


// api's request and urls

http://localhost:8000/user

{
  "name":"G Yamamoto",
  "email":"g.yamamoto@mail.com",
  "phone":"7829310923",
  "password":"gyamamoto"
}

http://localhost:8000/hotels

{
  "name":"hotel inn",
  "description":"htel",
   "star":4,
   "review":[],
   "address":{
     "landmark":"gol chown",
     "city":"noida",
     "state":"ncr",
     "pincode":402863,
     "coordinates":{
       "longitude":43.4534,
       "latitude":34.932908
     }
   },
   "rooms":[],  
   "amenities":{"garden":true},
   "typesofrooms":["luxury","economy"]
   
   
}

http://localhost:8000/user/login

// {
//   "email":"misaki.kurosaki@mail.com",
//   "password":"misakikurosaki"
// }


{
  "email":"g.yamamoto@mail.com",
  "password":"gyamamoto"
}

http://localhost:8000/hotel/67dd9b331eabe38bd4803c5e?isRoom=true


http://localhost:8000/room

{
 "class" :"suite",
 "roomnumber":524,
 "isbooked":false,
 "features":"ac,food,tennis,private pool",
 "hotelid": "67dd9b331eabe38bd4803c5e",
"roomcategory":"single"
}


http://localhost:8000/account?id=67debee3dc0647d58e9db5a1

{
  "bankname":"Axis Bank",
  "accountnumber":903481782,
  "bankifsc":"utib0000004",
  "bankholderfirstname":"radha",
  "bankholderlastname":"mohan",
  "balance":100,
  "userid":"67cd4fdfe107cfeb6e8ba9cd"
}