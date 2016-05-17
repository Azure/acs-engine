#!/usr/bin/python
import base64
import os
import gzip
import StringIO
import sys
import shutil
import json
import argparse

def buildb64GzipStringFromFile(file):
    # read the script file
    with open(file) as f:
        content = f.read()
    compressedbuffer=StringIO.StringIO()

    # gzip the script file
    # mtime=0 sets a fixed timestamp in GZip header to the Epoch which is January 1st, 1970
    # Make sure it doens't change unless the stream changes 
    with gzip.GzipFile(fileobj=compressedbuffer, mode='wb', mtime=0) as f:
        f.write(content)
    b64GzipStream=base64.b64encode(compressedbuffer.getvalue())

    return b64GzipStream

# Function reads the files from disk,
# and embeds it in a Yaml file as a base-64 enconded string to be
# executed later by template
def buildYamlFileWithWriteFiles(files):
    gzipBuffer=StringIO.StringIO()

    clusterYamlFile="""#cloud-config

write_files:
%s
"""
    writeFileBlock=""" -  encoding: gzip
    content: !!binary |
        %s
    path: /opt/azure/containers/%s
    permissions: "0744"
"""
    filelines=""
    for encodeFile in files:
        b64GzipString = buildb64GzipStringFromFile(encodeFile)
        filelines=filelines+(writeFileBlock % (b64GzipString,encodeFile))

    return clusterYamlFile % (filelines)

# processes a Yaml file to be included properly in ARM template
def convertToOneArmTemplateLine(clusterYamlFile):
    # remove the \r\n and include \n in body and escape " to \"
    return  clusterYamlFile.replace("\n", "\\n").replace('"', '\\"')

# Loads the base ARM template file and injects the Yaml for the shell scripts into it.
def processBaseTemplate(baseTemplatePath,
                        clusterInstallScript,
                        jumpboxTemplatePath = None,
                        linuxJumpboxInstallScript = None,
                        swarmWindowsAgentInstallScript = None,
                        additionalFiles = [],
                        windowsAgentDiagnosticsExtensionTemplatePath = None):

    #String to replace in JSON file
    CLUSTER_YAML_REPLACE_STRING  = "#clusterCustomDataInstallYaml"
    JUMPBOX_FRAGMENT_REPLACE_STRING = "#jumpboxFragment"
    JUMPBOX_FQDN_REPLACE_STRING = "#jumpboxFQDN"
    JUMPBOX_LINUX_YAML_REPLACE_STRING = "#jumpboxLinuxCustomDataInstallYaml"
    SWARM_WINDOWS_AGENT_CUSTOMDATA_REPLACE_STRING = "#swarmWindowsAgentCustomData"
    WINDOWS_AGENT_DIAGNOSTICS_EXTENSION_REPLACE_STRING = "#windowsAgentDiagnosticsExtension"
   
    # Load Base Template
    armTemplate = []
    with open(baseTemplatePath) as f:
        armTemplate = f.read()
        
    # All templates have vmsizemapping.  Add it to base template
    ARM_INPUT_VMSIZE_MAPPING_TEMPLATE = "vmsizes-storage-account-mappings.json"
    VMSIZE_MAPPINGS_STRING = "#vmsizemapping"

    vmsizeMappings = ""    
    with open(ARM_INPUT_VMSIZE_MAPPING_TEMPLATE) as f:
       vmsizeMappings = f.read()

    armTemplate = armTemplate.replace(VMSIZE_MAPPINGS_STRING, vmsizeMappings)

    # Generate cluster Yaml file for ARM
    clusterYamlFile = convertToOneArmTemplateLine(buildYamlFileWithWriteFiles([clusterInstallScript]+additionalFiles))
    armTemplate = armTemplate.replace(CLUSTER_YAML_REPLACE_STRING, clusterYamlFile)

    # Add Jumpbox YAML, ARM and FQDN Fragment if jumpboxTemplatePath is defined
    jumpboxTemplate = ""
    jumpboxFQDN = ""
    linuxJumpboxYamlFile = ""
    swarmWindowsAgentCustomData = ""
    windowsAgentDiagnosticsExtension = ""

    if jumpboxTemplatePath != None :
        # Add Jumpbox FQDN Fragment if jumpboxTemplatePath is defined
        jumpboxFQDN = "[reference(concat('Microsoft.Network/publicIPAddresses/', variables('jumpboxPublicIPAddressName'))).dnsSettings.fqdn]"

        with open(jumpboxTemplatePath) as f:
            jumpboxTemplate = f.read()

        # Generate jumpbox Yaml file for ARM
        if linuxJumpboxInstallScript != None :
            # the linux jumpbox does not need the nginx configuration file
            linuxJumpboxYamlFile = convertToOneArmTemplateLine(buildYamlFileWithWriteFiles([linuxJumpboxInstallScript]))

    # Add windows agent install script if passed in 
    if swarmWindowsAgentInstallScript != None:
        swarmWindowsAgentCustomData = buildb64GzipStringFromFile(swarmWindowsAgentInstallScript)
            
    # Add windows IaaS Diagnsotics extension JSON payload if passed in
    if windowsAgentDiagnosticsExtensionTemplatePath != None :
        with open(windowsAgentDiagnosticsExtensionTemplatePath) as f:
            windowsAgentDiagnosticsExtension = f.read()

    # Want these to be replaced with blank strings if not defined
    armTemplate = armTemplate.replace(JUMPBOX_FRAGMENT_REPLACE_STRING, jumpboxTemplate)
    armTemplate = armTemplate.replace(JUMPBOX_FQDN_REPLACE_STRING, jumpboxFQDN)
    armTemplate = armTemplate.replace(JUMPBOX_LINUX_YAML_REPLACE_STRING, linuxJumpboxYamlFile)
    armTemplate = armTemplate.replace(SWARM_WINDOWS_AGENT_CUSTOMDATA_REPLACE_STRING, swarmWindowsAgentCustomData)
    armTemplate = armTemplate.replace(WINDOWS_AGENT_DIAGNOSTICS_EXTENSION_REPLACE_STRING, windowsAgentDiagnosticsExtension)

    # Make sure the final string is valid JSON
    try:
        json_object = json.loads(armTemplate)
    except ValueError, e:
        print e
        errorFileName = baseTemplatePath + ".err"
        with open(errorFileName, "w") as f:
            f.write(armTemplate)
        print "Invalid armTemplate saved to: " + errorFileName
        raise

    return armTemplate;

