# Library REST API

**Get User**
----
  Returns json data about a single user. Requires ADMIN role.

* **URL**

  library/api/v1/users/:email

* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
   `email=[string]`

   **Headers** `Authorization: Bearer jwt_token`

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** `{ Email : "email@gmail.com", Taken_books : [], Returned_books : [] }`
 
* **Error Response:**

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error message : "User not found" }`

  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `"You need to be authorized to access this route"`

* **Sample Call:**

  ```go
    http.NewRequest("GET", "library/api/v1/users/email@gmail.com", nil)
  ```


**Get All users**
----
  Returns json data all users. Requires ADMIN role.

* **URL**

  library/api/v1/users

* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
    **Headers** `Authorization: Bearer jwt_token`

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** `[{ Email : "email@gmail.com", Taken_books : [], Returned_books : []] },`
                 `{ Email : "email@abv.bg", Taken_books : [], Returned_books : [] }]`
 
* **Error Response:**
  
  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `"You need to be authorized to access this route"`

* **Sample Call:**

  ```go
    http.NewRequest("GET", "library/api/v1/users", nil)
  ```


**Take book by user**
----
  Returns successful message if book is taken by the user. Requires USER role.

* **URL**

  library/api/v1/users/:email/:isbn

* **Method:**

  `POST`
  
*  **URL Params**

   **Required:**
 
   `email=[string]`
   `isbn=[string]`

  **Headers** `Authorization: Bearer jwt_token`

* **Success Response:**

  * **Code:** 201 <br />
    **Content:** `{ message : Book successfully taken }`
 
* **Error Response:**

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error message : "Book not found" }`

  OR

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error message : "User not found" }`

  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `"You need to be authorized to access this route"`

  OR

  * **Code:** 400 BAD REQUEST <br />
   **Content:** `{ error message : "This book has no available copies" }`

  OR

  * **Code:** 400 BAD REQUEST <br />
   **Content:** `{ error message : "This book is already taken" }`

* **Sample Call:**

  ```go
    http.NewRequest("POST", "library/api/v1/users/:email/:isbn", nil)
  ```


**Return book**
----
  Returns successful message if the book returned. Requires USER role

* **URL**

  library/api/v1/users/:email/:isbn

* **Method:**

  `DELETE`
  
*  **URL Params**

   **Required:**
 
   `email=[string]`
   `isbn=[string]`

  **Headers** `Authorization: Bearer jwt_token`

* **Success Response:**

  * **Code:** 204 <br />
    **Content:** `{ message : Book successfully returned }`
 
* **Error Response:**

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error message : "Book not found" }`

  OR

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error message : "User not found" }`

  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `"You need to be authorized to access this route"`

  OR

  * **Code:** 500 INTERNAL SERVER ERROR <br />
   **Content:** `{ error message : "unable to return book" }`

  OR

  * **Code:** 400 BAD REQUEST <br />
   **Content:** `{ error message : "This book is already taken" }`

* **Sample Call:**

  ```go
    http.NewRequest("DELETE", "library/api/v1/users/:email/:isbn", nil)
  ```

**Get All Books**
----
  Returns json data about a all books in the library. Requires ADMIN or USER role.

* **URL**

  library/api/v1/books

* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
  **Headers** `Authorization: Bearer jwt_token`

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** `[{Isbn: "12345", Title: "Sample Title", Author: "Sample Author", "Available_units: 12"}]`
 
* **Error Response:**
  
  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `"You need to be authorized to access this route"`

* **Sample Call:**

  ```go
    http.NewRequest("GET", "library/api/v1/books", nil)
  ```


**Get Book**
----
  Returns json data about a single book in the library. Requires ADMIN or USER role.

* **URL**

  library/api/v1/books/:isbn

* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
  **Headers** `Authorization: Bearer jwt_token`

* **Success Response:**

  * **Code:** 201 <br />
    **Content:** `{message: "Successfully saved"}`
 
* **Error Response:**

  * **Code:** 409 CONFLICT <br />
    **Content:** `{ error message : "Every book must have an unique isbn" }`

  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `"You need to be authorized to access this route"`

* **Sample Call:**

  ```go
    http.NewRequest("GET", "library/api/v1/books/:isbn", nil)
  ``` 

**Save Book**
----
  Saves book in the library. Requires ADMIN role.

* **URL**

  library/api/v1/books

* **Method:**

  `POST`
  
*  **URL Params**

   **Required:**
 
   **Request body** `{Isbn: "12345", Title: "Sample Title", Author: "Sample Author", "Available_units: 12"}`

   **Headers** `Authorization: Bearer jwt_token`

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** `{Isbn: "12345", Title: "Sample Title", Author: "Sample Author", "Available_units: 12"}`
 
* **Error Response:**

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error message : "Book not found" }`

  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `"You need to be authorized to access this route"`

* **Sample Call:**

  ```go
    http.NewRequest("POST", "library/api/v1/books", sampleBody)
  ``` 


**Delete Book**
----
  Returns json data about a single book in the library. Requires ADMIN role.

* **URL**

  library/api/v1/books

* **Method:**

  `DELETE`
  
*  **URL Params**

   **Required:**
 
   `isbn=[isbn]`

  **Headers** `Authorization: Bearer jwt_token`

* **Success Response:**

  * **Code:** 204 NO CONTENT <br />
 
* **Error Response:**

  * **Code:** 404 NOT FOUND <br />
    **Content:** `{ error message : "Book not found" }`

  OR

  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `"You need to be authorized to access this route"`

* **Sample Call:**

  ```go
    http.NewRequest("DELETE", "library/api/v1/books/:isbn", nil)
  ``` 

**User Register**
---
  Register using an email and password

* **URL**
  
   library/api/v1/users/register

*  **URL Params**

   **Required:**
 
    **Request body** `{Email: "email@gmail.com", Password: "password"}`

* **Success Response:**

  * **Code:** 201 CREATED <br />

* **URL**

  library/api/v1/users/register

* **Sample Call:**

  ```go
    http.NewRequest("POST", "library/api/v1/users/register", sampleBody)
  ``` 
  

**User Login**
----
  Returns a jwt token which can be used to acces the REST API.

* **URL**

  library/api/v1/users/login

* **Method:**

  `POST`
  
*  **URL Params**

   **Required:**
 
    **Request body** `{Email: "email@gmail.com", Password: "password"}`

* **Success Response:**

  * **Code:** 200 OK <br />
 
* **Error Response:**

  * **Code:** 403 FORBIDDEN <br />
    **Content:** `{ error message : "Wrong credentials" }`


* **Sample Call:**

  ```go
    http.NewRequest("POST", "library/api/v1/users/login", sampleBody)
  ``` 

  **User Logout**
----
  Logouts the user and deletes the jwt token with it

* **URL**

  library/api/v1/users/logout

* **Method:**

  `POST`
  
*  **URL Params**

   **Required:**
 
   **Headers** `Authorization: Bearer jwt_token`

* **Success Response:**

  * **Code:** 200 OK <br />
 
* **Error Response:**
  
  * **Code:** 401 UNAUTHORIZED <br />
    **Content:** `{ error message : "Unauthorized" }`

* **Sample Call:**

  ```go
    http.NewRequest("POST", "library/api/v1/users/logout", sampleBody)
  ``` 