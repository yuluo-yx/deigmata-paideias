curl "http://127.0.0.1:9180/apisix/admin/routes/1" -H "X-API-KEY: 054f7cf07e344346cd3f287985e76a21" -X PUT -d '{"methods": ["GET"],"uri": "/api/test","upstream": {"type": "roundrobin","nodes": {"127.0.0.1:8080": 1}}}'

