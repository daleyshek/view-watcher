# View-Watcher 记录页面驻留时长

需要  

- Redis
- Nginx

## Redis

默认使用database 2,无授权

## Nginx 配置

```conf
server {
    # your server
    # ...

    # watcher location
    location /watcher {
        proxy_set_header X-Forwarded-For $remote_addr;
        proxy_set_header Host            $http_host;

        proxy_pass http://127.0.0.1:8844;
    }
}
```
