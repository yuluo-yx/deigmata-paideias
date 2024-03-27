curl "http://127.0.0.1:9180/apisix/admin/upstreams/1" \
-H "X-API-KEY: 054f7cf07e344346cd3f287985e76a21" -X PUT -d '
{
  "type": "roundrobin",
  "nodes": {
    "httpbin.org:80": 1
  }
}'
