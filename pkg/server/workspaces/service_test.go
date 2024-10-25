// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package workspaces_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	t_targetconfigs "github.com/daytonaio/daytona/internal/testing/provider/targetconfigs"
	t_workspaces "github.com/daytonaio/daytona/internal/testing/server/workspaces"
	"github.com/daytonaio/daytona/internal/testing/server/workspaces/mocks"
	"github.com/daytonaio/daytona/internal/util"
	"github.com/daytonaio/daytona/pkg/apikey"
	"github.com/daytonaio/daytona/pkg/containerregistry"
	"github.com/daytonaio/daytona/pkg/gitprovider"
	"github.com/daytonaio/daytona/pkg/logs"
	"github.com/daytonaio/daytona/pkg/provider"
	"github.com/daytonaio/daytona/pkg/provisioner"
	"github.com/daytonaio/daytona/pkg/server/workspaces"
	"github.com/daytonaio/daytona/pkg/server/workspaces/dto"
	"github.com/daytonaio/daytona/pkg/telemetry"
	"github.com/daytonaio/daytona/pkg/workspace"
	"github.com/daytonaio/daytona/pkg/workspace/project"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const serverApiUrl = "http://localhost:3986"
const serverUrl = "http://localhost:3987"
const serverVersion = "0.0.0-test"
const defaultProjectUser = "daytona"
const defaultProjectImage = "daytonaio/workspace-project:latest"

var targetConfig = provider.TargetConfig{
	Name: "test-target-config",
	ProviderInfo: provider.ProviderInfo{
		Name:    "test-provider",
		Version: "test",
	},
	Options: "test-options",
}
var gitProviderConfigId = "github"

var baseApiUrl = "https://api.github.com"

var gitProviderConfig = gitprovider.GitProviderConfig{
	Id:         "github",
	ProviderId: gitProviderConfigId,
	Alias:      "test-alias",
	Username:   "test-username",
	Token:      "test-token",
	BaseApiUrl: &baseApiUrl,
}

var createWorkspaceDto = dto.CreateWorkspaceDTO{
	Name:         "test",
	Id:           "test",
	TargetConfig: targetConfig.Name,
	Projects: []dto.CreateProjectDTO{
		{
			Name:                "project1",
			GitProviderConfigId: &gitProviderConfig.Id,
			Source: dto.CreateProjectSourceDTO{
				Repository: &gitprovider.GitRepository{
					Id:     "123",
					Url:    "https://github.com/daytonaio/daytona",
					Name:   "daytona",
					Branch: "main",
					Sha:    "sha1",
				},
			},
			Image: util.Pointer(defaultProjectImage),
			User:  util.Pointer(defaultProjectUser),
		},
	},
}

var workspaceInfo = workspace.WorkspaceInfo{
	Name:             createWorkspaceDto.Name,
	ProviderMetadata: "provider-metadata-test",
	Projects: []*project.ProjectInfo{
		{
			Name:             createWorkspaceDto.Projects[0].Name,
			Created:          "1 min ago",
			IsRunning:        true,
			ProviderMetadata: "provider-metadata-test",
			WorkspaceId:      createWorkspaceDto.Id,
		},
	},
}

