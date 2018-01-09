# yamc Server

Based on [gin](https://github.com/gin-gonic/gin) web framework. Uses [yamc.Store](https://github.com/someanon/yamc/tree/master/store) as store backend. Has [go client](https://github.com/someanon/yamc/tree/master/client). 

# API documentation
All methods require [HTTP Basic Authorization](https://en.wikipedia.org/wiki/Basic_access_authentication).

## Get key
Return value by key

* **Path:** `/key`

* **Method:** ` GET`

*  **URL Params**

    **Required:**
   
   `key=[string]`

* **Success Response:**
  
    * **Code:** 200 OK <br />
    **Content:** `value`

* **Error Response:**

    * **Code:** 400 Bad request <br />
    **Reason:** absent key
    
    * **Code:** 401 Unauthorized <br />
      **Reason:** absent or wrong authorization header
    
    * **Code:** 404 Not found <br />
    **Reason:** key not found or not scalar type
    
  * **Code:** 500 Internal server error

* **Sample Call:**

    `curl -u test:test -X GET "http://127.0.0.1/key?key=k"`
  
## Set key
Set key to value

* **Path:** `/key`

* **Method:** `PUT`
  
*  **URL Params** 

    **Required:**
   
   `key=[string]`
   
   `ttl=[time.Duration string]`  
   
* **Data Params**

    Value string

* **Success Response:**

    * **Code:** 200 OK <br/>
 
* **Error Response:**

    * **Code:** 400 Bad request <br />
    **Reason:** absent or invalid key or ttl
    
    * **Code:** 401 Unauthorized <br />
      **Reason:** absent or wrong authorization header
    
    * **Code:** 500 Internal server error

* **Sample Call:**

    `curl -u test:test -X PUT -d "value" "http://127.0.0.1/key?key=k&ttl=60s"`
  
## Get list item
Get list value by key and index 

* **Path:** `/list`

* **Method:** `GET`

*  **URL Params**

    **Required:**
 
    `key=[string]`
 
    `index=[unsigned integer]`

* **Success Response:**

    * **Code:** 200 OK <br />
    **Content:** `value`

* **Error Response:**

    * **Code:** 400 Bad request <br />
    **Reason:** absent or invalid key or index
    
    * **Code:** 401 Unauthorized <br />
      **Reason:** absent or wrong authorization header
  
    * **Code:** 404 Not found <br />
    **Reason:** list not found; not list item; index in list not exists (too big index or empty list)
  
    * **Code:** 500 Internal server error

* **Sample Call:**

    `curl -u test:test -X GET "http://127.0.0.1/list?key=k&index=0"`

## Set list
Set key to list

* **Path:** `/list`

* **Method:** `PUT`

*  **URL Params** 

    **Required:**
 
    `key=[string]`
    
    `ttl=[time.Duration string]`  
 
* **Data Params**

    YAML encoded list

* **Success Response:**

    * **Code:** 200 OK <br/>

* **Error Response:**

    * **Code:** 400 Bad request <br />
    **Reason:** absent or invalid key or ttl; invalid list YAML
    
    * **Code:** 401 Unauthorized <br />
      **Reason:** absent or wrong authorization header
  
    * **Code:** 500 Internal server error

* **Sample Call:**

    `curl -u test:test -X PUT -d $"- a\n- b\n" "http://127.0.0.1/list?key=k&ttl=60s"`

## Get dictionary item
Get dictionary value by key and dictionary key

* **Path:** `/dict`

* **Method:** `GET`

*  **URL Params**

    **Required:**
 
    `key=[string]`
 
    `dkey=[string]` (dictionary key)

* **Success Response:**

    * **Code:** 200 OK <br />
    **Content:** `value`

* **Error Response:**

    * **Code:** 400 Bad request <br />
    **Reason:** absent key or dkey
    
    * **Code:** 401 Unauthorized <br />
      **Reason:** absent or wrong authorization header
  
    * **Code:** 404 Not found <br />
    **Reason:** dictionary not found; dkey in dictionary not found
  
    * **Code:** 500 Internal server error

* **Sample Call:**

    `curl -u test:test -X GET "http://127.0.0.1/dict?key=k&dkey=dk"`

## Set dictionary
Set key to dictionary

* **Path:** `/dict`

* **Method:** `PUT`

*  **URL Params** 

    **Required:**
 
    `key=[string]`
     
    `ttl=[time.Duration string]`  
 
* **Data Params**

    YAML encoded dictionary

* **Success Response:**
    
    * **Code:** 200 OK

* **Error Response:**

    * **Code:** 400 Bad request <br />
    **Reason:** absent or invalid key or ttl; invalid dict's YAML
    
    * **Code:** 401 Unauthorized <br />
      **Reason:** absent or wrong authorization header
  
    * **Code:** 500 Internal server error

* **Sample Call:**

    `curl -u test:test -X PUT -d $"a: b\nc: d\n" "http://127.0.0.1/dict?key=k&ttl=60s"`

## Remove key
Remove value or list or dictionary. Paths are not bound to the type. Any path removes any key type: value, list or dict.

* **Path:** `/key` or `/list` or `/dict`

* **Method:** `DELETE`

*  **URL Params**

   **Required:**
   
   `key=[string]`

* **Success Response:**
    Value (or list or dict) is deleted
  
    * **Code:** 200 OK <br />

* **Error Response:**

    * **Code:** 400 Bad request <br />
    **Reason:** absent key
    
    * **Code:** 401 Unauthorized <br />
      **Reason:** absent or wrong authorization header

* **Sample Call:**

    `curl -u test:test -X DELETE "http://127.0.0.1/key?key=k"`
    
## Get keys
Get all keys list encoded in YAML

* **Path:** `/keys`

* **Method:** `GET`

* **Success Response:**
  
    * **Code:** 200 OK <br />
    **Content:** YAML encoded list

* **Error Response:**
    
    * **Code:** 401 Unauthorized <br />
      **Reason:** absent or wrong authorization header

    * **Code:** 500 Internal server error

* **Sample Call:**

    `curl -u test:test -X GET "http://127.0.0.1/keys"`