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
- Validate the circuit breaker, add the debouncer!! 
- In Get Orders add features to get orders based on other searches, such as order number, date filter and etc
- Valor field in some places (order model for example) is string, but should not be!
- To have a key value store with OLIST ERP ENDPOINTS
- Review code according to Go standards
- Add tests
- Add benchmarks/profile
- Create documentation
- Deploy to cloud?
- Integrate with whats app 
- Review error handling according to the 5 patterns
- Propage the context from url, in a way that if the customer cancels the request, then it sends the cancellation
- Add a debounce to every endpoint request, actually is a good idea to chain both! circuit breaker + debounce!
- Add in app the condition to renew authentication every couple hours 

- (ONLY AFTER DEPLOYING THIS FIRST VERSION)
- Create a new service for RTE 
- Create a new service for BTU
- Consider moving from Gorilla/mux to Chi or Gin, and understand why it can't use the basic http option


## Future Developments
- Perform the refresh token through AWS Functions with a trigger every 8 hours. Then, even when the server is down, the tokens will still refresh. Also, will allow to turn off and turn on the servers based on users activity. For example, weekends, holidays, and other dates where there may be no requests, there would be no reason to keep the host on. Would also have to create a trigger to start the server whenever needed (based on incoming request). Also, it would request one intermediator to get the requests and turn on/off the remaining services on demand. Not sure how to achieve this, btw.

# Testing
1. Get one time authentication code:
One Time Authentication: https://accounts.tiny.com.br/realms/tiny/protocol/openid-connect/auth?client_id=tiny-api-d0b3e47d766765eedc8dabe8f90eba37474257c3-1737037931&redirect_uri=https://6c86-2804-10f8-4356-8600-64eb-2e6e-a771-dbed.ngrok-free.app&scope=openid&response_type=code

2. Code returned after authentication: 08f29cb7-cb8d-4278-95e7-8fefe08ff790.2fbf3a57-634f-43d6-9539-6330b39ed4cc.3dcda8a1-a6ef-4964-adcc-d0a5e1b8eebb

3. Use code to get the access token: curl --location 'https://accounts.tiny.com.br/realms/tiny/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'grant_type=authorization_code' \
--data-urlencode 'client_id=tiny-api-d0b3e47d766765eedc8dabe8f90eba37474257c3-1737037931' \
--data-urlencode 'client_secret=bpSCcOcdxN8tftz4X8vRRceGd1FTguyw' \
--data-urlencode 'redirect_uri=https://6c86-2804-10f8-4356-8600-64eb-2e6e-a771-dbed.ngrok-free.app' \
--data-urlencode 'code=08f29cb7-cb8d-4278-95e7-8fefe08ff790.2fbf3a57-634f-43d6-9539-6330b39ed4cc.3dcda8a1-a6ef-4964-adcc-d0a5e1b8eebb'

4. Get the access token from the response: 
{"access_token":"ACCESS_TOKEN",
"expires_in":14400,
"refresh_expires_in":86400,
"refresh_token":"REFRESH_TOKEN",
"not-before-policy":0,
"session_state":"7cf40efb-398f-45f3-b94f-1262d6216d45",
"scope":"openid email offline_access"}

