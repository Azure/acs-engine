var querystring = require('querystring');
var UUID = require("uuid-1345");

module.exports = function (context, req) {
    /*
        ?subscription_id=foo
        &resource_group=foo
        &role_id=foo
        &vm_oid={vm}_{oid}
    */

    var subscriptionId = req.query["subscription_id"];
    var resourceGroup = req.query["resource_group"];
    var roleDefinitionId = req.query["role_id"];
    var vmName = req.query["vm_name"];
    var principalId = req.query["principal_id"];

    var scope = "/subscriptions/"+subscriptionId+"/resourceGroups/"+resourceGroup;
    var outputs = [];
    var guid = UUID.v5({
        namespace: UUID.namespace.url,
        name: principalId,
    });

    var entry = `{
        "apiVersion": "2014-07-01-preview",
        "type": "Microsoft.Authorization/roleAssignments",
        "name": "` + guid + `",
        "properties": {
            "roleDefinitionId": "` + roleDefinitionId + `",
            "principalId": "` + principalId  + `",
            "scope": "` + scope + `"
        }
    }`;

    content = `{
        "$schema": "https://schema.management.azure.com/schemas/2015-01-01/deploymentTemplate.json#",
        "contentVersion": "1.0.0.0",
        "parameters": {},
        "variables": {},
        "resources": [ ` + entry + `]
    }`;
    context.log(content);
    content = JSON.parse(content);
    content = JSON.stringify(content);

    context.res = {
        status: 200,
        body: content,
        headers: {
            'Content-Type': 'application/json'
        },
        isRaw: true
    };

    context.done();
};
