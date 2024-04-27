This project illustrates the distributed tracing with jaeger and opentelmetry. Router receives the request it does some checks and send to the grpc server. Grpc Server has implementation of Open AI api. It also has some test cases under test folder.

We have following routes:

i.  /answer  : which gives the answer of the question.
ii. /search : which searches the answer provided in the given context.

We have jaeger to see the traces.



