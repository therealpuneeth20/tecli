/*
Copyright © 2020 Amazon.com, Inc. or its affiliates. All Rights Reserved.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/go-tfe"
	"github.com/spf13/cobra"
	"gitlab.aws.dev/devops-aws/terraform-ce-cli/cobra/aid"
	"gitlab.aws.dev/devops-aws/terraform-ce-cli/cobra/dao"
	"gitlab.aws.dev/devops-aws/terraform-ce-cli/helper"
)

var workspaceValidArgs = []string{
	"list",
	"create",
	"read",
	"read-by-id",
	"update",
	"update-by-id",
	"delete",
	"delete-by-id",
	"remove-vcs-connection",
	"remove-vcs-connection-by-id",
	"lock",
	"unlock",
	"force-unlock",
	"assign-ssh-key",
	"unassign-ssh-key"}

// WorkspaceCmd command to display tecli current version
func WorkspaceCmd() *cobra.Command {
	man, err := helper.GetManualV2("workspace", workspaceValidArgs)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	cmd := &cobra.Command{
		Use:       man.Use,
		Short:     man.Short,
		Long:      man.Long,
		Example:   man.Example,
		ValidArgs: workspaceValidArgs,
		Args:      cobra.OnlyValidArgs,
		PreRunE:   workspacePreRun,
		RunE:      workspaceRun,
	}

	aid.SetWorkspaceFlagsV1(cmd)

	return cmd
}

func workspacePreRun(cmd *cobra.Command, args []string) error {
	if err := helper.ValidateCmdArgs(cmd, args, "workspace"); err != nil {
		return err
	}

	fArg := args[0]
	switch fArg {

	case "list":
		if err := helper.ValidateCmdArgAndFlag(cmd, args, "workspace", fArg, "organization"); err != nil {
			return err
		}

	case "create",
		"read",
		"update",
		"delete",
		"remove-vcs-connection":

		if err := helper.ValidateCmdArgAndFlag(cmd, args, "workspace", fArg, "organization"); err != nil {
			return err
		}

		if err := helper.ValidateCmdArgAndFlag(cmd, args, "workspace", fArg, "name"); err != nil {
			return err
		}

	case "read-by-id",
		"update-by-id",
		"delete-by-id",
		"remove-vcs-connection-by-id",
		"lock",
		"unlock",
		"force-unlock",
		"assign-ssh-key",
		"unassign-ssh-key":
		if err := helper.ValidateCmdArgAndFlag(cmd, args, "workspace", fArg, "id"); err != nil {
			return err
		}

	default:
		return fmt.Errorf("unknown argument")
	}

	return nil
}

func workspaceRun(cmd *cobra.Command, args []string) error {

	token := dao.GetTeamToken(profile)
	client := aid.GetTFEClient(token)

	fArg := args[0]
	switch fArg {
	case "list":
		list, err := workspaceList(client)
		if err == nil {
			if list.TotalCount > 0 {
				for _, item := range list.Items {
					fmt.Printf("%v,\n", aid.ToJSON(item))
				}
			} else {
				return fmt.Errorf("no workspace was found")
			}
		}
	case "create":
		options := aid.GetWorkspaceCreateOptions(cmd)
		workspace, err := workspaceCreate(client, options)

		if err == nil && workspace.ID != "" {
			fmt.Println(aid.ToJSON(workspace))
		} else {
			return fmt.Errorf("unable to create workspace\n%v", err)
		}
	case "read":
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		workspace, err := workspaceRead(client, name)
		if err == nil {
			fmt.Println(aid.ToJSON(workspace))
		} else {
			return fmt.Errorf("workspace %s not found\n%v", name, err)
		}
	case "read-by-id":
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}

		workspace, err := workspaceReadByID(client, id)
		if err == nil {
			fmt.Println(aid.ToJSON(workspace))
		} else {
			return fmt.Errorf("workspace %s not found\n%v", id, err)
		}
	case "update":
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		options := aid.GetWorkspaceUpdateOptions(cmd)
		workspace, err := workspaceUpdate(client, name, options)
		if err == nil && workspace.ID != "" {
			fmt.Println(aid.ToJSON(workspace))
		} else {
			return fmt.Errorf("unable to update workspace\n%v", err)
		}
	case "update-by-id":
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}

		options := aid.GetWorkspaceUpdateOptions(cmd)
		workspace, err := workspaceUpdateByID(client, id, options)
		if err == nil && workspace.ID != "" {
			fmt.Println(aid.ToJSON(workspace))
		} else {
			return fmt.Errorf("unable to update workspace\n%v", err)
		}
	case "delete":
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		err = workspaceDelete(client, name)
		if err == nil {
			fmt.Printf("workspace %s deleted successfully\n", name)
		} else {
			return fmt.Errorf("unable to delete workspace %s\n%v", name, err)
		}
	case "delete-by-id":
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}

		err = workspaceDeleteByID(client, id)
		if err == nil {
			fmt.Printf("workspace %s deleted successfully\n", id)
		} else {
			return fmt.Errorf("unable to delete workspace %s\n%v", id, err)
		}
	case "remove-vcs-connection":
		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		workspace, err := workspaceRemoveVCSConnection(client, name)
		if err == nil {
			fmt.Println(aid.ToJSON(workspace))
		} else {
			return fmt.Errorf("unable to remove vcs connection\n%v", err)
		}
	case "remove-vcs-connection-by-id":
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}

		workspace, err := workspaceRemoveVCSConnectionByID(client, id)
		if err == nil {
			fmt.Println(aid.ToJSON(workspace))
		} else {
			return fmt.Errorf("unable to remove vcs connection\n%v", err)
		}
	case "lock":
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}

		workspace, err := workspaceLock(client, id)
		if err != nil {
			return fmt.Errorf("unable to lock workspace\n%v", err)
		}

		if workspace.Locked {
			fmt.Println(aid.ToJSON(workspace))
		}
	case "unlock":
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}

		workspace, err := workspaceUnlock(client, id)
		if err != nil {
			return err
		}

		if !workspace.Locked {
			fmt.Println(aid.ToJSON(workspace))
		}
	case "force-unlock":
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}

		workspace, err := workspaceForceUnlock(client, id)
		if err != nil {
			return err
		}

		if !workspace.Locked {
			fmt.Println(aid.ToJSON(workspace))
		}
	case "assign-ssh-key":
		id, err := cmd.Flags().GetString("id")
		if err != nil {
			return err
		}

		// TODO: need to fetch the SSH keys via the client.SSHKeys interface
		options := aid.GetWorkspaceAssignSSHKeyOptions(cmd)
		workspace, err := workspaceAssignSSHKey(client, id, options)
		if err != nil {
			return err
		}

		if workspace.ID != "" && workspace.SSHKey.ID != "" {
			fmt.Println(aid.ToJSON(workspace))
		}
	case "unassign-ssh-key":
		fmt.Println("unassign-ssh-key")
	default:
		return fmt.Errorf("unknown argument provided")
	}

	return nil
}

func workspaceList(client *tfe.Client) (*tfe.WorkspaceList, error) {
	return client.Workspaces.List(context.Background(), organization, tfe.WorkspaceListOptions{})
}

// Create is used to create a new workspace.
func workspaceCreate(client *tfe.Client, options tfe.WorkspaceCreateOptions) (*tfe.Workspace, error) {
	return client.Workspaces.Create(context.Background(), organization, options)
}

// Read a workspace by its name.
func workspaceRead(client *tfe.Client, workspace string) (*tfe.Workspace, error) {
	return client.Workspaces.Read(context.Background(), organization, workspace)
}

// Read a workspace by its name.
func workspaceReadByID(client *tfe.Client, workspaceID string) (*tfe.Workspace, error) {
	return client.Workspaces.ReadByID(context.Background(), workspaceID)
}

// Update settings of an existing workspace.
func workspaceUpdate(client *tfe.Client, workspace string, options tfe.WorkspaceUpdateOptions) (*tfe.Workspace, error) {
	return client.Workspaces.Update(context.Background(), organization, workspace, options)
}

// Update settings of an existing workspace.
func workspaceUpdateByID(client *tfe.Client, workspaceID string, options tfe.WorkspaceUpdateOptions) (*tfe.Workspace, error) {
	return client.Workspaces.UpdateByID(context.Background(), workspaceID, options)
}

// // Delete a workspace by its name.
func workspaceDelete(client *tfe.Client, workspace string) error {
	return client.Workspaces.Delete(context.Background(), organization, workspace)
}

// Delete a workspace by its name.
func workspaceDeleteByID(client *tfe.Client, workspaceID string) error {
	return client.Workspaces.DeleteByID(context.Background(), workspaceID)
}

// RemoveVCSConnection from a workspace.
func workspaceRemoveVCSConnection(client *tfe.Client, workspace string) (*tfe.Workspace, error) {
	return client.Workspaces.RemoveVCSConnection(context.Background(), organization, workspace)
}

// RemoveVCSConnection from a workspace.
func workspaceRemoveVCSConnectionByID(client *tfe.Client, workspaceID string) (*tfe.Workspace, error) {
	return client.Workspaces.RemoveVCSConnectionByID(context.Background(), workspaceID)
}

// Lock a workspace by its ID.
func workspaceLock(client *tfe.Client, workspaceID string) (*tfe.Workspace, error) {
	return client.Workspaces.Lock(context.Background(), workspaceID, tfe.WorkspaceLockOptions{})
}

// Unlock a workspace by its ID.
func workspaceUnlock(client *tfe.Client, workspaceID string) (*tfe.Workspace, error) {
	return client.Workspaces.Unlock(context.Background(), workspaceID)
}

// ForceUnlock a workspace by its ID.
func workspaceForceUnlock(client *tfe.Client, workspaceID string) (*tfe.Workspace, error) {
	return client.Workspaces.ForceUnlock(context.Background(), workspaceID)
}

// AssignSSHKey to a workspace.
func workspaceAssignSSHKey(client *tfe.Client, workspaceID string, options tfe.WorkspaceAssignSSHKeyOptions) (*tfe.Workspace, error) {
	return client.Workspaces.AssignSSHKey(context.Background(), workspaceID, options)
}

// UnassignSSHKey from a workspace.
func workspaceUnassignSSHKey(client *tfe.Client, workspaceID string) (*tfe.Workspace, error) {
	return client.Workspaces.UnassignSSHKey(context.Background(), workspaceID)
}
