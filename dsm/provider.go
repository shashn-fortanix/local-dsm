// **********
// Terraform Provider - DSM: provider
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.5.0
//       - Date:      27/07/2021
// **********

package dsm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define Provider
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:     schema.TypeInt,
				Optional: true,
			},
			"insecure": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"username": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"password": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Default:   "",
			},
			"api_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
				Default:   "",
			},
			"acct_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"aws_profile": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "",
			},
			"aws_region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "us-east-1",
			},
			"azure_region": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "us-east",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"dsm_sobject":       resourceSobject(),
			"dsm_aws_sobject":   resourceAWSSobject(),
			"dsm_aws_group":     resourceAWSGroup(),
			"dsm_azure_sobject": resourceAzureSobject(),
			"dsm_azure_group":   resourceAzureGroup(),
			"dsm_secret":        resourceSecret(),
			"dsm_group":         resourceGroup(),
			"dsm_app":           resourceApp(),
			"dsm_gcp_ekm_sa":    resourceGcpEkmSa(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"dsm_aws_group":   dataSourceAWSGroup(),
			"dsm_azure_group": dataSourceAzureGroup(),
			"dsm_secret":      dataSourceSecret(),
			"dsm_group":       dataSourceGroup(),
			"dsm_version":     dataSourceVersion(),
		},
		ConfigureContextFunc: configureProvider,
	}
}

func configureProvider(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	var diags diag.Diagnostics

	// Create new API client
	var bearer bool = true
	if err := d.Get("api_key").(string); len(err) > 0 {
		bearer = false
	}
	newclient, err := NewAPIClient(d.Get("endpoint").(string), d.Get("port").(int), d.Get("username").(string), d.Get("password").(string), d.Get("api_key").(string), bearer, d.Get("acct_id").(string), d.Get("aws_profile").(string), d.Get("aws_region").(string), d.Get("azure_region").(string), d.Get("insecure").(bool))
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK]: Unable to configure DSM provider",
			Detail:   fmt.Sprintf("[E]: SDK: Terraform: %s", err),
		})
		return nil, diags
	}

	return newclient, nil
}
