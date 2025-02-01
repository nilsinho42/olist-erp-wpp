# OlistERPMediator
Olist ERP is a cloud service that allows companies to manage their business. It provides an API to access each company data. The OlistERPMediator is a middleware API that allow customers to query and get responses based on Olist ERP available data.
---
# Architecture
## Authentication
- Olist ERP API requires a token for authentication. The main token expires after 4 hours, and a refresh token expires after 24 hours.
- Olist ERP Mediator cannot acquire the initial token; user login is required via the authentication URL to obtain the code.
- The authentication URL redirects its response to the Redirect URI set in Olist ERP Aplicativos.
- The webserver for the Redirect URI captures the generated token and stores it.
- The main application allows providing a new token through command line inputs.

## Token Management
- The main token expires after 4 hours, and the refresh token expires after 24 hours, requiring proper management to maintain connection, as Olist ERP Mediator cannot perform initial authentication.
- Tokens are stored in a PostgreSQL server accessible by all services, considering the application consists of different services possibly running in separate containers.
- To avoid high costs from continuous PostgreSQL server activity, a replication strategy writes the token to both the database and a JSON file, ensuring token availability even if the database is down.
- Tokens are encrypted for security, with the encryption key stored in environmental variables.

## Next Steps
- In Get Orders, add a mapping for situacao between the number given and the corresponding name
- In Get NF, add a mapping for situacao between the number given and the corresponding name
- In Get Orders add features to get orders based on other searches, such as order number, date filter and etc
- Create Get NF
- Valor field in some places (order model for example) is string, but should not be!
- To have a key value store with OLIST ERP ENDPOINTS

## Future Developments
- Perform the refresh token through AWS Functions with a trigger every 8 hours. Then, even when the server is down, the tokens will still refresh. Also, will allow to turn off and turn on the servers based on users activity. For example, weekends, holidays, and other dates where there may be no requests, there would be no reason to keep the host on. Would also have to create a trigger to start the server whenever needed (based on incoming request). Also, it would request one intermediator to get the requests and turn on/off the remaining services on demand. Not sure how to achieve this, btw.

# Testing
1. Get one time authentication code:
One Time Authentication: https://accounts.tiny.com.br/realms/tiny/protocol/openid-connect/auth?client_id=tiny-api-d0b3e47d766765eedc8dabe8f90eba37474257c3-1737037931&redirect_uri=https://6c86-2804-10f8-4356-8600-64eb-2e6e-a771-dbed.ngrok-free.app&scope=openid&response_type=code

2. Code returned after authentication: e042d4e3-44b3-46ba-97f6-33c1d8c1bc39.714667d9-a95d-4729-abac-a0054964d9fb.3dcda8a1-a6ef-4964-adcc-d0a5e1b8eebb

3. Use code to get the access token: curl --location 'https://accounts.tiny.com.br/realms/tiny/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'grant_type=authorization_code' \
--data-urlencode 'client_id=tiny-api-d0b3e47d766765eedc8dabe8f90eba37474257c3-1737037931' \
--data-urlencode 'client_secret=bpSCcOcdxN8tftz4X8vRRceGd1FTguyw' \
--data-urlencode 'redirect_uri=https://6c86-2804-10f8-4356-8600-64eb-2e6e-a771-dbed.ngrok-free.app' \
--data-urlencode 'code=e042d4e3-44b3-46ba-97f6-33c1d8c1bc39.714667d9-a95d-4729-abac-a0054964d9fb.3dcda8a1-a6ef-4964-adcc-d0a5e1b8eebb'

4. Get the access token from the response: 
{"access_token":"ACCESS_TOKEN",
"expires_in":14400,
"refresh_expires_in":86400,
"refresh_token":"REFRESH_TOKEN",
"not-before-policy":0,
"session_state":"7cf40efb-398f-45f3-b94f-1262d6216d45",
"scope":"openid email offline_access"}

