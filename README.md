## 1. 📡 URL: `http://localhost:8888/account`

### Purpose
- Store the account details. (Add,Modify,Delete and Get)

### Add Account: (POST method)

```bash
Request :
{
    "accid":"ACC2",
    "accname":"James",
    "email":"james007@gmail.com"
}

### ✅ Success Response :
{
    "sts":"S",
    "emsg":""
}

### ❌ Error Response :
{
    "sts": "E",
    "msg": "AACH03 : Error 1062 (23000): Duplicate entry 'james007@gmail.com' for key 'account_holders.email'"
}
```

### GET Account: (GET method)

```bash
Request :
{
    "accid":"ACC2"
}

### ✅ Success Response :
{
    "email": "james007@gmail.com",
    "accname": "James",
    "sts": "S",
    "msg": ""
}

### ❌ Error Response :
{
    "sts": "E",
    "msg": "AACH03 : Account not available - ACC3"
}
```

### Modify Account: (PUT method)

```bash
Request :
{
    "accid":"ACC2",
    "accname":"James",
    "email":"james777@gmail.com"
}

### ✅ Success Response :
{
    "sts":"S",
    "emsg":""
}

### ❌ Error Response :
{
    "sts": "E",
    "msg": "AACH03 : Syntax error"
}
```

### Delete Account: (DELETE method)

- It will also delete all the destionations for the account.

```bash
Request :
{
    "accid":"ACC2"
}

### ✅ Success Response :
{
    "sts":"S",
    "emsg":""
}

### ❌ Error Response :
{
    "sts": "E",
    "msg": "AACH03 : Unable to remove, Please provide valide Account ID - ACC2"
}
```

---

## 2. 📡 URL: `http://localhost:8888/destination`

### Purpose
- Store the destinations of the account. (Add,Modify,Delete and Get)

### Add Destination : (POST method)

```bash
Request :
{
  "account_id": "ACC2",
  "des_id":"DES1",
  "url": "https://destination.com/api3",
  "http_method": "GET",
  "headers": {
    "APP_ID": "1234APPID1234",
    "APP_SECRET": "xxxx",
    "ACTION": "user.update",
    "Content-Type": "application/json"
  }
}

### ✅ Success Response :
{
    "sts":"S",
    "emsg":""
}

### ❌ Error Response :
{
    "sts": "E",
    "msg": "DDT03 : URL already exist for the account - https://destination.com/api3"
}
```

### Modify Destination : (PUT method)

```bash
Request :
{
  "account_id": "ACC2",
  "des_id":"DES1",
  "url": "https://destination.com/api1",
  "http_method": "PUT",
  "headers": {
    "APP_ID": "1234APPID1234",
    "APP_SECRET": "xxxx",
    "ACTION": "user.update",
    "Content-Type": "application/json"
  }
}

### ✅ Success Response :
{
    "sts":"S",
    "emsg":""
}

### ❌ Error Response :
{
    "sts": "E",
    "msg": "DDT03 : syntax error"
}
```

### GET Account: (GET method)

```bash
Request :
{
    "des_id":"DES1"
}

### ✅ Success Response :
{
    "account_id": "ACC2",
    "des_id": "DES1",
    "url": "https://destination.com/api",
    "http_method": "GET",
    "headers": {
        "ACTION": "user.update",
        "APP_ID": "1234APPID1234",
        "APP_SECRET": "xxxx",
        "Content-Type": "application/json"
    },
    "sts": "S",
    "msg": ""
}

### ❌ Error Response :
{
    "sts": "E",
    "msg": "DDT03 : syntax error"
}
```

### Delete Account: (DELETE method)

```bash
Request :
{
  "des_id":"DES1"
}

### ✅ Success Response :
{
    "sts":"S",
    "emsg":""
}

### ❌ Error Response :
{
    "sts": "E",
    "msg": "DDT03 : Unable to delete, Destination ID not available - DES1"
}
```

---

## 3. 📡 URL: `http://localhost:8888/getaccdest`

### Description
This endpoint loads CSV data into the database by internally calling the `ReadConsStore()` function.

### GET Destionation for an account: (GET method)

```bash
Header : [{"key":"ACC_ID","value":"ACC2"}]

### ✅ Success Response :
{
    "data": [
        {
            "des_id": "DES1",
            "url": "https://destination.com/api",
            "http_method": "GET"
        }
    ],
    "sts": "S",
    "msg": ""
}

### ❌ Error Response :
{
    "data": null,
    "sts": "E",
    "msg": "DGAD02 : Destinations are not available for the account - ACC1"
}
```

---

## 4. 📡 URL: `http://localhost:8888/server/incoming_data`

### Purpose
- Identify the destinations of the respective app_secret.
- And then construct the request with the incoming json and call the destionation api.

### POST Destionation for an account: (POST method)

```bash
Header : [{"key":"App_Secret","value":"QUNDMkphbWVzamFtZXM3NzdAZ21haWwuY29t"}]

### any json request you can provide
Body :
{
    "fdt":"2024-01-03",
    "tdt":"2024-05-23"
}

### ✅ Success Response :
{
    "data": [
        {
            "url": "https://destination.com/api",
            "method": "PUT",
            "jbody": "{\"fdt\":\"2024-01-03\",\"tdt\":\"2024-05-23\"}",
            "headers": "{\"ACTION\":\"user.update\",\"APP_ID\":\"1234APPID1234\",\"APP_SECRET\":\"xxxx\",\"Content-Type\":\"application/json\"}"
        },
        {
            "url": "https://destination.com/api3?fdt=2024-01-03&tdt=2024-05-23",
            "method": "GET",
            "jbody": "",
            "headers": "{\"ACTION\":\"user.update\",\"APP_ID\":\"1234APPID1234\",\"APP_SECRET\":\"xxxx\",\"Content-Type\":\"application/json\"}"
        }
    ],
    "sts": "S",
    "msg": "Data received successfully"
}

### ❌ 1. Error Response :
{
    "data": null,
    "sts": "E",
    "msg": "Invalid Data"
}

### ❌ 2. Error Response :
{
    "data": null,
    "sts": "E",
    "msg": "Un Authenticate"
}

### ❌ 3. Error Response :
{
    "data": null,
    "sts": "E",
    "msg": "Destionations not available for the secret (or) Invalid secret"
}


```

---

## 5. Database Tables:

```bash
-- Table for account deaills
CREATE TABLE `account_holders` (
  `id` int NOT NULL AUTO_INCREMENT,
  `accid` varchar(100) NOT NULL,
  `accname` varchar(100) NOT NULL,
  `email` varchar(100) NOT NULL,
  `appsecret` varchar(300) NOT NULL,
  `createdby` varchar(100) NOT NULL,
  `createddate` datetime NOT NULL,
  `updatedby` varchar(100) DEFAULT NULL,
  `updateddate` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `email` (`email`),
  UNIQUE KEY `accid` (`accid`)
);


-- Table for destinations
CREATE TABLE `destination` (
  `id` int NOT NULL AUTO_INCREMENT,
  `desId` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `accid` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
  `url` varchar(100) NOT NULL,
  `method` varchar(10) NOT NULL,
  `headers` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci,
  `createdby` varchar(100) NOT NULL,
  `createddate` datetime NOT NULL,
  `updatedby` varchar(100) DEFAULT NULL,
  `updateddate` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `destination_unique` (`desId`),
  KEY `destination_accid_IDX` (`accid`) USING BTREE
);
```

- Maria DB used.