# Microsoft Azure Container Service Engine - Kubernetes AAD integration Walkthrough

This is walkthrough is to help you get start with Azure Active Directory(AAD) integeration with an ACS-Engine Kubernetes cluster.

[OpenID Connect](http://openid.net/connect/) is a simple identity layer built on top of the OAuth 2.0 protocol, and it is supported by both AAD and Kubernetes. Here we're going to use OpenID Connect as the communication protocol.

Please also refer to [Azure Active Directory plugin for client authentication](https://github.com/kubernetes/kubernetes/blob/master/staging/src/k8s.io/client-go/plugin/pkg/client/auth/azure/README.md) in Kubernetes repo for more details abount OpenID Connect and AAD support in upstream.

## Prerequision
1. An Azure Active Directory tenant, will refer as `AAD Tenant`. You can use the tenant for your Azure subscription;
2. A `Web app / API` type AAD application, will refer as `Server Application`. This application represents the `apiserver`;
3. A `Native` type AAD application, will refer as `Client Application`. This application is for user login via `kubectl`. You'll need to add delegated permission to `Server Application`, please see [troubleshooting](#loginpageerror) section for detail.

## Deployment
Follow the [deployment steps](kubernetes.md#deployment). In step #4, add the following under 'properties' section:
```
"aadProfile": {
    "serverAppID": "",
    "clientAppID": "",
    "tenantID": ""
}
```

- `serverAppID`   : the `Server Application`'s ID
- `clientAppID`   : the `Client Application`'s ID
- `tenantID`      : (optional) the `AAD tenant`'s ID. If not specified, will use the tenant of the deployment subscription.

After template generation, the local generated kubeconfig file (`_output/<instance>/kubeconfig/kubeconfig.<location>.json`) will have the default user using AAD.
Initially it isn't assoicated with any AAD user yet. To get started, try any kubectl command (like `kubectl get pods`), and you'll be prompted to the device login process. After login, you will be able to operate the cluster using your AAD identity.

### Note
Please note that as of Kubernetes 1.7, the default is authorization mode is `AlwaysAllow`, which means any authenticated user have full access of the cluster.
OpenID Connect is an authentication protocol responsible for identify users only, so initally all active accounts under the tenant will be able to login and have full admin privilege of the cluster.

In this case you may want to also turn on RBAC for your cluster.
Please refer to [Enable Kubernetes Role-Based Access Control](features.md#optional-enable-kubernetes-role-based-access-control-rbac) for turing on RBAC using acs-engine.

Following instructions are for turnning on RBAC manually together with AAD integration:

1. Since we use AAD object ID as OpenID Connect identity.
    You'll first need to figure out your account's object ID. Here is how to do it using Azure Portal:
    Navigate to `Azure Active Directory` -> `Users and groups` -> `All users`. And choose your account in right pannel. Switch to `Manage` -> `Profile`, and you can see the `Object ID` property.
2. Figure out your user name. The user name would be in form of `IssuerUrl#ObjectID` format.
    You can navigate to `https://login.microsoftonline.com/{tenantid}/.well-known/openid-configuration`, and find the `IssuerUrl` under `issuer` property.
3. Add your account as admin role
```
kubectl create clusterrolebinding aad-default-cluster-admin-binding --clusterrole=cluster-admin --user={UserName}
```
For example, if your `IssuerUrl` is `https://sts.windows.net/e2917176-1632-47a0-ad18-671d485757a3/`, and your `ObjectID` is `22fa281b-bf62-4b14-972c-0dbca24a25a2`, the command would be:
```
kubectl create clusterrolebinding aad-default-cluster-admin-binding --clusterrole=cluster-admin --user=https://sts.windows.net/e2917176-1632-47a0-ad18-671d485757a3/#22fa281b-bf62-4b14-972c-0dbca24a25a2
```

4. Turn on RBAC on master nodes.
    On master nodes, edit `/etc/kubernetes/manifests/kube-apiserver.yaml`, add `--authorization-mode=RBAC` under `command` property. Reboot nodes.
5. Now that AAD account will be cluster admin, other accounts can still login but do not have permission for operating the cluster.
    To verify this, add another client user:
    ```
    kubectl config set-credentials "user1" --auth-provider=azure \
    --auth-provider-arg=environment=AzurePublicCloud \
    --auth-provider-arg=client-id={ClientAppID} \
    --auth-provider-arg=apiserver-id={ServerAppID} \
    --auth-provider-arg=tenant-id={TenantID}
    ```

    And use that user to login
    ```
    kubectl get pods --user=user1
    ```
    Now you'll be prompted to login again, you can try logining with another AAD user account. 
    The login would succeed, but later you can see following message since server denies access:
    ```
    Error from server (Forbidden): User "https://sts.windows.net/{tenantID}/#{objectID}" cannot list pods in the namespace "default". (get pods)
    ```

    You can manually update server configuration or add administrator users based on your requirement.

## Troubleshooting

### LoginPageError
If you failed in login page, you may see following error message
```
Invalid resource. The client has requested access to a resource which is not listed in the requested permissions in the client's application registration. Client app ID: {UUID} Resource value from request: {UUID}. Resource app ID: {UUID}. List of valid resources from app registration: {UUID}.
```
This could be caused by `Client Application` not authorized.

1. Go to Azure Portal, navigate to `Azure Active Directory` -> `App registrations`.
2. Select the `Client Application`, Navigate to `Settings` -> `Required permissions`
3. Choose `Add`, select the `Server Application`. In permissions tab, select `Delegated permissions` -> `Access {Server Application}`

### ClientError
If you see following message return from server via `kubectl`
```
Error from server (Forbidden)
```

It is usually caused by an incorrect configuration. You could find more debug information in apiserver log. On a master node, run following command:
```
docker logs -f $(docker ps|grep 'hyperkube apiserver'|cut -d' ' -f1) 2>&1 |grep -a auth
```

You might see following message like this:
```
Unable to authenticate the request due to an error: [invalid bearer token, [crypto/rsa: verification error, oidc: JWT claims invalid: invalid claims, 'aud' claim and 'client_id' do not match, aud=UUID1, client_id=spn:UUID2]]
```
This indicates server and client is using different `Server Application` ID, could usually happen when the configurations being updated manually.

For other auth issues, you may also find some useful information from the log.