5. Store the token to the token repositories:
curl -X PUT "http://localhost:8081/auth?key=eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJnWXk2cDhkQkU0dDBkZkFvU0J4WkJvbDBkYmpTcEF5Z3FpQm1vY3pMcXJVIn0.eyJleHAiOjE3Mzg0MjMyNDAsImlhdCI6MTczODQwODg0MCwiYXV0aF90aW1lIjoxNzM4NDA4ODE2LCJqdGkiOiJkNmY2N2M5MC0wMzVhLTQ2ZTktODliYy05MDBjMWVhOTBhN2QiLCJpc3MiOiJodHRwczovL2FjY291bnRzLnRpbnkuY29tLmJyL3JlYWxtcy90aW55IiwiYXVkIjoidGlueS1hcGkiLCJzdWIiOiJjOGYxYjlhOC0xMzViLTRmMDEtYmQ4Ni1lZjE3YTE5NGZkZjUiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJ0aW55LWFwaS1kMGIzZTQ3ZDc2Njc2NWVlZGM4ZGFiZThmOTBlYmEzNzQ3NDI1N2MzLTE3MzcwMzc5MzEiLCJzZXNzaW9uX3N0YXRlIjoiNzE0NjY3ZDktYTk1ZC00NzI5LWFiYWMtYTAwNTQ5NjRkOWZiIiwic2NvcGUiOiJvcGVuaWQgZW1haWwgb2ZmbGluZV9hY2Nlc3MiLCJzaWQiOiI3MTQ2NjdkOS1hOTVkLTQ3MjktYWJhYy1hMDA1NDk2NGQ5ZmIiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsInJvbGVzIjp7InRpbnktYXBpIjpbIm1hcmNhcy1lc2NyaXRhIiwicGVkaWRvcy1sZWl0dXJhIiwiZXN0b3F1ZS1sZWl0dXJhIiwibm90YXMtZmlzY2Fpcy1lc2NyaXRhIiwiY29udGF0b3MtZXhjbHVzYW8iLCJpbmZvLWNvbnRhLWxlaXR1cmEiLCJwcm9kdXRvcy1leGNsdXNhbyIsImxpc3RhLXByZWNvcy1sZWl0dXJhIiwiY2F0ZWdvcmlhcy1sZWl0dXJhIiwiZXhwZWRpY2FvLWV4Y2x1c2FvIiwiZm9ybWEtZW52aW8tbGVpdHVyYSIsInNlcGFyYWNhby1sZWl0dXJhIiwicHJvZHV0b3MtbGVpdHVyYSIsImV4cGVkaWNhby1sZWl0dXJhIiwib3JkZW0tY29tcHJhLWxlaXR1cmEiLCJjb250YXRvcy1lc2NyaXRhIiwibm90YXMtZmlzY2Fpcy1sZWl0dXJhIiwicGVkaWRvcy1lc2NyaXRhIiwiY29udGFzLXBhZ2FyLWVzY3JpdGEiLCJvcmRlbS1zZXJ2aWNvLWxlaXR1cmEiLCJjb250YXMtcGFnYXItbGVpdHVyYSIsImV4cGVkaWNhby1lc2NyaXRhIiwiY29udGFzLXJlY2ViZXItbGVpdHVyYSIsInNlcGFyYWNhby1lc2NyaXRhIiwibWFyY2FzLWxlaXR1cmEiLCJmb3JtYS1wYWdhbWVudG8tbGVpdHVyYSIsImludGVybWVkaWFkb3Jlcy1sZWl0dXJhIiwiZ2F0aWxob3MiLCJlc3RvcXVlLWVzY3JpdGEiLCJwcm9kdXRvcy1lc2NyaXRhIiwib3JkZW0tc2Vydmljby1lc2NyaXRhIiwiY29udGFzLXJlY2ViZXItZXNjcml0YSIsImNvbnRhdG9zLWxlaXR1cmEiLCJwZWRpZG9zLWV4Y2x1c2FvIiwib3JkZW0tY29tcHJhLWVzY3JpdGEiLCJub3Rhcy1maXNjYWlzLWV4Y2x1c2FvIl19LCJlbWFpbCI6ImNvbWVyY2lhbEBzZWZmZXIuY29tLmJyIn0.YZLaMbkBsQKSKny0I6l8tgTXZV2GgxkbgfObhTiPw4MaqJ2Wc8NIvOHgZIQb_RBgnECNIqiAP5dtNtZGbE6Gm40Yyz_-aSleACR4XCu0wfq_lHRC0vNih8j-bRRLMvpC-Zmcoc93h7SMTs97H1voY19TAMMXv5QTknWEGJDm4i6KRfb1RKQ3P7d-hHFmD_sdZMDYh4cGgdjOvRsX5N39hNwfFX2kznun0Z3KWnFzPj1HEx9YbubM-OxyOaou_gd9mq4T-Ppx_zlDAYUJUxw6tfa3DvJK_5VIvBpPnkJ-JkjwWA_X6jqfsWNGc0zT5OBRKfjkFoBW7szJpWyQ-jiejg" 