5. Store the token to the token repositories:
curl -X PUT "http://localhost:8081/auth?key=eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJnWXk2cDhkQkU0dDBkZkFvU0J4WkJvbDBkYmpTcEF5Z3FpQm1vY3pMcXJVIn0.eyJleHAiOjE3Mzk1NjM5MzYsImlhdCI6MTczOTU0OTUzNiwiYXV0aF90aW1lIjoxNzM5NTQ2NTUzLCJqdGkiOiI2Nzg2YTBiMC0yNzUzLTQ2ZDgtODYzYS1jYTM5YWFmZWUxYzEiLCJpc3MiOiJodHRwczovL2FjY291bnRzLnRpbnkuY29tLmJyL3JlYWxtcy90aW55IiwiYXVkIjoidGlueS1hcGkiLCJzdWIiOiJjOGYxYjlhOC0xMzViLTRmMDEtYmQ4Ni1lZjE3YTE5NGZkZjUiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJ0aW55LWFwaS1kMGIzZTQ3ZDc2Njc2NWVlZGM4ZGFiZThmOTBlYmEzNzQ3NDI1N2MzLTE3MzcwMzc5MzEiLCJzZXNzaW9uX3N0YXRlIjoiMmZiZjNhNTctNjM0Zi00M2Q2LTk1MzktNjMzMGIzOWVkNGNjIiwic2NvcGUiOiJvcGVuaWQgZW1haWwgb2ZmbGluZV9hY2Nlc3MiLCJzaWQiOiIyZmJmM2E1Ny02MzRmLTQzZDYtOTUzOS02MzMwYjM5ZWQ0Y2MiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsInJvbGVzIjp7InRpbnktYXBpIjpbIm1hcmNhcy1lc2NyaXRhIiwicGVkaWRvcy1sZWl0dXJhIiwiZXN0b3F1ZS1sZWl0dXJhIiwibm90YXMtZmlzY2Fpcy1lc2NyaXRhIiwiY29udGF0b3MtZXhjbHVzYW8iLCJpbmZvLWNvbnRhLWxlaXR1cmEiLCJwcm9kdXRvcy1leGNsdXNhbyIsImxpc3RhLXByZWNvcy1sZWl0dXJhIiwiY2F0ZWdvcmlhcy1sZWl0dXJhIiwiZXhwZWRpY2FvLWV4Y2x1c2FvIiwiZm9ybWEtZW52aW8tbGVpdHVyYSIsInNlcGFyYWNhby1sZWl0dXJhIiwicHJvZHV0b3MtbGVpdHVyYSIsImV4cGVkaWNhby1sZWl0dXJhIiwib3JkZW0tY29tcHJhLWxlaXR1cmEiLCJjb250YXRvcy1lc2NyaXRhIiwibm90YXMtZmlzY2Fpcy1sZWl0dXJhIiwicGVkaWRvcy1lc2NyaXRhIiwiY29udGFzLXBhZ2FyLWVzY3JpdGEiLCJvcmRlbS1zZXJ2aWNvLWxlaXR1cmEiLCJjb250YXMtcGFnYXItbGVpdHVyYSIsImV4cGVkaWNhby1lc2NyaXRhIiwiY29udGFzLXJlY2ViZXItbGVpdHVyYSIsInNlcGFyYWNhby1lc2NyaXRhIiwibWFyY2FzLWxlaXR1cmEiLCJmb3JtYS1wYWdhbWVudG8tbGVpdHVyYSIsImludGVybWVkaWFkb3Jlcy1sZWl0dXJhIiwiZ2F0aWxob3MiLCJlc3RvcXVlLWVzY3JpdGEiLCJwcm9kdXRvcy1lc2NyaXRhIiwib3JkZW0tc2Vydmljby1lc2NyaXRhIiwiY29udGFzLXJlY2ViZXItZXNjcml0YSIsImNvbnRhdG9zLWxlaXR1cmEiLCJwZWRpZG9zLWV4Y2x1c2FvIiwib3JkZW0tY29tcHJhLWVzY3JpdGEiLCJub3Rhcy1maXNjYWlzLWV4Y2x1c2FvIl19LCJlbWFpbCI6ImNvbWVyY2lhbEBzZWZmZXIuY29tLmJyIn0.SOkb3N6Cgm5c6GD0QCSpnAaKBchPyrZWRRWdSHyEQzM2-UE94c4TlblwiT4RR5QsQGDz3vMOVZ6dtkEDXeUQToCMQO5bFyeChsvbCZ4QPL1Hm2utLgnzYXCJNCH2QZsq6BMyDeg90XJuhHiEi8XgWnQ_QEfiwSEyvGI9B_1jCRinMjbnMbsM4RGLEhwp_5tSoC8pH5OKXTNpbvlsI0gAp-Iph70-bvGjuFirC-w0XXQdJ1KcoqAuVAO6-L7JfZTKJ0Dt__gL5ntKxRU7X1VELPPVfkvNrmoXPw9uno926v6VRFNEZHDjkU5v6DGxPGYBVGxqYGIYPy8lP1OPTk2ssw&refresh_token=eyJhbGciOiJIUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICI5MDNkY2IxNi0xNTgwLTQ0NzYtYWY0Mi1lYTFkN2YyNzY4MjcifQ.eyJleHAiOjE3Mzk2MzU5MzYsImlhdCI6MTczOTU0OTUzNiwianRpIjoiN2FlMWJjMDItZmU4ZC00ZGU4LTlmOTAtMjVmZjkxYmMyMmMxIiwiaXNzIjoiaHR0cHM6Ly9hY2NvdW50cy50aW55LmNvbS5ici9yZWFsbXMvdGlueSIsImF1ZCI6Imh0dHBzOi8vYWNjb3VudHMudGlueS5jb20uYnIvcmVhbG1zL3RpbnkiLCJzdWIiOiJjOGYxYjlhOC0xMzViLTRmMDEtYmQ4Ni1lZjE3YTE5NGZkZjUiLCJ0eXAiOiJPZmZsaW5lIiwiYXpwIjoidGlueS1hcGktZDBiM2U0N2Q3NjY3NjVlZWRjOGRhYmU4ZjkwZWJhMzc0NzQyNTdjMy0xNzM3MDM3OTMxIiwic2Vzc2lvbl9zdGF0ZSI6IjJmYmYzYTU3LTYzNGYtNDNkNi05NTM5LTYzMzBiMzllZDRjYyIsInNjb3BlIjoib3BlbmlkIGVtYWlsIG9mZmxpbmVfYWNjZXNzIiwic2lkIjoiMmZiZjNhNTctNjM0Zi00M2Q2LTk1MzktNjMzMGIzOWVkNGNjIn0.7JJ9oDOuZ7IFanb7S6Z9NWKwc8XbdEO2RQ9bGjqow8g" 

