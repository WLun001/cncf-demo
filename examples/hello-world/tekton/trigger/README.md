## GitHub EventListener

Creates an EventListener that listens for GitHub webhook events.

### Try it out locally:

1. To create the GitHub trigger and all related resources, run:

   ```bash
   kubectl apply -f .
   ```

2. Port forward:

   ```bash
   kubectl port-forward service/el-github-listener 8080
   ```

3. Test by sending the sample payload.

   ```bash
   curl -v \
   -H 'X-GitHub-Event: push' \
   -H 'X-Hub-Signature: sha1=87b1adbb9aca10522739f9f94d372afd1542e498' \
   -H 'Content-Type: application/json' \
   -d '{"ref": "refs/heads/main", "repository": {"git_url": "https://github.com/WLun001/cncf-demo.git"}}' \
   http://localhost:8080
   ```

The response status code should be `202 Accepted`

[`HMAC`](https://www.freeformatter.com/hmac-generator.html) tool used to create X-Hub-Signature.

You should see a new PipelineRun that got created:
