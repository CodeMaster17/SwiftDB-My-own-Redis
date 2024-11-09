 
## SwiftDB

Redis uses RESP protocol.

## RESP

- Redis serialization protocol (RESP) is the wire protocol that clients implement.

- While the protocol was designed specifically for Redis, you can use it for other client-server software projects.

- RESP is a compromise among the following considerations:

    - Simple to implement.
    - Fast to parse.
    - Human readable.

- RESP can serialize different data types including integers, strings, and arrays. It also features an error-specific type. A client sends a request to the Redis server as an array of strings. The array's contents are the command and its arguments that the server should execute.

- RESP is binary-safe and uses prefixed length to transfer bulk data so it does not require processing bulk data transferred from one process to another. 


## How Redis works

- The key will be a string, and the value can be a serialized object of any type, such as an array, integer, or boolean.

```
SET admin harshit
```
- Redis receives these commands through a Serialization Protocol called RESP (Redis serialization protocol).

- So if we look at how SET admin ahmed is sent as a serialized message to Redis, it will look like this:

```
*3\r\n$3\r\nset\r\n$5\r\nadmin\r\n$5\r\harshit
```

- And to simplify it even more:

```
*3
$3
set
$5
admin
$5
harshit
```

- ‘\r\n’ is called CRLF and it indicates the end of a line.