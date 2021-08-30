// **********
// Terraform Provider - DSM: data source: aws kms group
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.1.5
//       - Date:      05/01/2021
// **********

package dsm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceAWSGroup() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceAWSGroupRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"hmg": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"access_key": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"hmg_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kind": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"tls": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"url": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
			"acct_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"creator": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"scan": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
		},
	}
}

func dataSourceAWSGroupRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, err := m.(*api_client).APICallList("GET", "sys/v1/groups")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: GET sys/v1/groups: %s", err),
		})
		return diags
	}

	for _, data := range req {
		if data.(map[string]interface{})["name"].(string) == d.Get("name").(string) {
			if err := d.Set("name", data.(map[string]interface{})["name"].(string)); err != nil {
				return diag.FromErr(err)
			}
			if err := d.Set("group_id", data.(map[string]interface{})["group_id"].(string)); err != nil {
				return diag.FromErr(err)
			}
			//diags = append(diags, diag.Diagnostic{
			//	Severity: diag.Error,
			//	Summary:  "Unable to call SDKMS provider API client",
			//	Detail:   fmt.Sprintf("[E]: API: GET sys/v1/groups: %s", parseHmg(d, data.(map[string]interface{})["hmg"].(map[string]interface{}))),
			//})
			//return diags
			if err := d.Set("hmg", parseHmg(d, data.(map[string]interface{})["hmg"].(map[string]interface{}))); err != nil {
				return diag.FromErr(err)
			}
			//if err := d.Set("acct_id", data.(map[string]interface{})["acct_id"].(string)); err != nil {
			//	return diag.FromErr(err)
			//}
			if err := d.Set("creator", data.(map[string]interface{})["creator"]); err != nil {
				return diag.FromErr(err)
			}
			if _, ok := data.(map[string]interface{})["description"]; ok {
				if err := d.Set("description", data.(map[string]interface{})["description"].(string)); err != nil {
					return diag.FromErr(err)
				}
			}
		}
	}

	d.SetId(d.Get("group_id").(string))

	// If Scan is set, then move to scanning for data source
	if d.Get("scan").(bool) {
		check_hmg_req := map[string]interface{}{}
		// Scan the AWS Group first before
		_, err := m.(*api_client).APICallBody("POST", fmt.Sprintf("sys/v1/groups/%s/hmg/check", d.Get("group_id").(string)), check_hmg_req)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups/-/hmg/check: %s", err),
			})
			return diags
		}

		_, err = m.(*api_client).APICallBody("POST", fmt.Sprintf("sys/v1/groups/%s/hmg/scan", d.Get("group_id").(string)), check_hmg_req)
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups/-/hmg/scan: %s", err),
			})
			return diags
		}
	}
	return nil
}
