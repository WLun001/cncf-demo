kubectl create secret generic google-dns-sa \
  --from-file=key.json=dns.json -n traefik
