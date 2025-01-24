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

One Time Authentication: https://accounts.tiny.com.br/realms/tiny/protocol/openid-connect/auth?client_id=tiny-api-d0b3e47d766765eedc8dabe8f90eba37474257c3-1737037931&redirect_uri=https://6c86-2804-10f8-4356-8600-64eb-2e6e-a771-dbed.ngrok-free.app&scope=openid&response_type=code

CODE: 3fc78f19-77bb-4751-9aac-a3fe67743075.7cf40efb-398f-45f3-b94f-1262d6216d45.3dcda8a1-a6ef-4964-adcc-d0a5e1b8eebb

curl -X PUT "http://localhost:8081/auth?key=600e26d3-2300-44cf-90e2-7d061238fa2e.7cf40efb-398f-45f3-b94f-1262d6216d45.3dcda8a1-a6ef-4964-adcc-d0a5e1b8eebb" 

curl --location 'https://accounts.tiny.com.br/realms/tiny/protocol/openid-connect/token' \
--header 'Content-Type: application/x-www-form-urlencoded' \
--data-urlencode 'grant_type=authorization_code' \
--data-urlencode 'client_id=tiny-api-d0b3e47d766765eedc8dabe8f90eba37474257c3-1737037931' \
--data-urlencode 'client_secret=bpSCcOcdxN8tftz4X8vRRceGd1FTguyw' \
--data-urlencode 'redirect_uri=https://6c86-2804-10f8-4356-8600-64eb-2e6e-a771-dbed.ngrok-free.app' \
--data-urlencode 'code=3fc78f19-77bb-4751-9aac-a3fe67743075.7cf40efb-398f-45f3-b94f-1262d6216d45.3dcda8a1-a6ef-4964-adcc-d0a5e1b8eebb'

