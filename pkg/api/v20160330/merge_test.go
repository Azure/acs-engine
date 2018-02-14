package v20160330

import "testing"

func TestMerge(t *testing.T) {
	newCS := &ContainerService{
		Properties: &Properties{
			WindowsProfile: &WindowsProfile{
				AdminUsername: "azureuser",
				AdminPassword: "",
			},
		},
	}

	existingCS := &ContainerService{
		Properties: &Properties{
			WindowsProfile: &WindowsProfile{
				AdminUsername: "azureuser",
				AdminPassword: "existingPassword",
			},
		},
	}
	if err := newCS.Merge(existingCS); err != nil {
		t.Fatalf("unexpectedly detected merge failure, %+v", err)
	}
	if newCS.Properties.WindowsProfile.AdminPassword != "existingPassword" {
		t.Fatalf("unexpected Properties.WindowsProfile.AdminPassword not updated")
	}
}

func TestMergeWithNil(t *testing.T) {
	newCS := &ContainerService{
		Properties: &Properties{},
	}

	existingCS := &ContainerService{
		Properties: &Properties{
			WindowsProfile: &WindowsProfile{
				AdminUsername: "azureuser",
				AdminPassword: "existingPassword",
			},
		},
	}
	if err := newCS.Merge(existingCS); err != nil {
		t.Fatalf("unexpectedly detected merge failure, %+v", err)
	}
	if newCS.Properties.WindowsProfile == nil {
		t.Fatalf("unexpected Properties.WindowsProfile not updated")
	}
	if newCS.Properties.WindowsProfile.AdminUsername != "azureuser" {
		t.Fatalf("unexpected Properties.WindowsProfile.AdminUsername not updated")
	}
	if newCS.Properties.WindowsProfile.AdminPassword != "existingPassword" {
		t.Fatalf("unexpected Properties.WindowsProfile.AdminPassword not updated")
	}
}
