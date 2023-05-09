Setting up Workload Identity Federation
To exchange a GitHub Actions OIDC token for a Google Cloud access token, you must create and configure a Workload Identity Provider. These instructions use the gcloud command-line tool.

Alternatively, you can also use the gh-oidc Terraform module to automate your infrastructure provisioning. See examples for usage.

Create or use an existing Google Cloud project. You must have privileges to create Workload Identity Pools, Workload Identity Providers, and to manage Service Accounts and IAM permissions. Save your project ID as an environment variable. The rest of these steps assume this environment variable is set:

export PROJECT_ID="my-project" # update with your value
(Optional) Create a Google Cloud Service Account. If you already have a Service Account, take note of the email address and skip this step.

gcloud iam service-accounts create "my-service-account" \
  --project "${PROJECT_ID}"
(Optional) Grant the Google Cloud Service Account permissions to access Google Cloud resources. This step varies by use case. For demonstration purposes, you could grant access to a Google Secret Manager secret or Google Cloud Storage object.

Enable the IAM Credentials API:

gcloud services enable iamcredentials.googleapis.com \
  --project "${PROJECT_ID}"
Create a Workload Identity Pool:

gcloud iam workload-identity-pools create "my-pool" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --display-name="Demo pool"
Get the full ID of the Workload Identity Pool:

gcloud iam workload-identity-pools describe "my-pool" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --format="value(name)"
Save this value as an environment variable:

export WORKLOAD_IDENTITY_POOL_ID="..." # value from above

# This should look like:
#
#   projects/123456789/locations/global/workloadIdentityPools/my-pool
#
Create a Workload Identity Provider in that pool:

gcloud iam workload-identity-pools providers create-oidc "my-provider" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --workload-identity-pool="my-pool" \
  --display-name="Demo provider" \
  --attribute-mapping="google.subject=assertion.sub,attribute.actor=assertion.actor,attribute.repository=assertion.repository" \
  --issuer-uri="https://token.actions.githubusercontent.com"
The attribute mappings map claims in the GitHub Actions JWT to assertions you can make about the request (like the repository or GitHub username of the principal invoking the GitHub Action). These can be used to further restrict the authentication using --attribute-condition flags.

The example above only maps the actor and repository values. To map additional values, add them to the attribute map:

--attribute-mapping="google.subject=assertion.sub,attribute.repository_owner=assertion.repository_owner"
You must map any claims in the incoming token to attributes before you can assert on those attributes in a CEL expression or IAM policy!

Allow authentications from the Workload Identity Provider originating from your repository to impersonate the Service Account created above:

# TODO(developer): Update this value to your GitHub repository.
export REPO="username/name" # e.g. "google/chrome"

gcloud iam service-accounts add-iam-policy-binding "my-service-account@${PROJECT_ID}.iam.gserviceaccount.com" \
  --project="${PROJECT_ID}" \
  --role="roles/iam.workloadIdentityUser" \
  --member="principalSet://iam.googleapis.com/${WORKLOAD_IDENTITY_POOL_ID}/attribute.repository/${REPO}"
If you want to admit all repos of an owner (user or organization), map on attribute.repository_owner:

--member="principalSet://iam.googleapis.com/${WORKLOAD_IDENTITY_POOL_ID}/attribute.repository_owner/${OWNER}"
For this to work, you need to make sure that attribute.repository_owner is mapped in your attribute mapping (see previous step).

Note that $WORKLOAD_IDENTITY_POOL_ID should be the full Workload Identity Pool resource ID, like:

projects/123456789/locations/global/workloadIdentityPools/my-pool
Extract the Workload Identity Provider resource name:

gcloud iam workload-identity-pools providers describe "my-provider" \
  --project="${PROJECT_ID}" \
  --location="global" \
  --workload-identity-pool="my-pool" \
  --format="value(name)"
Use this value as the workload_identity_provider value in your GitHub Actions YAML.

Use this GitHub Action with the Workload Identity Provider ID and Service Account email. The GitHub Action will mint a GitHub OIDC token and exchange the GitHub token for a Google Cloud access token (assuming the authorization is correct). This all happens without exporting a Google Cloud service account key JSON!

Note: It can take up to 5 minutes from when you configure the Workload Identity Pool mapping until the permissions are available.

Organizational Policy Constraints
By default, Google Cloud allows you to create Workload Identity Pools and Workload Identity Providers for any endpoints. Your organization may restrict which external identity providers are permitted on your Google Cloud account. To enable GitHub Actions as a Workload Identity Pool and Provider, add the https://token.actions.githubusercontent.com to the allowed iam.workloadIdentityPoolProviders Org Policy constraint.

gcloud resource-manager org-policies allow "constraints/iam.workloadIdentityPoolProviders" \
  https://token.actions.githubusercontent.com
You can specify a --folder or --organization. If you do not have permission to manage these Org Policies, please contact your Google Cloud administrator.

For GitHub Enterprise Server, the endpoint will be your server URL:

gcloud resource-manager org-policies allow "constraints/iam.workloadIdentityPoolProviders" \
  https://my.github.company
GitHub Token Format
Below is a sample GitHub Token for reference for attribute mappings. For a list of all mappings, see the GitHub OIDC token documentation.

{
  "jti": "...",
  "sub": "repo:username/reponame:ref:refs/heads/main",
  "aud": "https://iam.googleapis.com/projects/123456789/locations/global/workloadIdentityPools/my-pool/providers/my-provider",
  "ref": "refs/heads/main",
  "sha": "d11880f4f451ee35192135525dc974c56a3c1b28",
  "repository": "username/reponame",
  "repository_owner": "username",
  "repository_visibility": "private",
  "repository_id": "74",
  "repository_owner_id": "65",
  "run_id": "1238222155",
  "run_number": "18",
  "run_attempt": "1",
  "actor": "username",
  "actor_id": "12",
  "workflow": "oidc",
  "head_ref": "",
  "base_ref": "",
  "event_name": "push",
  "ref_type": "branch",
  "job_workflow_ref": "username/reponame/.github/workflows/token.yml@refs/heads/main",
  "iss": "https://token.actions.githubusercontent.com",
  "nbf": 1631718827,
  "exp": 1631719727,
  "iat": 1631719427
}
Versioning
We recommend pinning to the latest available major version:

- uses: 'google-github-actions/auth@v1'
While this action attempts to follow semantic versioning, but we're ultimately human and sometimes make mistakes. To prevent accidental breaking changes, you can also pin to a specific version:

- uses: 'google-github-actions/auth@v1.0.0'
However, you will not get automatic security updates or new features without explicitly updating your version number. Note that we only publish MAJOR and MAJOR.MINOR.PATCH versions. There is not a floating alias for MAJOR.MINOR.