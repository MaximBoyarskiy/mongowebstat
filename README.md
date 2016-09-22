It is obsolete, unfortunatelly
==============================

mongowebstat
============

- Monitor your mongodb in more convenient way.
- Monitor as many mongodb hosts as you want!
- Simple json based config.
- Table based now and graphs in the nearest future!
- Mongo 2.4.X compatibility. Older version are in plans.

Build from sources
------------------

_Install golang:_ http://golang.org/doc/install

_Install dependency:_ <code>go get "labix.org/v2/mgo"</code>

_Compile your mongowebstat:_ <code>go build main.go</code>

_Or just run it:_ <code>go run main.go</code>

Or use the compiled binary:
-----------------------

- copy your platform binary from bin folder
- put it in near _static, templates_ folders
- run it <code>./main-linux-386</code> on 8080 port or run it <code>./main-linux-386 -http=:8090</code> on 8090 or any other port.

**Go to http://127.0.0.1:8080 and monitor your mongodb!**

Check out mongowebstat.json: 
----------------------------

- put connection string to **"host"** field
- set **"http"** to true for http connection or false for direct connection
- put some short node name to **"name"** field
- use as many host records as you want!


It is made available under the [Simplified BSD License](http://en.wikipedia.org/wiki/BSD_licenses#2-clause_license_.28.22Simplified_BSD_License.22_or_.22FreeBSD_License.22.29)

**And [remember](https://github.com/MaximBoyarskiy/mongowebstat/blob/master/src/static/like-a-boss.jpg)!**

See roadmap.txt
