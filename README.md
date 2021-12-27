# CapaSea

## Capa HTTP runner

Submits a file to be run across a standard capa ruleset. Reports are served seperately, with the specific uri path returned in json.

This was developed for an internal testing setup. If there are features you want, create an issue with the details. If I feel like it fits the spirit of this application, I'll happily add it.

### Example Submission Code (Python)
```
import requests

files = {'myFile': ("example", open("/path/to/file", 'rb'))}
res = requests.post("http://<ip_addr>:8080/upload", files=files)
if res.status_code == 200:
    data = res.json()
    print(data)

```
