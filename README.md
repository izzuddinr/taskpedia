# taskpedia

List URL:

- View User (GET)
http://localhost:18080/user/view

- User Create (POST)
http://localhost:18080/user/create
  
  Request Message Example:
  {
    "Name": "123"
  }
  
- User Update (POST)
http://localhost:18080/user/update

  Request Message Example:
  {
    "ID": 11,
    "Name": "123"
  }
  
  
  
- View Task (GET)
http://localhost:18080/task/view

- Task Create (POST)
http://localhost:18080/task/create
  
  Request Message Example:
  {
      "Name": "TestName123123123",
      "Desc": "hgjghjghjfgdf",
      "UserID": 6,
      "Username": "Test123123",
      "Status": "On Progress"
  }
  
- Task Update (POST)
http://localhost:18080/task/update

  Request Message Example:
  {
      "ID": 14,
      "Name": "TestName",
      "Desc": "Lorem Ipsum Dolores Ciao Bella",
      "UserID": 2,
      "Username": "Test123123",
      "Status": "On Review"
  }
  
- Stat (GET)
http://localhost:18080/stat

- Search (GET)
http://localhost:18080/task/search/userid
  
  Request Message Example:
  {
    "UserID": 3
  }
