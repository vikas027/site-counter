# Site Counter
This is a simple app which prints out the hostname, IP Addresses and the number of hits on a HTTP call.

The app is available on [Docker Hub Registry](https://cloud.docker.com/repository/docker/vikas027/site-counter)

# Usage
```bash
$ docker run -dit -p <host_port>:80 vikas027/site_counter
```

# Sample Outputs
```bash
$ curl localhost:8080/
Use one of these URIs: [health counter]
$
$ curl localhost:<host_port>/health
{"status":"ok"}
$ curl localhost:<host_port>/counter
a8c7e2217271  -  [172.17.0.2]  -  View Count:  1
$ curl localhost:<host_port>/counter
a8c7e2217271  -  [172.17.0.2]  -  View Count:  2
```

# To Do
- Reduce size of the docker image

# References
- https://github.com/generalhenry/go-redis-counter
- https://github.com/docker-library/redis
