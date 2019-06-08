## Intro

VideoCoin command line application to access the VideoCoin Network.

## Build

Build a binary:

```
make cli
```

Build a mini rtmp server (optional):

```
make mrtmp
```

Build a docker container (optional):

```
docker build . -t cli
```

## Usage

For demo purpose we are going to run mini rtmp server we built above:

```
./build/minirtmp
rtmp server is listening on 192.168.86.107:1936
```

Start streaming to mini rtmp server:

```
ffmpeg -f lavfi -i anullsrc=channel_layout=stereo:sample_rate=44100 -f avfoundation -pix_fmt uyvy422 -i "Capture screen 0" -vcodec libx264 -pix_fmt yuv420p -preset medium -r 30 -g 30 -b:v 2500k -acodec libmp3lame -ar 44100 -threads 6 -b:a 712000 -bufsize 512k -f flv rtmp://192.168.86.107:1936/stream
```

Run a `cli`:

```
build/cli start rtmp://127.0.0.1:1936/stream -a $(ACCOUNT_FILE_PATH) -p $(ACCOUNT_PASSWORD)
```

Or docker container:

```
docker run -v $(PWD)/$(ACCOUNT_FILE_PATH):/account cli start rtmp://192.168.86.107:1936/stream -a account -p $(ACCOUNT_PASSWORD)
```
