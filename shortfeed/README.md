# ShortFeed

## DNS Coffee API Example

This is an example API client to list dropped short domains from the [DNS Coffee API](https://api.dns.coffee/doc/)

### Building

```shell
go build
```

### Running

```shell
./shortfeed -len 5 -apikey YOUR_KEY_HERE
```

### Usage

```shell
Usage of ./shortfeed:
  -apikey string
        DNS coffee API Key. REQUIRED
  -date string
        the date to fetch YYYY-MM-DD (default "2022-06-07")
  -endpoint string
        endpoint to use (default "https://api.dns.coffee/api/v0/feeds/domains/old/date")
  -ingore-tlds string
        comma separated list of TLDs to ignore
  -len uint
        minimum domain length to return (default 5)
```

