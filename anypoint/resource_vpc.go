package anypoint

import (
	"context"
	"io/ioutil"
	"sort"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	vpc "github.com/mulesoft-anypoint/anypoint-client-go/vpc"
)

func resourceVPC() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVPCCreate,
		ReadContext:   resourceVPCRead,
		UpdateContext: resourceVPCUpdate,
		DeleteContext: resourceVPCDelete,
		Description: `
		Creates a ` + "`" + `vpc` + "`" + `component.
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
				Description: "The unique id of this vpc generated by the anypoint platform.",
			},
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The organization id where the vpc is defined.",
			},
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The name of the vpc.",
			},
			"region": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The CloudHub region where this vpc will exist",
			},
			"cidr_block": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				Description:      "The IP address range that the vpc will use. The largest is /16 and the smallest, /24",
				ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
			},
			"internal_dns_servers": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of internal dns servers",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"internal_dns_special_domains": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "List of internal dns special domains",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"is_default": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "If set to true, the VPC will be associated to all CloudHub environments not explicitly associated to another vpc, including newly created ones",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return DiffSuppressFunc4OptionalPrimitives(k, old, new, d, "false")
				},
			},
			"associated_environments": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of CloudHub environments to associate to this vpc.",
			},
			"owner_id": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The id of the organization that owns the vpc.",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return DiffSuppressFunc4OptionalPrimitives(k, old, new, d, "")
				},
			},
			"shared_with": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Description: "A list of Business Groups to share this vpc with",
			},
			"firewall_rules": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "Inbound firewall rules for all CloudHub workers in this vpc. The list is allow only with an implicit deny all if no rules match",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					return equalsVPCFirewallRules(d.GetChange("firewall_rules"))
				},
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"cidr_block": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
						},
						"protocol": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.StringInSlice([]string{"tcp", "udp"}, true)),
						},
						"from_port": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsPortNumber),
						},
						"to_port": {
							Type:             schema.TypeInt,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsPortNumber),
						},
					},
				},
			},
			"vpc_routes": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "The network routes of this vpc.",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"next_hop": {
							Type:     schema.TypeString,
							Required: true,
						},
						"cidr": {
							Type:             schema.TypeString,
							Required:         true,
							ValidateDiagFunc: validation.ToDiagFunc(validation.IsCIDR),
						},
					},
				},
			},
		},
	}
}

func resourceVPCCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	orgid := d.Get("org_id").(string)

	authctx := getVPCAuthCtx(ctx, &pco)

	body := newVPCBody(d)

	//request vpc creation
	res, httpr, err := pco.vpcclient.DefaultApi.OrganizationsOrgIdVpcsPost(authctx, orgid).VpcCore(*body).Execute()
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
			Summary:  "Unable to Create VPC",
			Detail:   details,
		})
		return diags
	}
	defer httpr.Body.Close()

	d.SetId(res.GetId())

	resourceVPCRead(ctx, d, m)

	return diags
}

func resourceVPCRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	vpcid := d.Id()
	orgid := d.Get("org_id").(string)

	authctx := getVPCAuthCtx(ctx, &pco)

	res, httpr, err := pco.vpcclient.DefaultApi.OrganizationsOrgIdVpcsVpcIdGet(authctx, orgid, vpcid).Execute()
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
			Summary:  "Unable to Get VPC",
			Detail:   details,
		})
		return diags
	}

	//process data
	vpcinstance := flattenVPCData(&res)
	//save in data source schema
	if err := setVPCCoreAttributesToResourceData(d, vpcinstance); err != nil {
		diags := append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to set VPC",
			Detail:   err.Error(),
		})
		return diags
	}

	return diags
}

func resourceVPCUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	vpcid := d.Id()
	orgid := d.Get("org_id").(string)

	authctx := getVPCAuthCtx(ctx, &pco)

	if d.HasChanges(getVPCCoreAttributes()...) {
		body := newVPCBody(d)
		//request vpc creation
		_, httpr, err := pco.vpcclient.DefaultApi.OrganizationsOrgIdVpcsVpcIdPut(authctx, orgid, vpcid).VpcCore(*body).Execute()
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
				Summary:  "Unable to Update VPC",
				Detail:   details,
			})
			return diags
		}
		defer httpr.Body.Close()

		d.Set("last_updated", time.Now().Format(time.RFC850))
	}

	return resourceVPCRead(ctx, d, m)
}

func resourceVPCDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	pco := m.(ProviderConfOutput)
	vpcid := d.Id()
	orgid := d.Get("org_id").(string)

	authctx := getVPCAuthCtx(ctx, &pco)

	httpr, err := pco.vpcclient.DefaultApi.OrganizationsOrgIdVpcsVpcIdDelete(authctx, orgid, vpcid).Execute()
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
			Summary:  "Unable to Delete VPC",
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

/*
 * Creates a new VPC Core Struct from the resource data schema
 */
func newVPCBody(d *schema.ResourceData) *vpc.VpcCore {
	body := vpc.NewVpcCoreWithDefaults()

	body.SetName(d.Get("name").(string))
	body.SetRegion(d.Get("region").(string))
	body.SetCidrBlock(d.Get("cidr_block").(string))
	body.SetIsDefault(d.Get("is_default").(bool))
	body.SetOwnerId(d.Get("owner_id").(string))

	//preparing shared with list
	sw := d.Get("shared_with").([]interface{})
	shared_with := make([]string, len(sw))
	for index, e := range sw {
		shared_with[index] = e.(string)
	}
	body.SetSharedWith(shared_with)

	//preparing associated environments list
	aes := d.Get("associated_environments").([]interface{})
	associated_environments := make([]string, len(aes))
	for index, ae := range aes {
		associated_environments[index] = ae.(string)
	}
	body.SetAssociatedEnvironments(associated_environments)

	//preparing internal_dns structure
	idss := d.Get("internal_dns_servers").([]interface{})
	dns_servers := make([]string, len(idss))
	for index, dns_server := range idss {
		dns_servers[index] = dns_server.(string)
	}
	idsds := d.Get("internal_dns_special_domains").([]interface{})
	special_domains := make([]string, len(idsds))
	for index, special_domain := range idsds {
		special_domains[index] = special_domain.(string)
	}
	body.SetInternalDns(*vpc.NewInternalDns(dns_servers, special_domains))

	//preparing firewall rules
	orules := d.Get("firewall_rules").([]interface{})
	frules := make([]vpc.FirewallRule, len(orules))
	for index, rule := range orules {
		frules[index] = *vpc.NewFirewallRule(rule.(map[string]interface{})["cidr_block"].(string), int32(rule.(map[string]interface{})["from_port"].(int)), rule.(map[string]interface{})["protocol"].(string), int32(rule.(map[string]interface{})["to_port"].(int)))
	}
	body.SetFirewallRules(frules)

	return body
}

/*
 * Returns authentication context (includes authorization header)
 */
func getVPCAuthCtx(ctx context.Context, pco *ProviderConfOutput) context.Context {
	tmp := context.WithValue(ctx, vpc.ContextAccessToken, pco.access_token)
	return context.WithValue(tmp, vpc.ContextServerIndex, pco.server_index)
}

// Compares 2 firewall rules lists
// returns true if they are the same, false otherwise
func equalsVPCFirewallRules(old, new interface{}) bool {
	old_list := old.([]interface{})
	new_list := new.([]interface{})

	if len(new_list) != len(old_list) {
		return false
	}

	if len(new_list) == 0 {
		return true
	}

	sortFirewallRules(old_list)
	sortFirewallRules(new_list)

	for i, val := range old_list {
		o := val.(map[string]interface{})
		n := new_list[i].(map[string]interface{})

		old_cidr_block := o["cidr_block"].(string)
		new_cidr_block := n["cidr_block"].(string)

		if old_cidr_block != new_cidr_block {
			return false
		}

		old_from_port := o["from_port"]
		new_from_port := n["from_port"]

		if old_from_port != new_from_port {
			return false
		}

		old_protocol := o["protocol"]
		new_protocol := n["protocol"]

		if old_protocol != new_protocol {
			return false
		}

		old_to_port := o["to_port"]
		new_to_port := n["to_port"]

		if old_to_port != new_to_port {
			return false
		}
	}

	return true
}

func sortFirewallRules(list []interface{}) {
	sort.SliceStable(list, func(i, j int) bool {
		i_elem := list[i].(map[string]interface{})
		j_elem := list[j].(map[string]interface{})

		i_cidr_block := i_elem["cidr_block"].(string)
		j_cidr_block := j_elem["cidr_block"].(string)

		if i_cidr_block != j_cidr_block {
			return i_cidr_block < j_cidr_block
		}

		i_from_port := i_elem["from_port"].(int)
		j_from_port := j_elem["from_port"].(int)

		if i_from_port != j_from_port {
			return i_from_port < j_from_port
		}

		i_protocol := i_elem["protocol"].(string)
		j_protocol := j_elem["protocol"].(string)

		if i_protocol != j_protocol {
			return i_protocol < j_protocol
		}

		i_to_port := i_elem["to_port"].(int)
		j_to_port := j_elem["to_port"].(int)

		if i_to_port != j_to_port {
			return i_to_port < j_to_port
		}

		return true
	})
}
