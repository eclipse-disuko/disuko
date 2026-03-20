openssl genrsa -out ../../disco.key 2048
openssl ecparam -genkey -name secp384r1 -out ../../disco.key
openssl req -new -x509 -sha256 -key ../../disco.key -out ../../disco.cert -days 3650
