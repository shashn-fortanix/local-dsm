// **********
// Terraform Provider - DSM: resource: group
// **********
//       - Author:    fyoo at fortanix dot com
//       - Version:   0.3.7
//       - Date:      27/11/2020
// **********

package dsm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// [-] Define Group
func resourceGroup() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateGroup,
		ReadContext:   resourceReadGroup,
		UpdateContext: resourceUpdateGroup,
		DeleteContext: resourceDeleteGroup,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"group_id": {
				Type:     schema.TypeString,
				Computed: true,
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
				Optional: true,
				Default:  "",
			},
			"approval_policy": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// [C]: Create Group
func resourceCreateGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	group_object := map[string]interface{}{
		"name":            d.Get("name").(string),
		"description":     d.Get("description").(string),
		"approval_policy": json.RawMessage(d.Get("approval_policy").(string)),
	}

	req, err := m.(*api_client).APICallBody("POST", "sys/v1/groups", group_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/groups: %v", err),
		})
		return diags
	}

	d.SetId(req["group_id"].(string))
	return resourceReadGroup(ctx, d, m)
}

// [R]: Read Group
func resourceReadGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/groups/%s", d.Id()))
	if statuscode == 404 {
		d.SetId("")
	} else {
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: GET sys/v1/groups: %v", err),
			})
			return diags
		}

		if err := d.Set("name", req["name"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("group_id", req["group_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("acct_id", req["acct_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if err := d.Set("creator", req["creator"]); err != nil {
			return diag.FromErr(err)
		}
		if _, ok := req["description"]; ok {
			if err := d.Set("description", req["description"].(string)); err != nil {
				return diag.FromErr(err)
			}
		}
		if _, ok := req["approval_policy"]; ok {
			if err := d.Set("approval_policy", fmt.Sprintf("%v", req["approval_policy"])); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return diags
}

// [U]: Update Group
func resourceUpdateGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// [D]: Delete Group
func resourceDeleteGroup(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	_, statuscode, err := m.(*api_client).APICall("DELETE", fmt.Sprintf("sys/v1/groups/%s", d.Id()))
	if (err != nil) && (statuscode != 404) && (statuscode != 400) {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: DELETE sys/v1/groups: %v", err),
		})
		return diags
	} else {
		if statuscode == 400 {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Call to DSM provider API client failed",
				Detail:   fmt.Sprintf("[E]: API: DELETE sys/v1/groups: %s", "Group Not Empty"),
			})
			return diags
		}
	}

	d.SetId("")
	return nil
}