{"access_token":"eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJnWXk2cDhkQkU0dDBkZkFvU0J4WkJvbDBkYmpTcEF5Z3FpQm1vY3pMcXJVIn0.eyJleHAiOjE3Mzc3NDI2MzksImlhdCI6MTczNzcyODIzOSwiYXV0aF90aW1lIjoxNzM3NzE4NTc2LCJqdGkiOiI5ZTRkNDhkNS0wYzA4LTRlMmQtOTU1ZC01ODUyZWNhN2JhNTEiLCJpc3MiOiJodHRwczovL2FjY291bnRzLnRpbnkuY29tLmJyL3JlYWxtcy90aW55IiwiYXVkIjoidGlueS1hcGkiLCJzdWIiOiJjOGYxYjlhOC0xMzViLTRmMDEtYmQ4Ni1lZjE3YTE5NGZkZjUiLCJ0eXAiOiJCZWFyZXIiLCJhenAiOiJ0aW55LWFwaS1kMGIzZTQ3ZDc2Njc2NWVlZGM4ZGFiZThmOTBlYmEzNzQ3NDI1N2MzLTE3MzcwMzc5MzEiLCJzZXNzaW9uX3N0YXRlIjoiN2NmNDBlZmItMzk4Zi00NWYzLWI5NGYtMTI2MmQ2MjE2ZDQ1Iiwic2NvcGUiOiJvcGVuaWQgZW1haWwgb2ZmbGluZV9hY2Nlc3MiLCJzaWQiOiI3Y2Y0MGVmYi0zOThmLTQ1ZjMtYjk0Zi0xMjYyZDYyMTZkNDUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsInJvbGVzIjp7InRpbnktYXBpIjpbIm1hcmNhcy1lc2NyaXRhIiwicGVkaWRvcy1sZWl0dXJhIiwiZXN0b3F1ZS1sZWl0dXJhIiwibm90YXMtZmlzY2Fpcy1lc2NyaXRhIiwiY29udGF0b3MtZXhjbHVzYW8iLCJpbmZvLWNvbnRhLWxlaXR1cmEiLCJwcm9kdXRvcy1leGNsdXNhbyIsImxpc3RhLXByZWNvcy1sZWl0dXJhIiwiY2F0ZWdvcmlhcy1sZWl0dXJhIiwiZXhwZWRpY2FvLWV4Y2x1c2FvIiwiZm9ybWEtZW52aW8tbGVpdHVyYSIsInNlcGFyYWNhby1sZWl0dXJhIiwicHJvZHV0b3MtbGVpdHVyYSIsImV4cGVkaWNhby1sZWl0dXJhIiwib3JkZW0tY29tcHJhLWxlaXR1cmEiLCJjb250YXRvcy1lc2NyaXRhIiwibm90YXMtZmlzY2Fpcy1sZWl0dXJhIiwicGVkaWRvcy1lc2NyaXRhIiwiY29udGFzLXBhZ2FyLWVzY3JpdGEiLCJvcmRlbS1zZXJ2aWNvLWxlaXR1cmEiLCJjb250YXMtcGFnYXItbGVpdHVyYSIsImV4cGVkaWNhby1lc2NyaXRhIiwiY29udGFzLXJlY2ViZXItbGVpdHVyYSIsInNlcGFyYWNhby1lc2NyaXRhIiwibWFyY2FzLWxlaXR1cmEiLCJmb3JtYS1wYWdhbWVudG8tbGVpdHVyYSIsImludGVybWVkaWFkb3Jlcy1sZWl0dXJhIiwiZ2F0aWxob3MiLCJlc3RvcXVlLWVzY3JpdGEiLCJwcm9kdXRvcy1lc2NyaXRhIiwib3JkZW0tc2Vydmljby1lc2NyaXRhIiwiY29udGFzLXJlY2ViZXItZXNjcml0YSIsImNvbnRhdG9zLWxlaXR1cmEiLCJwZWRpZG9zLWV4Y2x1c2FvIiwib3JkZW0tY29tcHJhLWVzY3JpdGEiLCJub3Rhcy1maXNjYWlzLWV4Y2x1c2FvIl19LCJlbWFpbCI6ImNvbWVyY2lhbEBzZWZmZXIuY29tLmJyIn0.l-ypl-8A7o1dxohPKK0EoBjrfLkBKAH0srFDToFUcDcOCyQt193Qtn0jfyyzVSYDUqj2yDOb30dZS0eDJlSlerfgNfmlqfi6Df4bhTGiJ-vTHhSv9-xhcUueYP8Cb9LblRyvHHDRfEYBFScHfocw8GxSPMS_JLRZvjnakCvYQucDi2cCivLFrUQ2mw4d6MQ7RfjQTc8F4TFBjohl0P4BJf8opo2G3lzKFiuwsH5nf5FG-fYep7zQvXomU_RgdwUyf33TP2XycHZhy11An7F2sUtrSceSQsTHIm-orLwDKFU8ZrntHMxMuPMcShBnW43xHr4vUYB_jEpqpiuxJTDTBg",
"expires_in":14400,
"refresh_expires_in":86400,"refresh_token":"eyJhbGciOiJIUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICI5MDNkY2IxNi0xNTgwLTQ0NzYtYWY0Mi1lYTFkN2YyNzY4MjcifQ.eyJleHAiOjE3Mzc4MTQ2MzksImlhdCI6MTczNzcyODIzOSwianRpIjoiNjlhZTQyNGUtZWJjMC00NDQ5LWJjNTAtM2FmOTdjOTI5MjUzIiwiaXNzIjoiaHR0cHM6Ly9hY2NvdW50cy50aW55LmNvbS5ici9yZWFsbXMvdGlueSIsImF1ZCI6Imh0dHBzOi8vYWNjb3VudHMudGlueS5jb20uYnIvcmVhbG1zL3RpbnkiLCJzdWIiOiJjOGYxYjlhOC0xMzViLTRmMDEtYmQ4Ni1lZjE3YTE5NGZkZjUiLCJ0eXAiOiJPZmZsaW5lIiwiYXpwIjoidGlueS1hcGktZDBiM2U0N2Q3NjY3NjVlZWRjOGRhYmU4ZjkwZWJhMzc0NzQyNTdjMy0xNzM3MDM3OTMxIiwic2Vzc2lvbl9zdGF0ZSI6IjdjZjQwZWZiLTM5OGYtNDVmMy1iOTRmLTEyNjJkNjIxNmQ0NSIsInNjb3BlIjoib3BlbmlkIGVtYWlsIG9mZmxpbmVfYWNjZXNzIiwic2lkIjoiN2NmNDBlZmItMzk4Zi00NWYzLWI5NGYtMTI2MmQ2MjE2ZDQ1In0.WWjSwsD4TP-qK27ZMNmxZLgPIm-sHHEVsQoPu1tGVJk","token_type":"Bearer","id_token":"eyJhbGciOiJSUzI1NiIsInR5cCIgOiAiSldUIiwia2lkIiA6ICJnWXk2cDhkQkU0dDBkZkFvU0J4WkJvbDBkYmpTcEF5Z3FpQm1vY3pMcXJVIn0.eyJleHAiOjE3Mzc3NDI2MzksImlhdCI6MTczNzcyODIzOSwiYXV0aF90aW1lIjoxNzM3NzE4NTc2LCJqdGkiOiI2OGM1MDJkYy05MWNjLTQ4MWYtOTBiNC01YThlY2I2NTgyNTYiLCJpc3MiOiJodHRwczovL2FjY291bnRzLnRpbnkuY29tLmJyL3JlYWxtcy90aW55IiwiYXVkIjoidGlueS1hcGktZDBiM2U0N2Q3NjY3NjVlZWRjOGRhYmU4ZjkwZWJhMzc0NzQyNTdjMy0xNzM3MDM3OTMxIiwic3ViIjoiYzhmMWI5YTgtMTM1Yi00ZjAxLWJkODYtZWYxN2ExOTRmZGY1IiwidHlwIjoiSUQiLCJhenAiOiJ0aW55LWFwaS1kMGIzZTQ3ZDc2Njc2NWVlZGM4ZGFiZThmOTBlYmEzNzQ3NDI1N2MzLTE3MzcwMzc5MzEiLCJzZXNzaW9uX3N0YXRlIjoiN2NmNDBlZmItMzk4Zi00NWYzLWI5NGYtMTI2MmQ2MjE2ZDQ1IiwiYXRfaGFzaCI6Ik5ySU5TakJZY3FvTGtfQmIweklKa1EiLCJzaWQiOiI3Y2Y0MGVmYi0zOThmLTQ1ZjMtYjk0Zi0xMjYyZDYyMTZkNDUiLCJlbWFpbF92ZXJpZmllZCI6ZmFsc2UsImVtYWlsIjoiY29tZXJjaWFsQHNlZmZlci5jb20uYnIifQ.kgkjyttcY9FDqDtlzhM8meBS9hRVxwycJ82ICgJC6U9aNhgXgVX-VJzYZwCDb-df-A8YSFC5zm4c4ks5_LH0Mz0E8LRahW4Gh6ZQIKWKvznFxa8Gv-FK7PE818_qQCcxJnldy6ad28KtaU_Caf_GJ2PoUzjPTpRJFFNxsWzvh9o_dRlHyaQ8kEq1Ppla2h2Y1uZsXy_Ixuw4Fv5KifVTD9YP6QAZ_60zItlQ4UYs3IEjUpQTJzRqvB1dOjU54ITPglKSk4j7jC7DpMeavWoUVIH_mBPwWkYpS3JOxeIiY0XhlwCw8z7PRM99PSHjpIyH1VwxLKfQv8j4TzNfKI83mA","not-before-policy":0,"session_state":"7cf40efb-398f-45f3-b94f-1262d6216d45","scope":"openid email offline_access"}%


## Token Management
- The main token expires after 4 hours, and the refresh token expires after 24 hours, requiring proper management to maintain connection, as Olist ERP Mediator cannot perform initial authentication.
- Tokens are stored in a PostgreSQL server accessible by all services, considering the application consists of different services possibly running in separate containers.
- To avoid high costs from continuous PostgreSQL server activity, a replication strategy writes the token to both the database and a JSON file, ensuring token availability even if the database is down.
- Tokens are encrypted for security, with the encryption key stored in environmental variables.

## Next Steps
- Fix the issue where it is not possible to save the access token due to expected size (base64 error)
- Continue development of main application
- To have a key value store with OLIST ERP ENDPOINTS

## Future Developments
- Perform the refresh token through AWS Functions with a trigger every 8 hours. Then, even when the server is down, the tokens will still refresh. Also, will allow to turn off and turn on the servers based on users activity. For example, weekends, holidays, and other dates where there may be no requests, there would be no reason to keep the host on. Would also have to create a trigger to start the server whenever needed (based on incoming request). Also, it would request one intermediator to get the requests and turn on/off the remaining services on demand. Not sure how to achieve this, btw.

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