# {{PROJECT_NAME}}

# Django Bootstrap Deploy Flow
```
[Azure DevOps Push] 
      ⮕ [Pipeline triggered]
          ⮕ [Trigger OpenShift BuildConfig]
              ⮕ [Wait for build]
                  ⮕ [Save build metadata into Azure Artifact]
                      ⮕ [Optional: Update ArgoCD Git repo with new tag]
                          ⮕ [ArgoCD Auto-Sync (Dev) / Manual Sync (Prod)]
```