func TestWorkspaceService(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, telemetry.CLIENT_ID_CONTEXT_KEY, "test")

	workspaceStore := t_workspaces.NewInMemoryWorkspaceStore()

	containerRegistryService := mocks.NewMockContainerRegistryService()

	projectConfigService := mocks.NewMockProjectConfigService()

	targetConfigStore := t_targetconfigs.NewInMemoryTargetConfigStore()
	err := targetConfigStore.Save(&targetConfig)
	require.Nil(t, err)

	apiKeyService := mocks.NewMockApiKeyService()
	gitProviderService := mocks.NewMockGitProviderService()
	mockProvisioner := mocks.NewMockProvisioner()

	wsLogsDir := t.TempDir()
	buildLogsDir := t.TempDir()

	service := workspaces.NewWorkspaceService(workspaces.WorkspaceServiceConfig{
		WorkspaceStore:           workspaceStore,
		TargetConfigStore:        targetConfigStore,
		ServerApiUrl:             serverApiUrl,
		ServerUrl:                serverUrl,
		ServerVersion:            serverVersion,
		ContainerRegistryService: containerRegistryService,
		ProjectConfigService:     projectConfigService,
		DefaultProjectImage:      defaultProjectImage,
		DefaultProjectUser:       defaultProjectUser,
		BuilderImage:             defaultProjectImage,
		ApiKeyService:            apiKeyService,
		Provisioner:              mockProvisioner,
		LoggerFactory:            logs.NewLoggerFactory(&wsLogsDir, &buildLogsDir),
		GitProviderService:       gitProviderService,
	})

	t.Run("CreateWorkspace", func(t *testing.T) {
		var containerRegistry *containerregistry.ContainerRegistry

		containerRegistryService.On("FindByImageName", defaultProjectImage).Return(containerRegistry, containerregistry.ErrContainerRegistryNotFound)

		mockProvisioner.On("CreateWorkspace", mock.Anything, &targetConfig).Return(nil)
		mockProvisioner.On("StartWorkspace", mock.Anything, &targetConfig).Return(nil)

		apiKeyService.On("Generate", apikey.ApiKeyTypeWorkspace, createWorkspaceDto.Id).Return(createWorkspaceDto.Id, nil)
		gitProviderService.On("GetLastCommitSha", createWorkspaceDto.Projects[0].Source.Repository).Return("123", nil)

		for _, project := range createWorkspaceDto.Projects {
			apiKeyService.On("Generate", apikey.ApiKeyTypeProject, fmt.Sprintf("%s/%s", createWorkspaceDto.Id, project.Name)).Return(project.Name, nil)
		}

		proj := &project.Project{
			Name:                createWorkspaceDto.Projects[0].Name,
			Image:               *createWorkspaceDto.Projects[0].Image,
			User:                *createWorkspaceDto.Projects[0].User,
			BuildConfig:         createWorkspaceDto.Projects[0].BuildConfig,
			Repository:          createWorkspaceDto.Projects[0].Source.Repository,
			ApiKey:              createWorkspaceDto.Projects[0].Name,
			GitProviderConfigId: createWorkspaceDto.Projects[0].GitProviderConfigId,
			WorkspaceId:         createWorkspaceDto.Id,
			TargetConfig:        createWorkspaceDto.TargetConfig,
		}

		proj.EnvVars = project.GetProjectEnvVars(proj, project.ProjectEnvVarParams{
			ApiUrl:        serverApiUrl,
			ServerUrl:     serverUrl,
			ServerVersion: serverVersion,
			ClientId:      "test",
		}, false)

		mockProvisioner.On("CreateProject", provisioner.ProjectParams{
			Project:                       proj,
			TargetConfig:                  &targetConfig,
			ContainerRegistry:             containerRegistry,
			GitProviderConfig:             &gitProviderConfig,
			BuilderImage:                  defaultProjectImage,
			BuilderImageContainerRegistry: containerRegistry,
		}).Return(nil)
		mockProvisioner.On("StartProject", provisioner.ProjectParams{
			Project:                       proj,
			TargetConfig:                  &targetConfig,
			ContainerRegistry:             containerRegistry,
			GitProviderConfig:             &gitProviderConfig,
			BuilderImage:                  defaultProjectImage,
			BuilderImageContainerRegistry: containerRegistry,
		}).Return(nil)

		gitProviderService.On("GetConfig", "github").Return(&gitProviderConfig, nil)

		workspace, err := service.CreateWorkspace(ctx, createWorkspaceDto)

		require.Nil(t, err)
		require.NotNil(t, workspace)

		workspaceEquals(t, createWorkspaceDto, workspace, defaultProjectImage)
	})

	t.Run("CreateWorkspace fails when workspace already exists", func(t *testing.T) {
		_, err := service.CreateWorkspace(ctx, createWorkspaceDto)
		require.NotNil(t, err)
		require.Equal(t, workspaces.ErrWorkspaceAlreadyExists, err)
	})

	t.Run("CreateWorkspace fails name validation", func(t *testing.T) {
		invalidWorkspaceRequest := createWorkspaceDto
		invalidWorkspaceRequest.Name = "invalid name"

		_, err := service.CreateWorkspace(ctx, invalidWorkspaceRequest)
		require.NotNil(t, err)
		require.Equal(t, workspaces.ErrInvalidWorkspaceName, err)
	})

	t.Run("GetWorkspace", func(t *testing.T) {
		mockProvisioner.On("GetWorkspaceInfo", mock.Anything, mock.Anything, &targetConfig).Return(&workspaceInfo, nil)

		workspace, err := service.GetWorkspace(ctx, createWorkspaceDto.Id, true)

		require.Nil(t, err)
		require.NotNil(t, workspace)

		workspaceDtoEquals(t, createWorkspaceDto, *workspace, workspaceInfo, defaultProjectImage, true)
	})

	t.Run("GetWorkspace fails when workspace not found", func(t *testing.T) {
		_, err := service.GetWorkspace(ctx, "invalid-id", true)
		require.NotNil(t, err)
		require.Equal(t, workspaces.ErrWorkspaceNotFound, err)
	})

	t.Run("ListWorkspaces", func(t *testing.T) {
		verbose := false
		mockProvisioner.On("GetWorkspaceInfo", mock.Anything, mock.Anything, &targetConfig).Return(&workspaceInfo, nil)

		workspaces, err := service.ListWorkspaces(ctx, verbose)

		require.Nil(t, err)
		require.Len(t, workspaces, 1)

		workspace := workspaces[0]

		workspaceDtoEquals(t, createWorkspaceDto, workspace, workspaceInfo, defaultProjectImage, verbose)
	})

	t.Run("ListWorkspaces - verbose", func(t *testing.T) {
		verbose := true
		mockProvisioner.On("GetWorkspaceInfo", mock.Anything, mock.Anything, &targetConfig).Return(&workspaceInfo, nil)

		workspaces, err := service.ListWorkspaces(ctx, verbose)

		require.Nil(t, err)
		require.Len(t, workspaces, 1)

		workspace := workspaces[0]

		workspaceDtoEquals(t, createWorkspaceDto, workspace, workspaceInfo, defaultProjectImage, verbose)
	})

	t.Run("StartWorkspace", func(t *testing.T) {
		mockProvisioner.On("StartWorkspace", mock.Anything, &targetConfig).Return(nil)
		mockProvisioner.On("StartProject", mock.Anything, &targetConfig).Return(nil)

		err := service.StartWorkspace(ctx, createWorkspaceDto.Id)

		require.Nil(t, err)
	})

	t.Run("StartProject", func(t *testing.T) {
		mockProvisioner.On("StartWorkspace", mock.Anything, &targetConfig).Return(nil)
		mockProvisioner.On("StartProject", mock.Anything, &targetConfig).Return(nil)

		err := service.StartProject(ctx, createWorkspaceDto.Id, createWorkspaceDto.Projects[0].Name)

		require.Nil(t, err)
	})

	t.Run("StopWorkspace", func(t *testing.T) {
		mockProvisioner.On("StopWorkspace", mock.Anything, &targetConfig).Return(nil)
		mockProvisioner.On("StopProject", mock.Anything, &targetConfig).Return(nil)

		err := service.StopWorkspace(ctx, createWorkspaceDto.Id)

		require.Nil(t, err)
	})

	t.Run("StopProject", func(t *testing.T) {
		mockProvisioner.On("StopWorkspace", mock.Anything, &targetConfig).Return(nil)
		mockProvisioner.On("StopProject", mock.Anything, &targetConfig).Return(nil)

		err := service.StopProject(ctx, createWorkspaceDto.Id, createWorkspaceDto.Projects[0].Name)

		require.Nil(t, err)
	})

	t.Run("RemoveWorkspace", func(t *testing.T) {
		mockProvisioner.On("DestroyWorkspace", mock.Anything, &targetConfig).Return(nil)
		mockProvisioner.On("DestroyProject", mock.Anything, &targetConfig).Return(nil)
		apiKeyService.On("Revoke", mock.Anything).Return(nil)

		err := service.RemoveWorkspace(ctx, createWorkspaceDto.Id)

		require.Nil(t, err)

		_, err = service.GetWorkspace(ctx, createWorkspaceDto.Id, true)
		require.Equal(t, workspaces.ErrWorkspaceNotFound, err)
	})

	t.Run("ForceRemoveWorkspace", func(t *testing.T) {
		err := workspaceStore.Save(&workspace.Workspace{Id: createWorkspaceDto.Id, TargetConfig: targetConfig.Name})
		require.Nil(t, err)

		mockProvisioner.On("DestroyWorkspace", mock.Anything, &targetConfig).Return(nil)
		mockProvisioner.On("DestroyProject", mock.Anything, &targetConfig).Return(nil)
		apiKeyService.On("Revoke", mock.Anything).Return(nil)

		err = service.ForceRemoveWorkspace(ctx, createWorkspaceDto.Id)

		require.Nil(t, err)

		_, err = service.GetWorkspace(ctx, createWorkspaceDto.Id, true)
		require.Equal(t, workspaces.ErrWorkspaceNotFound, err)
	})

	t.Run("SetProjectState", func(t *testing.T) {
		ws, err := service.CreateWorkspace(ctx, createWorkspaceDto)
		require.Nil(t, err)

		projectName := ws.Projects[0].Name
		updatedAt := time.Now().Format(time.RFC1123)
		res, err := service.SetProjectState(ws.Id, projectName, &project.ProjectState{
			UpdatedAt: updatedAt,
			Uptime:    10,
			GitStatus: &project.GitStatus{
				CurrentBranch: "main",
			},
		})
		require.Nil(t, err)

		project, err := res.GetProject(projectName)
		require.Nil(t, err)
		require.Equal(t, "main", project.State.GitStatus.CurrentBranch)
	})

	t.Cleanup(func() {
		apiKeyService.AssertExpectations(t)
		mockProvisioner.AssertExpectations(t)
	})
}

