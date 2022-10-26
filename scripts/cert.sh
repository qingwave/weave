#!/bin/bash

CERTDIR="$(cd "$( dirname "${BASH_SOURCE[0]}" )/../certs" && pwd )"

pushd $CERTDIR

base_subj="/C=CN/ST=Xi'an/L=Xi'an/O=weave/OU=weave"
root_subj="${base_subj}/CN=root"
frontend_subj="${base_subj}/CN=frontend"
server_subj="${base_subj}/CN=server"

# generate root cert
openssl genrsa -out root.key 4096

openssl req -new -x509 -days 100000 -key root.key -out root.crt -subj ${root_subj}

# generate frontend cert
openssl genrsa -out frontend.key 2048

openssl req -new -key frontend.key -out frontend.csr -subj ${frontend_subj}

openssl x509 -req -days 100000 -CA root.crt -CAkey root.key -in frontend.csr -out frontend.crt -CAcreateserial

# generate backend cert
openssl genrsa -out server.key 2048

openssl req -new -key server.key -out server.csr -subj ${server_subj}

openssl x509 -req -days 100000 -CA root.crt -CAkey root.key -in server.csr -out server.crt -CAcreateserial

popd
