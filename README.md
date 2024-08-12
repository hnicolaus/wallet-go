# wallet-go
## Setup
To run the project, run the following command: docker-compose up --build.
- HTTP endpoints are accessible at localhost:1323
  - Postgresql DB is accessible at:
  - Host: localhost
  - Port: 5342
  - Database: database
  - Username: postgres
  - Password: postgres
- DB schema and seed data is defined in `database.sql` file.
- OpenAPI specification is defined in `api.yml`
- Please check code comments for details about the implementations.

  
## Testing
1. Login with phone number and password to obtain JWT
   - Endpoint: `POST localhost:1323/v1/user/login`
   - Request body:
     ```
     {
         "phone_number": "+6281122334455",
         "password": "Admin1234!"    
     }
     ```
   - Response header: `Authorization: Bearer abcdefghijklmnopqrstuvwxyz`
   - Response body:
      ```
      {
          "header": {
              "messages": [
                  "request successful"
              ],
              "success": true
          },
          "user": {
              "id": 1
          }
      }
      ```

2. Check account balance
   - Endpoint: `GET localhost:1323/v1/user`
   - Request header: `Authorization: Bearer abcdefghijklmnopqrstuvwxyz`
   - Response body:
      ```
      {
        "header": {
            "messages": [
                "request successful"
            ],
            "success": true
        },
        "user": {
            "balance": 0,
            "full_name": "name1",
            "id": 1,
            "phone_number": "+6281122334455"
        }
      }
      ```

3. Top-up account balance
   - Endpoint: `POST localhost:1323/v1/user/1/transactions`
   - Request header: `Authorization: Bearer abcdefghijklmnopqrstuvwxyz`
   - Request body:
      ```
      {
          "amount": 5000000,
          "type": "TopUp",
          "description": "Top Up"
      }
      ```
    - Response body:
      ```
      {
          "header": {
              "messages": [
                  "request successful"
              ],
              "success": true
          },
          "transaction": {
              "id": "e3b54dbd-d445-4885-854e-918a9b959939"
          }
      }
      ```

4. Transfer to another account
   - Endpoint: `POST localhost:1323/v1/user/1/transactions`
   - Request header: `Authorization: Bearer abcdefghijklmnopqrstuvwxyz`
   - Request body:
      ```
      {
          "amount": 1000000,
          "recipient_id": 2,
          "type": "TransferOut",
          "description": "Traktir Makan",
          "password": "Admin1234!"
      }
      ```
    - Response body:
      ```
      {
        "header": {
            "messages": [
                "request successful"
            ],
            "success": true
        },
        "transaction": {
            "id": "27d68748-0566-49a7-aaf7-833a34a930d5"
        }
      }
      ```
    
5. Check account balance
   - Endpoint: `GET localhost:1323/v1/user`
   - Request header: `Authorization: Bearer abcdefghijklmnopqrstuvwxyz`
   - Response body:
      ```
      {
          "header": {
              "messages": [
                  "request successful"
              ],
              "success": true
          },
          "user": {
              "balance": 4000000,
              "full_name": "name1",
              "id": 1,
              "phone_number": "+6281122334455"
          }
      }
      ```
    
6. Repeat step 1-2 to check balance for userID 2 with the following payload:
     ```
      {
          "phone_number": "+6281122334455",
          "password": "Admin1234!"    
      }
    ```

7. Optional: sign-up as a new User, then repeat steps 1-5 to login and make transactions as the new User.
   - Endpoint: `POST localhost:1323/v1/user`
   - Request body:
      ```
      {
          "full_name": "Test Abc",
          "phone_number": "+628123456769",
          "password": "Admin1234!"
      }
      ```
    - Response body:
      ```
      {
          "header": {
              "messages": [
                  "request successful"
              ],
              "success": true
          },
          "user": {
              "id": 3
          }
      }
      ```
