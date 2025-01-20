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
https://accounts.tiny.com.br/realms/tiny/protocol/openid-connect/auth?client_id=tiny-api-d0b3e47d766765eedc8dabe8f90eba37474257c3-1737037931&redirect_uri=https://6c86-2804-10f8-4356-8600-64eb-2e6e-a771-dbed.ngrok-free.app&scope=openid&response_type=code

https://6c86-2804-10f8-4356-8600-64eb-2e6e-a771-dbed.ngrok-free.app/?
session_state=95c5eb2f-e8a8-4f48-8bf2-fa2882f6c607
&
code=34adade1-6ac4-4a5a-a394-2c47177a9311.95c5eb2f-e8a8-4f48-8bf2-fa2882f6c607.3dcda8a1-a6ef-4964-adcc-d0a5e1b8eebb
## Token Management
- The main token expires after 4 hours, and the refresh token expires after 24 hours, necessitating proper management to maintain connection, as Olist ERP Mediator cannot perform initial authentication.
- Tokens are stored in a PostgreSQL server accessible by all services, considering the application consists of different services possibly running in separate containers.
- To avoid high costs from continuous PostgreSQL server activity, a replication strategy writes the token to both the database and a file, ensuring token availability even if the database is down.
- Tokens are encrypted for security, with the encryption key stored in environmental variables.

## Next Steps
- In file.go, need to define a way to get the file_path, probably using a environmental variable
- Then, I can run a test by manually saving the token to the file and trying to read it to see the response
- Later, need to add the code to write the token to the file 
- Then, need to try reading it again

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