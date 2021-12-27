# CapaSea

## Capa HTTP runner

Submits a file to be run across a standard capa ruleset. Reports are served seperately, with the specific uri path returned in json. It supports both manually uploading a file, as well as posting to "/upload" using standard multiform format. 

This application is most suitable for small binaries and more simple jobs. It is recommended that more complex capa jobs be run outside of CapaSea. This is due to the ambiguity of if it's working or not on larger/complex binaries.

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

HTTP is served over port 8080.
