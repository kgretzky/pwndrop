****v1.0.2****
- [x] Unauthorized connections not pointing to any hosted files, returning 404, are now automatically closed instead of being kept alive. This resolves the issue of pwndrop getting DDoSed quickly with bots hammering requests at it from various sources.
- [x] Anti-DDoS feature has been added, which temporarily blacklists every IP address of a client who made 10 consecutive requests returning 404. Blacklist period is currently 10 minutes.
- [x] Removed timeouts for uploading and downloading files fully. The previous 15 minutes timeout would have not helped with DDoS attacks anyway. 

****v1.0.1****
- [x] Increased the time limit for uploads and downloads from 15 seconds to 15 minutes. Should fix the issue of uploads/downloads being interrupted on slow connections, when handling big files.