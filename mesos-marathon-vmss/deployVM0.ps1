$VerbosePreference="Continue"
$deployName="anhowe21h"
$RGName=$deployName
$locName="SouthEast Asia"
#$locName="East US2"
#$locName="West US"
#$locName="Brazil South"
#$locName="Central US"
#$locName="East US"
#$locName="SouthCentral US"
#$locName="Japan East"
#$locName="Japan West"
#$locName="West Europe"
#$locName="North Europe"
#$locName="NorthCentral US"
#$templateFile= "mesos-cluster-with-linux-jumpbox.json"
#$templateFile= "mesos-cluster-with-windows-jumpbox.json"
#$templateFile= "mesos-cluster-with-no-jumpbox.json"
$templateFile= "swarm-cluster-with-no-jumpbox.json"
$templateParameterFile= "cluster.parameters.json"
New-AzureRmResourceGroup -Name $RGName -Location $locName -Force

echo New-AzureRmResourceGroupDeployment -Name $deployName -ResourceGroupName $RGName -TemplateFile $templateFile
New-AzureRmResourceGroupDeployment -Name $deployName -ResourceGroupName $RGName -TemplateParameterFile $templateParameterFile -TemplateFile $templateFile