6. Test call to the OLIST ERP API:
curl -X GET "https://api.tiny.com.br/public-api/v3/contatos?nome=aap" -H "Authorization: Bearer ACCESS TOKEN"
curl -X GET "https://api.tiny.com.br/public-api/v3/pedidos?nomeCliente=aguas" -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJnWXk2cDhkQkU0dDBkZkFvU0J4WkJvbDBkYmpTcEF5Z3FpQm1vY3pMcXJVIn0.eyJleHAiOjE3MzgzMjk4MzAsImlhdCI6MTczODMxNTQzMCwiYXV0aF90aW1lIjoxNzM4MzE1NDAxLCJqdGkiOiJkZGIyMjJjNy1jNjBkLTRlYWUtYTJhOC0wYjdhZjk3NzM0MjQiLCJpc3MiOiJodHRwczovL2FjY291bnRzLnRpbnkuY29tLmJyL3JlYWxtcy90aW55IiwiYXVkIjoidGlueS1hcGkiLCJzdWIiOiJjOGYxYjlhOC0xMzViLTRmMDEtYmQ4Ni1lZjE3YTE5NGZkZjUiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJ0aW55LWFwaS1kMGIzZTQ3ZDc2Njc2NWVlZGM4ZGFiZThmOTBlYmEzNzQ3NDI1N2MzLTE3MzcwMzc5MzEiLCJzZXNzaW9uX3N0YXRlIjoiOTMzYjIwMDktZjA4YS00NDgyLWJiYWEtMzI1YjU5ZmY1YzJjIiwic2NvcGUiOiJvcGVuaWQgZW1haWwgb2ZmbGluZV9hY2Nlc3MiLCJzaWQiOiI5MzNiMjAwOS1mMDhhLTQ0ODItYmJhYS0zMjViNTlmZjVjMmMiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsInJvbGVzIjp7InRpbnktYXBpIjpbIm1hcmNhcy1lc2NyaXRhIiwicGVkaWRvcy1sZWl0dXJhIiwiZXN0b3F1ZS1sZWl0dXJhIiwibm90YXMtZmlzY2Fpcy1lc2NyaXRhIiwiY29udGF0b3MtZXhjbHVzYW8iLCJpbmZvLWNvbnRhLWxlaXR1cmEiLCJwcm9kdXRvcy1leGNsdXNhbyIsImxpc3RhLXByZWNvcy1sZWl0dXJhIiwiY2F0ZWdvcmlhcy1sZWl0dXJhIiwiZXhwZWRpY2FvLWV4Y2x1c2FvIiwiZm9ybWEtZW52aW8tbGVpdHVyYSIsInNlcGFyYWNhby1sZWl0dXJhIiwicHJvZHV0b3MtbGVpdHVyYSIsImV4cGVkaWNhby1sZWl0dXJhIiwib3JkZW0tY29tcHJhLWxlaXR1cmEiLCJjb250YXRvcy1lc2NyaXRhIiwibm90YXMtZmlzY2Fpcy1sZWl0dXJhIiwicGVkaWRvcy1lc2NyaXRhIiwiY29udGFzLXBhZ2FyLWVzY3JpdGEiLCJvcmRlbS1zZXJ2aWNvLWxlaXR1cmEiLCJjb250YXMtcGFnYXItbGVpdHVyYSIsImV4cGVkaWNhby1lc2NyaXRhIiwiY29udGFzLXJlY2ViZXItbGVpdHVyYSIsInNlcGFyYWNhby1lc2NyaXRhIiwibWFyY2FzLWxlaXR1cmEiLCJmb3JtYS1wYWdhbWVudG8tbGVpdHVyYSIsImludGVybWVkaWFkb3Jlcy1sZWl0dXJhIiwiZ2F0aWxob3MiLCJlc3RvcXVlLWVzY3JpdGEiLCJwcm9kdXRvcy1lc2NyaXRhIiwib3JkZW0tc2Vydmljby1lc2NyaXRhIiwiY29udGFzLXJlY2ViZXItZXNjcml0YSIsImNvbnRhdG9zLWxlaXR1cmEiLCJwZWRpZG9zLWV4Y2x1c2FvIiwib3JkZW0tY29tcHJhLWVzY3JpdGEiLCJub3Rhcy1maXNjYWlzLWV4Y2x1c2FvIl19LCJlbWFpbCI6ImNvbWVyY2lhbEBzZWZmZXIuY29tLmJyIn0.Ky-MQ4juYFwTu5jfD6fLP1X0aykfp24Pcq2jORNVXaPr5SD0ASmG8AhQ6iDAxfNQ3qmYi8d3I2y975GUqidQtkLSjKsDYbwgH2p2lS4xKPELD-r62AHQ9eRLe3WvjUs-qeCJ-FEMptYy-w1GZ7pI6sAzZXpSLkNUvzcBohdhOUEesEN1xlUMRAmADccDFqHKNxvBUGqJJdBprAdqAhIDY8V-hIQNZMeVI8lzZo-0LRlInhMVx1E_TxovwnRKoc4fm8SY5vJG812xDNsVZMo6AGF3HkkTr0pJhBLTvP-YQLe77y7zq9W8kZqHKEEYa3w0Gy9J2_UkEBzzmlKdgsFfIg"

