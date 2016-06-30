# doshin

External monitoring agent for prometheus alertmanager

## Build

    git clone git@github.com:tokuhirom/doshin.git $GOPATH/src/github.com/doshin
    cd $GOPATH/src/github.com/doshin/cli
    go get
    go build

## Run

    doshin -c config.yml

## Configuration

You can configure this application by YAML file.

```
---
alert_manager:
  url: http://localhost:9093/api/v1/alerts
  labels:
    service: my-great-webapp
watch:
  http:
    targets:
      - url: http://example.com/
        interval: 10s
      - url: http://example.cm/
        interval: 10s
  net:
    targets:
      - network: tcp
        address: example.com:80
        interval: 10s
      - network: tcp
        address: example.com:22
        interval: 10s
```

## LICENSE

```
The MIT License (MIT)
Copyright © 2016 Tokuhiro Matsuno, http://64p.org/ <tokuhirom@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the “Software”), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
```

