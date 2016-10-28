$VerbosePreference="Continue"
$deployName="myKubeVnet"
$RGName=$deployName
$locName="West US"
#$templateFile = "azuredeploy.dcos.json"
#$templateFile = "azuredeploy.swarm.json"
$templateFile = "azuredeploy.kubernetes.json"
New-AzureRmResourceGroup -Name $RGName -Location $locName -Force
New-AzureRmResourceGroupDeployment -Name $deployName -ResourceGroupName $RGName  -TemplateFile $templateFile
