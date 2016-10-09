$VerbosePreference="Continue"
$deployName="anhoweKubeVnet"
$RGName=$deployName
$SubscriptionId="b52fce95-de5f-4b37-afca-db203a5d0b6a"
Set-AzureRmContext -SubscriptionId $SubscriptionId
$locName="West US"
#$templateFile = "azuredeploy.dcos.json"
$templateFile = "azuredeploy.kubernetes.json"
New-AzureRmResourceGroup -Name $RGName -Location $locName -Force
New-AzureRmResourceGroupDeployment -Name $deployName -ResourceGroupName $RGName  -TemplateFile $templateFile
