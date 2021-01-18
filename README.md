# Library REST API

**Get User**
----
  Returns json data about a single user.

* **URL**

  library/api/v1/users/:email

* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
   `email=[string]`

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
  Returns json data all users

* **URL**

  library/api/v1/users

* **Method:**

  `GET`
  
*  **URL Params**

   **Required:**
 
   None

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
  Returns successful message if book is taken by the user.

* **URL**

  library/api/v1/users/:email/:isbn

* **Method:**

  `POST`
  
*  **URL Params**

   **Required:**
 
   `email=[string]`
   `isbn=[string]`

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
  Returns successful message if the book returned.

* **URL**

  library/api/v1/users/:email/:isbn

* **Method:**

  `DELETE`
  
*  **URL Params**

   **Required:**
 
   `email=[string]`
   `isbn=[string]`

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
