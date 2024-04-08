# A tiny blog from the tiny web framework
I wanted to make a small site for technical blog-posts etc.
I sought up the most contemporary web-stack of 2023-24, and found the following:
[What is the ideal Tech stack to build a website in 2024?](https://dev.to/jakemackie/web-development-in-2024-29d6)

Disappointed to see it was all Vue.js and React - not that there is anything wrong with those approaches, 
I just thought it was a bit much for my tiny site. The comment section however did not disappoint. Full of many D.I.Y. 
suggestions for tiny web-frameworks - especially the 
[Document Markup Library (DML)](https://dml.efpage.de/DML_homepage/index.html) caught my eye.
However I still felt something was missing, so here is a tiny-framework, with custom tags and typescript:

To begin with, typescript in a very vanilla form is needed. Either download it from 
[nodejs.org](https://nodejs.org/en/learn/getting-started/how-to-install-nodejs) or as i do on a mac with brew:
```shell
# on MacOs
brew install node
```
Then install typescript globally:
```shell
npm install typescript -g
```
after that, the command `tsc <file.ts>` will transpile a single ts-file to plain js, or as I do:
```shell
$\[checkout]\static\js\ts\> npx tsc
``` 
which transpiles ts files from the root of where the `tsconfig.json`-file is placed (in this case `\static\js\ts\`).

The resulting js-files are located in the `static\js\ts\dist` folder, which the html should refer to 
(when including `<script src="js/ts/dist/tiny-fw.js"></script>`-tags).  

## Serving the web-content
__**Depends on go-lang**__ (install go-lang)

A minimalistic go-lang application has been created for serving the files and the metric endpoint.

To build and run the go-lang http server do: 
```shell
# Env-var directing tiny-blog to the config file that configures the location of static web content and servlet port
# (The content of this config is not reloaded at run-time)
export SERVICE_CONFIG=./configs/default-service-config.yaml
go mod tidy
go build -o tiny-blog
./tiny-blog
```

# Container
The `Dockerfile` includes two stages [`build`,`app`], where `build` makes the golang-binaries and `app` is the final 
image, where users, static content and binaries are copied to.

Build and run as follows
```shell
docker build . --build-arg="BUILDPLATFORM=linux/arm64" \
  --build-arg="TARGETARCH=arm64" \
  -t favorite.registry.com/tiny-blog:0.2.0

docker run -it --publish=3000:3000 \                                                                                  
   favorite.registry.com/tiny-blog:0.2.0
```

Here is how to build containers for [multiple platforms](https://docs.docker.com/build/building/multi-platform/). 

## Metrics

The simplest of prometheus metrics are used: https://prometheus.io/docs/guides/go-application/
```shell
go get github.com/prometheus/client_golang/prometheus/promhttp
```
together with:
```go
import (
...
"github.com/prometheus/client_golang/prometheus/promhttp"
)
...
mux.Handle("/metrics", promhttp.Handler())
```

## Traces with open telemetry

```shell
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/trace
go get go.opentelemetry.io/otel/sdk
go get go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc
```
together with:



# Acknowledgments
## Style
I am currently borrowing styles from: https://smartblogger.com/blog-design/ 
## The markdown feature
Acknowledgements to [Adam Leggett](https://github.com/adamvleggett) - look to the Notes.md in the `drawdown` folder.
## 404 folders
Filtering out folders when serving static content: https://www.alexedwards.net/blog/disable-http-fileserver-directory-listings
## favicon.ico
Free and royalty-free:
https://www.flaticon.com/free-icon/favicon_7710476