if __name__ == "__main__":
    # Parse Arguments
    parser = argparse.ArgumentParser()
    parser.add_argument("-o", "--output_directory",  help="Directory to write templates files to.  Default is current directory.")
    parser.add_argument("-wpf", "--write_parameter_files", help="Write separate parameter file for each template.  Default is false.",
                        action="store_true" )

    args = parser.parse_args()

    if (args.output_directory == None) :
        args.output_directory = os.getcwd()

    args.output_directory = os.path.expandvars(os.path.normpath(args.output_directory))

    if ( os.path.exists(args.output_directory) == False ):
        os.mkdir(args.output_directory)

    # Input Arm Template Artifacts to be processed in
    # Note:  These files are not useable ARM templates on thier own or valid JSON
    # They require processing by this script.
    ARM_INPUT_TEMPLATE_TEMPLATE                  = "base-template.json"
    ARM_INPUT_PARAMETER_TEMPLATE                 = "base-template.parameters.json"
    ARM_INPUT_WINDOWS_JUMPBOX_TEMPLATE           = "fragment-windows-jumpbox.json"
    ARM_INPUT_WINDOWS_AGENT_DIAGNOSTICS_TEMPLATE = "fragment-windows-agent-diagnostics-extension.json"
    ARM_INPUT_LINUX_JUMPBOX_TEMPLATE             = "fragment-linux-jumpbox.json"
    ARM_INPUT_SWARM_TEMPLATE_TEMPLATE            = "base-swarm-template.json"
    ARM_INPUT_SWARM_WINDOWS_TEMPLATE_TEMPLATE    = "base-swarm-windows-template.json"

    # Shell Scripts to load into YAML
    MESOS_CLUSTER_INSTALL_SCRIPT = "configure-mesos-cluster.sh"
    SWARM_CLUSTER_INSTALL_SCRIPT = "configure-swarm-cluster.sh"
    LINUX_JUMPBOX_INSTALL_SCRIPT = "configure-ubuntu.sh"
    SWARM_WINDOWS_AGENT_INSTALL_SCRIPT = "Install-ContainerHost-And-Join-Swarm.ps1"

    # admin router configuration file
    ADMIN_ROUTER_CONF = "nginx.conf"

    # Output ARM Template Files.  WIll Also Output name.parameters.json for each
    ARM_OUTPUT_TEMPLATE                                   = "mesos-cluster-with-no-jumpbox.json"
    ARM_OUTPUT_TEMPLATE_WINDOWS_JUMPBOX                   = "mesos-cluster-with-windows-jumpbox.json"
    ARM_OUTPUT_TEMPLATE_LINUX_JUMPBOX                     = "mesos-cluster-with-linux-jumpbox.json"
    ARM_OUTPUT_SWARM_TEMPLATE                             = "swarm-cluster-with-no-jumpbox.json"
    ARM_OUTPUT_SWARM_WINDOWS_TEMPLATE                     = "swarm-cluster-with-windows.json"
    ARM_OUTPUT_SWARM_WINDOWS_TEMPLATE_NO_DIAGNOSTICS      = "swarm-cluster-with-windows-no-diagnostics.json"

    # build the ARM template for jumpboxless
    with open(os.path.join(args.output_directory, ARM_OUTPUT_TEMPLATE), "w") as armTemplate:
        clusterTemplate = processBaseTemplate(
            baseTemplatePath=ARM_INPUT_TEMPLATE_TEMPLATE, 
            clusterInstallScript=MESOS_CLUSTER_INSTALL_SCRIPT,
            additionalFiles=[ADMIN_ROUTER_CONF])
        armTemplate.write(clusterTemplate)

    # build the ARM template for linux jumpbox
    with open(os.path.join(args.output_directory,ARM_OUTPUT_TEMPLATE_LINUX_JUMPBOX), "w") as armTemplate:
        clusterTemplate = processBaseTemplate(
            baseTemplatePath=ARM_INPUT_TEMPLATE_TEMPLATE,
            clusterInstallScript=MESOS_CLUSTER_INSTALL_SCRIPT,
            jumpboxTemplatePath=ARM_INPUT_LINUX_JUMPBOX_TEMPLATE,
            linuxJumpboxInstallScript=LINUX_JUMPBOX_INSTALL_SCRIPT,
            additionalFiles=[ADMIN_ROUTER_CONF])
        armTemplate.write(clusterTemplate)

    # build the ARM template for windows jumpbox
    with open(os.path.join(args.output_directory, ARM_OUTPUT_TEMPLATE_WINDOWS_JUMPBOX), "w") as armTemplate:
        clusterTemplate = processBaseTemplate(
            baseTemplatePath=ARM_INPUT_TEMPLATE_TEMPLATE,
            clusterInstallScript=MESOS_CLUSTER_INSTALL_SCRIPT,
            jumpboxTemplatePath=ARM_INPUT_WINDOWS_JUMPBOX_TEMPLATE,
            additionalFiles=[ADMIN_ROUTER_CONF])
        armTemplate.write(clusterTemplate)

    # build the SWARM ARM template
    with open(os.path.join(args.output_directory, ARM_OUTPUT_SWARM_TEMPLATE), "w") as armTemplate:
        clusterTemplate = processBaseTemplate(
            baseTemplatePath=ARM_INPUT_SWARM_TEMPLATE_TEMPLATE,
            clusterInstallScript=SWARM_CLUSTER_INSTALL_SCRIPT)
        armTemplate.write(clusterTemplate)

    # build the SWARM WINDOWS ARM template with Windows diagnostics extension
    with open(os.path.join(args.output_directory, ARM_OUTPUT_SWARM_WINDOWS_TEMPLATE), "w") as armTemplate:
        clusterTemplate = processBaseTemplate(
            baseTemplatePath=ARM_INPUT_SWARM_WINDOWS_TEMPLATE_TEMPLATE,
            clusterInstallScript=SWARM_CLUSTER_INSTALL_SCRIPT,
            swarmWindowsAgentInstallScript=SWARM_WINDOWS_AGENT_INSTALL_SCRIPT,
            windowsAgentDiagnosticsExtensionTemplatePath=ARM_INPUT_WINDOWS_AGENT_DIAGNOSTICS_TEMPLATE)
        armTemplate.write(clusterTemplate)
        
    # build the SWARM WINDOWS ARM template with NO winodws diagnostics extension
    with open(os.path.join(args.output_directory, ARM_OUTPUT_SWARM_WINDOWS_TEMPLATE_NO_DIAGNOSTICS), "w") as armTemplate:
        clusterTemplate = processBaseTemplate(
            baseTemplatePath=ARM_INPUT_SWARM_WINDOWS_TEMPLATE_TEMPLATE,
            clusterInstallScript=SWARM_CLUSTER_INSTALL_SCRIPT,
            swarmWindowsAgentInstallScript=SWARM_WINDOWS_AGENT_INSTALL_SCRIPT)
        armTemplate.write(clusterTemplate)

    # Write parameter files if specified
    if (args.write_parameter_files == True) :
        shutil.copyfile(ARM_INPUT_PARAMETER_TEMPLATE, os.path.join(args.output_directory, ARM_OUTPUT_TEMPLATE).replace(".json", ".parameters.json") )
        shutil.copyfile(ARM_INPUT_PARAMETER_TEMPLATE, os.path.join(args.output_directory, ARM_OUTPUT_TEMPLATE_LINUX_JUMPBOX).replace(".json", ".parameters.json") )
        shutil.copyfile(ARM_INPUT_PARAMETER_TEMPLATE, os.path.join(args.output_directory, ARM_OUTPUT_TEMPLATE_WINDOWS_JUMPBOX).replace(".json", ".parameters.json") )
        shutil.copyfile(ARM_INPUT_PARAMETER_TEMPLATE, os.path.join(args.output_directory, ARM_OUTPUT_SWARM_TEMPLATE).replace(".json", ".parameters.json") )
        shutil.copyfile(ARM_INPUT_PARAMETER_TEMPLATE, os.path.join(args.output_directory, ARM_OUTPUT_SWARM_WINDOWS_TEMPLATE).replace(".json", ".parameters.json") )
        shutil.copyfile(ARM_INPUT_PARAMETER_TEMPLATE, os.path.join(args.output_directory, ARM_OUTPUT_SWARM_WINDOWS_TEMPLATE_NO_DIAGNOSTICS).replace(".json", ".parameters.json") )
