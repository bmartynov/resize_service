# Image resizer service

service for image resizing with cache

### Run
docker-compose up

### Requests

**Add download**
----
* **URL**
  /resize/

* **Method:**
    GET

*  **URL Params**

   **Required:**

   `uril=[string]&size=[string]`

* **Success Response:**

  * **Code:** 200 <br />
    **Content:** `image`

* **Example**
```http://127.0.0.1:8080/resize/?size=400x400&url=https://pp.userapi.com/c845020/v845020995/12a70a/KtmQ7Yb-efc.jpg```