func workspaceEquals(t *testing.T, req dto.CreateWorkspaceDTO, workspace *workspace.Workspace, projectImage string) {
	t.Helper()

	require.Equal(t, req.Id, workspace.Id)
	require.Equal(t, req.Name, workspace.Name)
	require.Equal(t, req.TargetConfig, workspace.TargetConfig)

	for i, project := range workspace.Projects {
		require.Equal(t, req.Projects[i].Name, project.Name)
		require.Equal(t, req.Projects[i].Source.Repository.Id, project.Repository.Id)
		require.Equal(t, req.Projects[i].Source.Repository.Url, project.Repository.Url)
		require.Equal(t, req.Projects[i].Source.Repository.Name, project.Repository.Name)
		require.Equal(t, project.ApiKey, project.Name)
		require.Equal(t, project.TargetConfig, req.TargetConfig)
		require.Equal(t, project.Image, projectImage)
	}
}

func workspaceDtoEquals(t *testing.T, req dto.CreateWorkspaceDTO, workspace dto.WorkspaceDTO, workspaceInfo workspace.WorkspaceInfo, projectImage string, verbose bool) {
	t.Helper()

	require.Equal(t, req.Id, workspace.Id)
	require.Equal(t, req.Name, workspace.Name)
	require.Equal(t, req.TargetConfig, workspace.TargetConfig)

	if verbose {
		require.Equal(t, workspace.Info.Name, workspaceInfo.Name)
		require.Equal(t, workspace.Info.ProviderMetadata, workspaceInfo.ProviderMetadata)
	} else {
		require.Nil(t, workspace.Info)
	}

	for i, project := range workspace.Projects {
		require.Equal(t, req.Projects[i].Name, project.Name)
		require.Equal(t, req.Projects[i].Source.Repository.Id, project.Repository.Id)
		require.Equal(t, req.Projects[i].Source.Repository.Url, project.Repository.Url)
		require.Equal(t, req.Projects[i].Source.Repository.Name, project.Repository.Name)
		require.Equal(t, project.ApiKey, project.Name)
		require.Equal(t, project.TargetConfig, req.TargetConfig)
		require.Equal(t, project.Image, projectImage)

		if verbose {
			require.Equal(t, workspace.Info.Projects[i].Name, workspaceInfo.Projects[i].Name)
			require.Equal(t, workspace.Info.Projects[i].Created, workspaceInfo.Projects[i].Created)
			require.Equal(t, workspace.Info.Projects[i].IsRunning, workspaceInfo.Projects[i].IsRunning)
			require.Equal(t, workspace.Info.Projects[i].ProviderMetadata, workspaceInfo.Projects[i].ProviderMetadata)
		}
	}
}
