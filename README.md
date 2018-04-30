# Simple TFTP Server
The goal is to build a simple TFTP server that only stores files in memory.  There’s no need to write out files.  For more information on the TFTP protocol check out the RFC:

http://tools.ietf.org/html/rfc1350

The server can bind a port other than 69 but still should be implemented using UDP.  Only the “octet” mode needs to be implemented.  Requests should be handled concurrently, but files being written to the server must not be visible until completed.

## Ground Rules
* The goal is to have functional code that can be the basis of a conversation around design trade-offs during the face-to-face interview process. 
* Expect to take at least a few hours, but this is not a timed assignment.  The most important goal is making a functional piece of code you would want to be reviewed by a colleague.
* Include all source code, test code, and a README describing how to build and run the solution.
* We think of languages as tools and good engineers can adapt to the right tool for the job.  As we primarily program in golang at Igneous (and it happens to be a great language for this type of project), please submit the code in golang.

Hint: OSX and (most) Linux distros include a tftp client that can be used for ad-hoc testing.

# Running it
* Clone git@github.com:AToner/go_tftp.git
* Make sure you're in the same directory as the source files
* Run the tests "go test"
* Run the real thing "go run tftp.go"
