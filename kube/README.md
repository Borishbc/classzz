# classzz Kubernetes Configs

The following guide will walk you through creating a classzz full node within GKE (Google Container Engine).

This node will have both transaction and address indexing turned on. If you don't need these features you can edit `kube/classzz-deployment` and update the flags passed to classzz.

Steps:
1. Add a new blank disk on GCE called `classzz-data` that is 300GB. You can always expand it later.
2. Change the `rpcuser` and `rpcpass` values in `classzz-secrets.yml`. They are base64 encoded. To base64 a string, just run `echo -n SOMESTRING | base64`.
3. Run `kubectl create -f /path/to/kube`
4. Lookup the `classzz-srv` service in the web-ui to get your public ip.
5. Profit!
