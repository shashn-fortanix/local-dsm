package dsm

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

/*
{
	"acct_id": "e8109ee3-a729-4562-8806-a932848191af",
	"approval_policy": {
		"policy": {
			"quorum": {
				"n": 1,
				"members": [{
					"user": "54e489ca-f5aa-4e59-869e-281bbd37caa2"
				}],
				"require_2fa": false,
				"require_password": true
			}
		},
		"manage_groups": false,
		"protect_authentication_methods": true,
		"protect_cryptographic_policy": true,
		"protect_logging_config": true
	}
}
*/

// [-] Define Account Quorum Policy
func resourceAccQuorumPolicy() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceCreateAccQuorumPolicy,
		ReadContext:   resourceReadAccQuorumPolicy,
		UpdateContext: resourceUpdateAccQuorumPolicy,
		DeleteContext: resourceDeleteAccQuorumPolicy,
		Schema: map[string]*schema.Schema{
			"acct_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"approval_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
	}
}

// [U]: Terraform Func: resourceUpdateSobject
func resourceCreateAccQuorumPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	policy_object := map[string]interface{}{
		"acct_id":         d.Get("acct_id").(string),
		"approval_policy": json.RawMessage(d.Get("approval_policy").(string)),
	}

	req, err := m.(*api_client).APICallBody("PATCH", fmt.Sprintf("sys/v1/accounts/%s", policy_object["acct_id"]), policy_object)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "[DSM SDK] Unable to call DSM provider API client",
			Detail:   fmt.Sprintf("[E]: API: POST sys/v1/accounts: %v", err),
		})
		return diags
	}

	d.SetId(req["acct_id"].(string))
	return resourceReadAccQuorumPolicy(ctx, d, m)

}

// [C]: Create Account Quorum Policy
func resourceUpdateAccQuorumPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// [D]: Delete Account Quorum Policy
func resourceDeleteAccQuorumPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return nil
}

// [R]: Read Account Quorum Policy
func resourceReadAccQuorumPolicy(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	req, statuscode, err := m.(*api_client).APICall("GET", fmt.Sprintf("sys/v1/accounts/%s", d.Id()))
	if statuscode == 404 {
		d.SetId("")
	} else {
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "[DSM SDK] Unable to call DSM provider API client",
				Detail:   fmt.Sprintf("[E]: API: GET sys/v1/accounts: %v", err),
			})
			return diags
		}

		if err := d.Set("acct_id", req["acct_id"].(string)); err != nil {
			return diag.FromErr(err)
		}
		if _, ok := req["approval_policy"]; ok {
			if err := d.Set("approval_policy", fmt.Sprintf("%v", req["approval_policy"])); err != nil {
				return diag.FromErr(err)
			}
		}
	}
	return diags
}
