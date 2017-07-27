# Porting a new DC/OS version to ACS-Engine

## 1. Locate the official ARM Template

Go to `https://dcos.io/docs/X.X/administration/installing/cloud/azure/`, where `X.X` should be replaced by the version you are looking to port.  
In the documentation, you will the link to the ARM templates you are looking for.  
The latest stable templates should be at `https://downloads.dcos.io/dcos/stable/azure.html`  
Early Access at: `https://downloads.dcos.io/dcos/EarlyAccess/azure.html`  
Etc.

## 2. Find the package GUIDs

Following the previous step, you should now have 3 ARM templates (1, 3 and 5 masters variants).  
We now need to find the package GUID of each variant.  
In each template you should find a string that looks like: `dcos-config--setup_<Some GUID>`, this GUID is what we are looking for.
Extract the GUIDs from the 3 differents templates, and them in `engine.go/getPackageGUID` for your specific DC/OS version.


## 3. Extract the cloud-config data from the template

In one of the template (no matter which one), grab the data from the MasterVM.osProfile.customData.  
If you remove the concat operation, you should end up which a big string of unescaped JSON.  
Escape it (for example using this [online tool](http://www.freeformatter.com/javascript-escape.html#ad-output)), and convert it to yaml (you can use [json2yaml](https://www.json2yaml.com/)).
You should now have a clean yaml.

## 4. Create and customize the custom data file.

under the `parts` directory, create a new file called `dcoscustomdataXXX.t` replacing `XXX` by the correct version number.  
Paste the yaml from the previous step inside.  

In the new file, under the `runcmd` section you should find 4 sucessive `curl` calls downloading some `.deb` packages followed by a bash script installing each one of them. This is handled by `parts\dcosprovision.sh` in ACS-Engine, so make sure the dependencies didn't change and replace the `curl` and `bash` calls by a link to the script.

For example, in DC/OS 1.9:  
```yaml
- curl -fLsSv --retry 20 -Y 100000 -y 60 -o /var/tmp/1.deb https://az837203.vo.msecnd.net/dcos-deps/libipset3_6.29-1_amd64.deb
- curl -fLsSv --retry 20 -Y 100000 -y 60 -o /var/tmp/2.deb https://az837203.vo.msecnd.net/dcos-deps/ipset_6.29-1_amd64.deb
- curl -fLsSv --retry 20 -Y 100000 -y 60 -o /var/tmp/3.deb https://az837203.vo.msecnd.net/dcos-deps/unzip_6.0-20ubuntu1_amd64.deb
- curl -fLsSv --retry 20 -Y 100000 -y 60 -o /var/tmp/4.deb https://az837203.vo.msecnd.net/dcos-deps/libltdl7_2.4.6-0.1_amd64.deb
- sed -i "s/^Port 22$/Port 22\nPort 2222/1" /etc/ssh/sshd_config
- service ssh restart
- bash -c "try=1;until dpkg -i /var/tmp/{1,2,3,4}.deb || ((try>9));do echo retry \$((try++));sleep
  \$((try*try));done"
```  

becomes   

```yaml
- /opt/azure/containers/provision.sh
```

Additional modifications under `runcmd`:
* the `content` of the cmd with path `/etc/mesosphere/setup-flags/cluster-packages.json` becomes `'DCOS_ENVIRONMENT={{{targetEnvironment}}}'`
* Replace every occurence of the Package GUID (that we found in step 2) by `DCOSGUID`.
* the `content` of the cmd with path `/etc/mesosphere/setup-flags/late-config.yaml` should be modified to accept ACS-Engine bindings instead of variable where needed (look at a previous custom data file for reference).  
* At the very end of the file, replace  
```yaml
- content: ''
  path: "/etc/mesosphere/roles/master"
- content: ''
  path: "/etc/mesosphere/roles/azure_master"
- content: ''
  path: "/etc/mesosphere/roles/azure"
```  

by   

```yaml
- content: ''
  path: /etc/mesosphere/roles/azure
- content: 'PROVISION_STR'
  path: "/opt/azure/containers/provision.sh"
  permissions: "0744"
  owner: "root"
```