1. More fine-grained errors (e.g. which cities could not be found in request)
2. Extend storage to store also information about cities, which couldn't be found.
3. Extend functions with context.Context, read request id, add logging based with request ID
4. If we want more instances of the service: add another source of forecasts - e.g. Redis. This would allow to share cache between instances
5. If we want a different format in the API (different from what is served by open weather website) - then we need to create separate structures for a) Open weather unmarshalling b) serving our custom structure