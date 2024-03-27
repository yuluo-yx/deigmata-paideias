curl "http://127.0.0.1:9180/apisix/admin/routes/1" -H "X-API-KEY: 054f7cf07e344346cd3f287985e76a21" -X PUT -d '{"methods": ["GET"],"uri": "/front/*","upstream": {"type": "roundrobin","nodes": {"192.168.2.27": 1}}}'

