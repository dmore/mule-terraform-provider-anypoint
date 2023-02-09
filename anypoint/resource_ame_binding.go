package anypoint

import (
	"context"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
	"github.com/mulesoft-consulting/anypoint-client-go/ame_binding"
)

func resourceAMEBinding() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceAMEBindingCreate,
		ReadContext:   resourceAMEBindingRead,
		DeleteContext: resourceAMEBindingDelete,
		Description: `
		Creates an ` + "`" + `Anypoint MQ Exchange Binding` + "`" + ` in your ` + "`" + `region` + "`" + `.
		`,
		Schema: map[string]*schema.Schema{
			"last_updated": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The last time this resource has been updated locally.",
			},
			"id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The unique id of this Anypoint MQ Exchange generated by the provider composed of {orgId}_{envId}_{regionId}_{queueId}.",
			},
			"exchange_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique id of this Anypoint MQ Exchange.",
			},
			"queue_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The unique id of this Anypoint MQ Queue.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The organization id where the Anypoint MQ Exchange is defined.",
			},
			"env_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The environment id where the Anypoint MQ Exchange is defined.",
			},
			"region_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The region id where the Anypoint MQ Exchange is defined. Refer to Anypoint Platform official documentation for the list of available regions",
				ValidateDiagFunc: validation.ToDiagFunc(
					validation.StringInSlice(
						[]string{
							"us-east-1", "us-east-2", "us-west-2", "ca-central-1", "eu-west-1", "eu-west-2",
							"ap-southeast-1", "ap-southeast-2", "ap-northeast-1", "eu-central-1",
						},
						false,
					),
				),
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceAMEBindingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)
	envid := d.Get("env_id").(string)
	regionid := d.Get("region_id").(string)
	exchangeid := d.Get("exchange_id").(string)
	queueid := d.Get("queue_id").(string)
	authctx := getAMEBindingAuthCtx(ctx, &pco)

	//request user creation
	_, httpr, err := pco.amebindingclient.DefaultApi.CreateAMEBinding(authctx, orgid, envid, regionid, exchangeid, queueid).Execute()
	if err != nil {
		var details string
		if httpr != nil {
			b, _ := ioutil.ReadAll(httpr.Body)
			details = string(b)
		} else {
			details = err.Error()
		}
		diags := append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create AME Binding " + exchangeid + " " + queueid,
			Detail:   details,
		})
		return diags
	}
	defer httpr.Body.Close()

	d.SetId(ComposeResourceId([]string{orgid, envid, regionid, exchangeid, queueid}))

	return resourceAMEBindingRead(ctx, d, m)
}

func resourceAMEBindingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid, envid, regionid, exchangeid, queueid := DecomposeAMEBindingId(d)
	authctx := getAMEBindingAuthCtx(ctx, &pco)

	//request resource
	_, httpr, err := pco.amebindingclient.DefaultApi.GetAMEBinding(authctx, orgid, envid, regionid, exchangeid, queueid).Execute()
	if err != nil {
		var details string
		if httpr != nil {
			b, _ := ioutil.ReadAll(httpr.Body)
			details = string(b)
		} else {
			details = err.Error()
		}
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get AME Binding " + d.Id(),
			Detail:   details,
		})
		return diags
	}
	defer httpr.Body.Close()

	// setting resource id components for import purposes
	d.Set("org_id", orgid)
	d.Set("env_id", envid)
	d.Set("region_id", regionid)
	d.Set("exchange_id", exchangeid)
	d.Set("queue_id", queueid)

	return diags
}

func resourceAMEBindingDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid, envid, regionid, exchangeid := DecomposeAMEId(d)
	authctx := getAMEBindingAuthCtx(ctx, &pco)

	httpr, err := pco.ameclient.DefaultApi.DeleteAME(authctx, orgid, envid, regionid, exchangeid).Execute()
	if err != nil {
		var details string
		if httpr != nil {
			b, _ := ioutil.ReadAll(httpr.Body)
			details = string(b)
		} else {
			details = err.Error()
		}
		diags := append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to delete AME Binding " + d.Id(),
			Detail:   details,
		})
		return diags
	}
	defer httpr.Body.Close()
	// d.SetId("") is automatically called assuming delete returns no errors, but
	// it is added here for explicitness.
	d.SetId("")

	return diags
}

func DecomposeAMEBindingId(d *schema.ResourceData) (string, string, string, string, string) {
	s := DecomposeResourceId(d.Id())
	return s[0], s[1], s[2], s[3], s[4]
}

// // flattens and maps a given Anypoint MQ Exchange Binding object
// func flattenAMEBindingData(binding *ame_binding.ExchangeBinding) map[string]interface{} {
// 	if binding != nil {
// 		item := make(map[string]interface{})
// 		item["queue_id"] = binding.GetQueueId()
// 		item["exchange_id"] = binding.GetExchangeId()
// 		return item
// 	}

// 	return nil
// }

/*
 * Returns authentication context (includes authorization header)
 */
func getAMEBindingAuthCtx(ctx context.Context, pco *ProviderConfOutput) context.Context {
	tmp := context.WithValue(ctx, ame_binding.ContextAccessToken, pco.access_token)
	return context.WithValue(tmp, ame_binding.ContextServerIndex, pco.server_index)
}
