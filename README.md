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

CODE: 600e26d3-2300-44cf-90e2-7d061238fa2e.7cf40efb-398f-45f3-b94f-1262d6216d45.3dcda8a1-a6ef-4964-adcc-d0a5e1b8eebb

curl -X PUT "http://localhost:8081/auth?key=600e26d3-2300-44cf-90e2-7d061238fa2e.7cf40efb-398f-45f3-b94f-1262d6216d45.3dcda8a1-a6ef-4964-adcc-d0a5e1b8eebb" 
## Token Management
- The main token expires after 4 hours, and the refresh token expires after 24 hours, requiring proper management to maintain connection, as Olist ERP Mediator cannot perform initial authentication.
- Tokens are stored in a PostgreSQL server accessible by all services, considering the application consists of different services possibly running in separate containers.
- To avoid high costs from continuous PostgreSQL server activity, a replication strategy writes the token to both the database and a JSON file, ensuring token availability even if the database is down.
- Tokens are encrypted for security, with the encryption key stored in environmental variables.

## Next Steps
- Start development of main application
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