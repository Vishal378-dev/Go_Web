


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


http://localhost:8000/booking

{
  "roomid":"67de6ec5261c3aa1580c2fc0",
  "startdate":"25-03-2025 12:00:00",
  "enddate":"26-03-2025 12:00:00"
}





http://localhost:8000/account

{
  "bankname":"HDFC BANK",
  "accountnumber":93840348,
  "bankifsc":"hdfc930282",
  "bankholderfirstname":"Gyenrui",
  "bankholderlastname":"yamamoto",
  "balance":0,
  "userid":"67dc00a9cf105d8ef3e03ebd"
}

