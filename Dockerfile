#
# Copyright (c) 2012-2018 Red Hat, Inc.
# This program and the accompanying materials are made
# available under the terms of the Eclipse Public License 2.0
# which is available at https://www.eclipse.org/legal/epl-2.0/
#
# SPDX-License-Identifier: EPL-2.0
#
# Contributors:
#   Red Hat, Inc. - initial API and implementation
#

FROM golang:1.13-alpine as builder
RUN apk add --no-cache ca-certificates
RUN adduser -D -g '' appuser
WORKDIR /go/src/github.com/skabashnyuk/json-rpc-loader
COPY . /go/src/github.com/skabashnyuk/json-rpc-loader
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags '-w -s' -installsuffix cgo -o json-rpc-loader github.com/skabashnyuk/json-rpc-loader


FROM alpine:3.10
USER appuser
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /go/src/github.com/skabashnyuk/json-rpc-loader /usr/local/bin
ENTRYPOINT ["json-rpc-loader"]