6. Test call to the OLIST ERP API:
curl -X GET "https://api.tiny.com.br/public-api/v3/contatos?nome=aap" -H "Authorization: Bearer ACCESS TOKEN"
curl -X GET "https://api.tiny.com.br/public-api/v3/produtos?nome=sertech" -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJnWXk2cDhkQkU0dDBkZkFvU0J4WkJvbDBkYmpTcEF5Z3FpQm1vY3pMcXJVIn0.eyJleHAiOjE3Mzg3NjQ4OTksImlhdCI6MTczODc1MDQ5OSwiYXV0aF90aW1lIjoxNzM4NzUwNDYzLCJqdGkiOiJhMTBhOTY0OS0xNzVmLTQxOGEtOTEyMi03ZWI4MTM5NjE0NTEiLCJpc3MiOiJodHRwczovL2FjY291bnRzLnRpbnkuY29tLmJyL3JlYWxtcy90aW55IiwiYXVkIjoidGlueS1hcGkiLCJzdWIiOiJjOGYxYjlhOC0xMzViLTRmMDEtYmQ4Ni1lZjE3YTE5NGZkZjUiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJ0aW55LWFwaS1kMGIzZTQ3ZDc2Njc2NWVlZGM4ZGFiZThmOTBlYmEzNzQ3NDI1N2MzLTE3MzcwMzc5MzEiLCJzZXNzaW9uX3N0YXRlIjoiMDg0ZmFiMTktMDZjMS00ZDZiLTllZDQtMWM1OTI5YzIxMGM0Iiwic2NvcGUiOiJvcGVuaWQgZW1haWwgb2ZmbGluZV9hY2Nlc3MiLCJzaWQiOiIwODRmYWIxOS0wNmMxLTRkNmItOWVkNC0xYzU5MjljMjEwYzQiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsInJvbGVzIjp7InRpbnktYXBpIjpbIm1hcmNhcy1lc2NyaXRhIiwicGVkaWRvcy1sZWl0dXJhIiwiZXN0b3F1ZS1sZWl0dXJhIiwibm90YXMtZmlzY2Fpcy1lc2NyaXRhIiwiY29udGF0b3MtZXhjbHVzYW8iLCJpbmZvLWNvbnRhLWxlaXR1cmEiLCJwcm9kdXRvcy1leGNsdXNhbyIsImxpc3RhLXByZWNvcy1sZWl0dXJhIiwiY2F0ZWdvcmlhcy1sZWl0dXJhIiwiZXhwZWRpY2FvLWV4Y2x1c2FvIiwiZm9ybWEtZW52aW8tbGVpdHVyYSIsInNlcGFyYWNhby1sZWl0dXJhIiwicHJvZHV0b3MtbGVpdHVyYSIsImV4cGVkaWNhby1sZWl0dXJhIiwib3JkZW0tY29tcHJhLWxlaXR1cmEiLCJjb250YXRvcy1lc2NyaXRhIiwibm90YXMtZmlzY2Fpcy1sZWl0dXJhIiwicGVkaWRvcy1lc2NyaXRhIiwiY29udGFzLXBhZ2FyLWVzY3JpdGEiLCJvcmRlbS1zZXJ2aWNvLWxlaXR1cmEiLCJjb250YXMtcGFnYXItbGVpdHVyYSIsImV4cGVkaWNhby1lc2NyaXRhIiwiY29udGFzLXJlY2ViZXItbGVpdHVyYSIsInNlcGFyYWNhby1lc2NyaXRhIiwibWFyY2FzLWxlaXR1cmEiLCJmb3JtYS1wYWdhbWVudG8tbGVpdHVyYSIsImludGVybWVkaWFkb3Jlcy1sZWl0dXJhIiwiZ2F0aWxob3MiLCJlc3RvcXVlLWVzY3JpdGEiLCJwcm9kdXRvcy1lc2NyaXRhIiwib3JkZW0tc2Vydmljby1lc2NyaXRhIiwiY29udGFzLXJlY2ViZXItZXNjcml0YSIsImNvbnRhdG9zLWxlaXR1cmEiLCJwZWRpZG9zLWV4Y2x1c2FvIiwib3JkZW0tY29tcHJhLWVzY3JpdGEiLCJub3Rhcy1maXNjYWlzLWV4Y2x1c2FvIl19LCJlbWFpbCI6ImNvbWVyY2lhbEBzZWZmZXIuY29tLmJyIn0.DT8xeMF9IeZMWBXg7a4H5BY3tYu0e-N1c1Sd04HBYIRlQtuPs12KgI-YZbwAzJP74aj3LjYwTICO4Tg17BFpiNasfwhxDl-t00NlrWH7VanXZF0g58AMdgMbnxgMvBpLjWmW2MZbSuURnEIIqNB2V1rxXFWF_vFqL4PUxzavZ3_zDG-lQUU69_qCHs8GKWjN3b6p2Zn-qUM2sP2X0nR6DYPgVmmOvRyshLY-Z15aJwVdhMWXIxlvT_49ze6IT_iq9x7T2zuhCnamsgnbb0fLN5oaAy94rNE1VvZiTqkSoEVLHM_Um88SS9NhB1bxruBv4SrpLDHLZLdcrWftMX8QSQ"

