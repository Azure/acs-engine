build-packer:
	@packer build -var-file=packer/settings.json packer/vhd-image.json

init-packer:
	@./packer/init-variables

az-login:
	az login --service-principal -u ${CLIENT_ID} -p ${CLIENT_SECRET} --tenant ${TENANT_ID}

run-packer:
	@packer version && make az-login && make init-packer && (make build-packer | tee packer-output)

az-copy:
	@make az-login && (azcopy --source "${OS_DISK_SAS}" --destination "${CLASSIC_BLOB}/${VHD_NAME}" --dest-sas "${CLASSIC_SAS_TOKEN}")

generate-sas:
	@make az-login && (az storage container generate-sas --name vhds --permissions lr --connection-string "${CLASSIC_SA_CONNECTION_STRING}" --start ${START_DATE} --expiry ${EXPIRY_DATE} | tee vhd-sas)