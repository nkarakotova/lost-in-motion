import http3

r = http3.get("https://localhost:8443/status", verify=False)

print(r)