curl -X GET "https://api.tiny.com.br/public-api/v3/notas?cpfCnpj=20.810.361/0001-00" -H "Authorization: Bearer eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJnWXk2cDhkQkU0dDBkZkFvU0J4WkJvbDBkYmpTcEF5Z3FpQm1vY3pMcXJVIn0.eyJleHAiOjE3Mzg3NjQ4OTksImlhdCI6MTczODc1MDQ5OSwiYXV0aF90aW1lIjoxNzM4NzUwNDYzLCJqdGkiOiJhMTBhOTY0OS0xNzVmLTQxOGEtOTEyMi03ZWI4MTM5NjE0NTEiLCJpc3MiOiJodHRwczovL2FjY291bnRzLnRpbnkuY29tLmJyL3JlYWxtcy90aW55IiwiYXVkIjoidGlueS1hcGkiLCJzdWIiOiJjOGYxYjlhOC0xMzViLTRmMDEtYmQ4Ni1lZjE3YTE5NGZkZjUiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJ0aW55LWFwaS1kMGIzZTQ3ZDc2Njc2NWVlZGM4ZGFiZThmOTBlYmEzNzQ3NDI1N2MzLTE3MzcwMzc5MzEiLCJzZXNzaW9uX3N0YXRlIjoiMDg0ZmFiMTktMDZjMS00ZDZiLTllZDQtMWM1OTI5YzIxMGM0Iiwic2NvcGUiOiJvcGVuaWQgZW1haWwgb2ZmbGluZV9hY2Nlc3MiLCJzaWQiOiIwODRmYWIxOS0wNmMxLTRkNmItOWVkNC0xYzU5MjljMjEwYzQiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsInJvbGVzIjp7InRpbnktYXBpIjpbIm1hcmNhcy1lc2NyaXRhIiwicGVkaWRvcy1sZWl0dXJhIiwiZXN0b3F1ZS1sZWl0dXJhIiwibm90YXMtZmlzY2Fpcy1lc2NyaXRhIiwiY29udGF0b3MtZXhjbHVzYW8iLCJpbmZvLWNvbnRhLWxlaXR1cmEiLCJwcm9kdXRvcy1leGNsdXNhbyIsImxpc3RhLXByZWNvcy1sZWl0dXJhIiwiY2F0ZWdvcmlhcy1sZWl0dXJhIiwiZXhwZWRpY2FvLWV4Y2x1c2FvIiwiZm9ybWEtZW52aW8tbGVpdHVyYSIsInNlcGFyYWNhby1sZWl0dXJhIiwicHJvZHV0b3MtbGVpdHVyYSIsImV4cGVkaWNhby1sZWl0dXJhIiwib3JkZW0tY29tcHJhLWxlaXR1cmEiLCJjb250YXRvcy1lc2NyaXRhIiwibm90YXMtZmlzY2Fpcy1sZWl0dXJhIiwicGVkaWRvcy1lc2NyaXRhIiwiY29udGFzLXBhZ2FyLWVzY3JpdGEiLCJvcmRlbS1zZXJ2aWNvLWxlaXR1cmEiLCJjb250YXMtcGFnYXItbGVpdHVyYSIsImV4cGVkaWNhby1lc2NyaXRhIiwiY29udGFzLXJlY2ViZXItbGVpdHVyYSIsInNlcGFyYWNhby1lc2NyaXRhIiwibWFyY2FzLWxlaXR1cmEiLCJmb3JtYS1wYWdhbWVudG8tbGVpdHVyYSIsImludGVybWVkaWFkb3Jlcy1sZWl0dXJhIiwiZ2F0aWxob3MiLCJlc3RvcXVlLWVzY3JpdGEiLCJwcm9kdXRvcy1lc2NyaXRhIiwib3JkZW0tc2Vydmljby1lc2NyaXRhIiwiY29udGFzLXJlY2ViZXItZXNjcml0YSIsImNvbnRhdG9zLWxlaXR1cmEiLCJwZWRpZG9zLWV4Y2x1c2FvIiwib3JkZW0tY29tcHJhLWVzY3JpdGEiLCJub3Rhcy1maXNjYWlzLWV4Y2x1c2FvIl19LCJlbWFpbCI6ImNvbWVyY2lhbEBzZWZmZXIuY29tLmJyIn0.DT8xeMF9IeZMWBXg7a4H5BY3tYu0e-N1c1Sd04HBYIRlQtuPs12KgI-YZbwAzJP74aj3LjYwTICO4Tg17BFpiNasfwhxDl-t00NlrWH7VanXZF0g58AMdgMbnxgMvBpLjWmW2MZbSuURnEIIqNB2V1rxXFWF_vFqL4PUxzavZ3_zDG-lQUU69_qCHs8GKWjN3b6p2Zn-qUM2sP2X0nR6DYPgVmmOvRyshLY-Z15aJwVdhMWXIxlvT_49ze6IT_iq9x7T2zuhCnamsgnbb0fLN5oaAy94rNE1VvZiTqkSoEVLHM_Um88SS9NhB1bxruBv4SrpLDHLZLdcrWftMX8QSQ" 

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