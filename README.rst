Go stats handler
================

.. image:: https://secure.travis-ci.org/thoas/stats.png?branch=master
    :alt: Build Status
    :target: http://travis-ci.org/thoas/stats

stats is a ``valyala/fasthttp`` handler in golang reporting various metrics about
your web application.

This middleware has been developed and required for the need of picfit_,
an image resizing server written in Go.

Compatibility
-------------

This handler supports the following frameworks at the moment:

* `fasthttp <https://github.com/valyala/fasthttp>`_

We don't support your favorite Go framework? Send me a PR or
create a new `issue <https://github.com/thoas/stats/issues>`_ and
I will implement it :)

Installation
------------

1. Make sure you have a Go language compiler >= 1.3 (required) and git installed.
2. Make sure you have the following go system dependencies in your $PATH: bzr, svn, hg, git
3. Ensure your GOPATH_ is properly set.
4. Download it:

::

    go get github.com/tkliner/stats


Usage
-----

Fasthttp and Fasthttprouter
.......

If you are using fasthttp_ with fasthttprouter_ you need to call the middleware with the handler itself:

.. code-block:: go
    
    package main                                                                          

    import (
            "encoding/json"
            "github.com/buaazp/fasthttprouter"
	        "github.com/valyala/fasthttp"
            "github.com/tkliner/stats"
    )
    
    func main() {
        router := fasthttprouter.New()
	    s := stats.New()
	    router.GET("/stats", func(ctx *fasthttp.RequestCtx) {
			ctx.Response.Header.Set("Content-Type", "application-json")
			s, err := json.Marshal(s.Data())
			if err != nil {
					log.Fatal("Stats error")
			}
			ctx.Write(s)
			ctx.SetStatusCode(200)
	    })
	    log.Fatal(fasthttp.ListenAndServe("localhost:8000", s.Handler(router.Handler)))
    }

Run it in a shell:

::

    $ go run server.go

Then in another shell run:

::

    $ curl http://localhost:3000/stats | python -m "json.tool"

Expect the following result:

.. code-block:: json

    {
        "total_response_time": "1.907382ms",
        "average_response_time": "86.699\u00b5s",
        "average_response_time_sec": 8.6699e-05,
        "count": 1,
        "pid": 99894,
        "status_code_count": {
            "200": 1
        },
        "time": "2015-03-06 17:23:27.000677896 +0100 CET",
        "total_count": 22,
        "total_response_time_sec": 0.0019073820000000002,
        "total_status_code_count": {
            "200": 22
        },
        "unixtime": 1425659007,
        "uptime": "4m14.502271612s",
        "uptime_sec": 254.502271612
    }

See `examples <https://github.com/thoas/stats/blob/master/examples>`_ to
test them.


Inspiration
-----------

`Antoine Imbert <https://github.com/ant0ine>`_ is the original author
of this middleware.

Originally developed for `go-json-rest <https://github.com/ant0ine/go-json-rest>`_,
it had been ported as a simple Golang handler by `Florent Messa <https://github.com/thoas>`_
to be used in various frameworks.

This middleware implements a ticker which is launched every seconds to
reset requests/sec and will implement new features in a near future :)

.. _GOPATH: http://golang.org/doc/code.html#GOPATH
.. _StatusMiddleware: https://github.com/ant0ine/go-json-rest/blob/master/rest/status.go
.. _go-json-rest: https://github.com/ant0ine/go-json-rest
.. _negroni: https://github.com/codegangsta/negroni
.. _martini: https://github.com/go-martini/martini
.. _picfit: https://github.com/thoas/picfit
.. _HTTPRouter: https://github.com/julienschmidt/httprouter

Original package
----------------

This is fork of the original package `thoas/stats <https://github.com/thoas/stats>`_, which was created as part of modification to be used with fasthttp and fasthttprouter