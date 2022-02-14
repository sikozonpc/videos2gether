The goal of this experimental project is to create an application to provide realtime video watching, 
so you can watch videos together with your friends at the same time.

## How it works

Firstly someone must create a room, after that by sharing the link of the room everyone can join in and watch and request videos.

Anyone can pause the video and add more videos to the queue since everyone is an "admin" in the room, this solution aims  to provide
freedom for everyone in the room and there are no curent plans to add authorization permitions to rooms.

### Run

```bash
# Create a .env if needed and run:
docker-compose up
```

### Connecting to the Redis instance

In order to debug the Redis data store we must connect to it.

```bash
$ docker exec -it <redis_image-id> /bin/sh
$ redis-cli -h <endpoint> -p <port> -a <password>

$ KEYS * # get all rooms by id
```


### Resources

- https://docs.redis.com/latest/rc/rc-quickstart/