curl -X GET "https://api.tiny.com.br/public-api/v3/notas?cpfCnpj=20.810.361/0001-00" -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJnWXk2cDhkQkU0dDBkZkFvU0J4WkJvbDBkYmpTcEF5Z3FpQm1vY3pMcXJVIn0.eyJleHAiOjE3Mzg0MjMyNDAsImlhdCI6MTczODQwODg0MCwiYXV0aF90aW1lIjoxNzM4NDA4ODE2LCJqdGkiOiJkNmY2N2M5MC0wMzVhLTQ2ZTktODliYy05MDBjMWVhOTBhN2QiLCJpc3MiOiJodHRwczovL2FjY291bnRzLnRpbnkuY29tLmJyL3JlYWxtcy90aW55IiwiYXVkIjoidGlueS1hcGkiLCJzdWIiOiJjOGYxYjlhOC0xMzViLTRmMDEtYmQ4Ni1lZjE3YTE5NGZkZjUiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJ0aW55LWFwaS1kMGIzZTQ3ZDc2Njc2NWVlZGM4ZGFiZThmOTBlYmEzNzQ3NDI1N2MzLTE3MzcwMzc5MzEiLCJzZXNzaW9uX3N0YXRlIjoiNzE0NjY3ZDktYTk1ZC00NzI5LWFiYWMtYTAwNTQ5NjRkOWZiIiwic2NvcGUiOiJvcGVuaWQgZW1haWwgb2ZmbGluZV9hY2Nlc3MiLCJzaWQiOiI3MTQ2NjdkOS1hOTVkLTQ3MjktYWJhYy1hMDA1NDk2NGQ5ZmIiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsInJvbGVzIjp7InRpbnktYXBpIjpbIm1hcmNhcy1lc2NyaXRhIiwicGVkaWRvcy1sZWl0dXJhIiwiZXN0b3F1ZS1sZWl0dXJhIiwibm90YXMtZmlzY2Fpcy1lc2NyaXRhIiwiY29udGF0b3MtZXhjbHVzYW8iLCJpbmZvLWNvbnRhLWxlaXR1cmEiLCJwcm9kdXRvcy1leGNsdXNhbyIsImxpc3RhLXByZWNvcy1sZWl0dXJhIiwiY2F0ZWdvcmlhcy1sZWl0dXJhIiwiZXhwZWRpY2FvLWV4Y2x1c2FvIiwiZm9ybWEtZW52aW8tbGVpdHVyYSIsInNlcGFyYWNhby1sZWl0dXJhIiwicHJvZHV0b3MtbGVpdHVyYSIsImV4cGVkaWNhby1sZWl0dXJhIiwib3JkZW0tY29tcHJhLWxlaXR1cmEiLCJjb250YXRvcy1lc2NyaXRhIiwibm90YXMtZmlzY2Fpcy1sZWl0dXJhIiwicGVkaWRvcy1lc2NyaXRhIiwiY29udGFzLXBhZ2FyLWVzY3JpdGEiLCJvcmRlbS1zZXJ2aWNvLWxlaXR1cmEiLCJjb250YXMtcGFnYXItbGVpdHVyYSIsImV4cGVkaWNhby1lc2NyaXRhIiwiY29udGFzLXJlY2ViZXItbGVpdHVyYSIsInNlcGFyYWNhby1lc2NyaXRhIiwibWFyY2FzLWxlaXR1cmEiLCJmb3JtYS1wYWdhbWVudG8tbGVpdHVyYSIsImludGVybWVkaWFkb3Jlcy1sZWl0dXJhIiwiZ2F0aWxob3MiLCJlc3RvcXVlLWVzY3JpdGEiLCJwcm9kdXRvcy1lc2NyaXRhIiwib3JkZW0tc2Vydmljby1lc2NyaXRhIiwiY29udGFzLXJlY2ViZXItZXNjcml0YSIsImNvbnRhdG9zLWxlaXR1cmEiLCJwZWRpZG9zLWV4Y2x1c2FvIiwib3JkZW0tY29tcHJhLWVzY3JpdGEiLCJub3Rhcy1maXNjYWlzLWV4Y2x1c2FvIl19LCJlbWFpbCI6ImNvbWVyY2lhbEBzZWZmZXIuY29tLmJyIn0.YZLaMbkBsQKSKny0I6l8tgTXZV2GgxkbgfObhTiPw4MaqJ2Wc8NIvOHgZIQb_RBgnECNIqiAP5dtNtZGbE6Gm40Yyz_-aSleACR4XCu0wfq_lHRC0vNih8j-bRRLMvpC-Zmcoc93h7SMTs97H1voY19TAMMXv5QTknWEGJDm4i6KRfb1RKQ3P7d-hHFmD_sdZMDYh4cGgdjOvRsX5N39hNwfFX2kznun0Z3KWnFzPj1HEx9YbubM-OxyOaou_gd9mq4T-Ppx_zlDAYUJUxw6tfa3DvJK_5VIvBpPnkJ-JkjwWA_X6jqfsWNGc0zT5OBRKfjkFoBW7szJpWyQ-jiejg" 

# Markdown Cheatsheet

## Headers
# H1
## H2
### H3

## Emphasis
*Italic* or _Italic_
**Bold** or __Bold__
~~Strikethrough~~

## Lists
### Unordered
- Item 1
- Item 2
  - Subitem 1
  - Subitem 2

### Ordered
1. Item 1
2. Item 2

## Links
[Link Text](URL)

## Images
![Alt Text](Image URL)

## Code
Inline `code`
\`\`\`
Block of code
\`\`\`

## Blockquotes
> Quote

## Tables
| Header 1 | Header 2 |
|----------|----------|
| Cell 1   | Cell 2   |
| Cell 3   | Cell 4   |

## Horizontal Rule
---