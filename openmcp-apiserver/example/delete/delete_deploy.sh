USERNAME="openmcp"
PASSWORD="keti"
IP="10.0.3.20"
PORT="30000"
NS="default"
POD="example-nginx-deploy"
URL="apis/apps/v1/namespaces/$NS/deployments/$POD"
CLUSTER="openmcp"


echo -n | openssl s_client -connect $IP:$PORT | sed -ne '/-BEGIN CERTIFICATE-/,/-END CERTIFICATE-/p' > server.crt

TOKEN_JSON=`curl -XPOST \
        --cacert server.crt \
        --insecure \
        -H "Content-type: application/json" \
        --data "{\"username\":\"$USERNAME\",\"password\":\"$PASSWORD\"}" \
        https://$IP:$PORT/token`

TOKEN=`echo $TOKEN_JSON | jq .token`
TOKEN=`echo "$TOKEN" | tr -d '"'`

curl -X DELETE --cacert server.srt --insecure -H "Authorization: Bearer $TOKEN" https://$IP:$PORT/$URL?clustername=$CLUSTER

rm server.crt
