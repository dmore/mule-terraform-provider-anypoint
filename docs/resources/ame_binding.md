---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "anypoint_ame_binding Resource - terraform-provider-anypoint"
subcategory: ""
description: |-
  Creates an `Anypoint MQ Exchange Binding` in your `region`.
---

# anypoint_ame_binding (Resource)

Creates an `Anypoint MQ Exchange Binding` in your `region`.

## Example Usage

```terraform
resource "anypoint_amq" "amq_01" {
  org_id = var.root_org
  env_id = var.env_id
  region_id = "us-east-1"
  queue_id = "yourQueueID"
  fifo = false
  default_ttl = 604800000
  default_lock_ttl = 120000
  dead_letter_queue_id = "myEmptyDLQ"
  max_deliveries = 10
}

resource "anypoint_amq" "amq_02" {
  org_id = var.root_org
  env_id = var.env_id
  region_id = "us-east-1"
  queue_id = "yourQueueID"
  fifo = false
  default_ttl = 604800000
  default_lock_ttl = 120000
  dead_letter_queue_id = "myEmptyDLQ"
  max_deliveries = 10
}

resource "anypoint_amq" "amq_03" {
  org_id = var.root_org
  env_id = var.env_id
  region_id = "us-east-1"
  queue_id = "yourQueueID"
  fifo = false
  default_ttl = 604800000
  default_lock_ttl = 120000
  dead_letter_queue_id = "myEmptyDLQ"
  max_deliveries = 10
}

resource "anypoint_ame" "ame" {
  org_id = var.root_org
  env_id = var.env_id
  region_id = "us-east-1"
  exchange_id = "myExchangeId"
  encrypted = true
}


resource "anypoint_ame_binding" "ame_b_01" {
  org_id = var.root_org
  env_id = anypoint_amq.ame.env_id
  region_id = anypoint_amq.ame.region_id
  exchange_id = anypoint_ame.ame.exchange_id
  queue_id = anypoint_amq.amq_01.queue_id

  rule_str_compare {
    property_name = "my_property_name"
    property_type = "STRING"
    matcher_type = "EQ"
    value = "full"
  }
}

resource "anypoint_ame_binding" "ame_b_02" {
  org_id = var.root_org
  env_id = anypoint_amq.ame.env_id
  region_id = anypoint_amq.ame.region_id
  exchange_id = anypoint_ame.ame.exchange_id
  queue_id = anypoint_amq.amq_02.queue_id

  rule_str_state {
    property_name = "TO_ROUTE"
    property_type = "STRING"
    matcher_type = "EXISTS"
    value = true
  }
}

resource "anypoint_ame_binding" "ame_b_03" {
  org_id = var.root_org
  env_id = anypoint_amq.ame.env_id
  region_id = anypoint_amq.ame.region_id
  exchange_id = anypoint_ame.ame.exchange_id
  queue_id = anypoint_amq.amq_03.queue_id

  rule_str_set {
    property_name = "horse_name"
    property_type = "STRING"
    matcher_type = "ANY_OF"
    value = tolist(["sugar", "cash", "magic"])
  }

}

resource "anypoint_ame_binding" "ame_b_04" {
  org_id = var.root_org
  env_id = anypoint_amq.ame.env_id
  region_id = anypoint_amq.ame.region_id
  exchange_id = anypoint_ame.ame.exchange_id
  queue_id = anypoint_amq.amq_04.queue_id

  rule_num_compare {
    property_name = "nbr_horses"
    property_type = "NUMERIC"
    matcher_type = "GT"
    value = 12
  }

}

resource "anypoint_ame_binding" "ame_b_05" {
  org_id = var.root_org
  env_id = anypoint_amq.ame.env_id
  region_id = anypoint_amq.ame.region_id
  exchange_id = anypoint_ame.ame.exchange_id
  queue_id = anypoint_amq.amq_05.queue_id

  rule_num_state {
    property_name = "to_ship"
    property_type = "NUMERIC"
    matcher_type = "EXISTS"
    value = true
  }
}


resource "anypoint_ame_binding" "ame_b_06" {
  org_id = var.root_org
  env_id = anypoint_amq.ame.env_id
  region_id = anypoint_amq.ame.region_id
  exchange_id = anypoint_ame.ame.exchange_id
  queue_id = anypoint_amq.amq_06.queue_id

  rule_num_set {
    property_name = "nbr_horses"
    property_type = "NUMERIC"
    matcher_type = "RANGE"
    value = tolist([2,10])
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `env_id` (String) The environment id where the Anypoint MQ Exchange is defined.
- `exchange_id` (String) The unique id of this Anypoint MQ Exchange.
- `org_id` (String) The organization id where the Anypoint MQ Exchange is defined.
- `queue_id` (String) The unique id of this Anypoint MQ Queue.
- `region_id` (String) The region id where the Anypoint MQ Exchange is defined. Refer to Anypoint Platform official documentation for the list of available regions

### Optional

- `last_updated` (String) The last time this resource has been updated locally.
- `rule_num_compare` (Block Set, Max: 1) This rule is to be used when your source attribute is a NUMERIC and you want to compare is to another NUMERIC value (see [below for nested schema](#nestedblock--rule_num_compare))
- `rule_num_set` (Block Set, Max: 1) This rule is to be used when your source attribute is a NUMERIC and you want to check of the property is included or excluded from a set of NUMERIC values (see [below for nested schema](#nestedblock--rule_num_set))
- `rule_num_state` (Block Set, Max: 1) This rule is to be used when your source attribute is a NUMERIC and you want to check the property's existence (see [below for nested schema](#nestedblock--rule_num_state))
- `rule_str_compare` (Block Set, Max: 1) This rule is to be used when your source attribute is a STRING and you want to use EQUAL or PREFIX comparisons (see [below for nested schema](#nestedblock--rule_str_compare))
- `rule_str_set` (Block Set, Max: 1) This rule is to be used when your source attribute is a STRING and you want to check of the property is included or excluded from a set of STRING values (see [below for nested schema](#nestedblock--rule_str_set))
- `rule_str_state` (Block Set, Max: 1) This rule is to be used when your source attribute is a STRING and you want to check the property's existence (see [below for nested schema](#nestedblock--rule_str_state))

### Read-Only

- `id` (String) The unique id of this Anypoint MQ Exchange generated by the provider composed of {orgId}_{envId}_{regionId}_{queueId}.

<a id="nestedblock--rule_num_compare"></a>
### Nested Schema for `rule_num_compare`

Required:

- `matcher_type` (String) The operation to perform on the property.
							Only 'EQ' (equal), 'LT'(less than), 'LE' (less or equal), 'GT' (greater than) and 'GE' (greater or equal)
							values are supported for this specific rule.
- `property_name` (String) The property name subject of the rule
- `property_type` (String) The propety type. Only NUMERIC is supported for this specific rule.
- `value` (Number) The value against which the operation will be performed.


<a id="nestedblock--rule_num_set"></a>
### Nested Schema for `rule_num_set`

Required:

- `matcher_type` (String) The operation to perform on the property. Only 'RANGE' and 'NONE_OF' values are supported for this specific rule
- `property_name` (String) The property name subject of the rule
- `property_type` (String) The propety type. Only NUMERIC is supported for this specific rule.
- `value` (List of Number) The value against which the operation will be performed.


<a id="nestedblock--rule_num_state"></a>
### Nested Schema for `rule_num_state`

Required:

- `matcher_type` (String) The operation to perform on the property. Only 'EXISTS' value is supported for this specific rule
- `property_name` (String) The property name subject of the rule
- `property_type` (String) The propety type. Only NUMERIC is supported for this specific rule.
- `value` (Boolean) The value against which the operation will be performed.


<a id="nestedblock--rule_str_compare"></a>
### Nested Schema for `rule_str_compare`

Required:

- `matcher_type` (String) The operation to perform on the property. Only 'EQ' (equal) and 'PREFIX' values are supported for this specific rule
- `property_name` (String) The property name subject of the rule
- `property_type` (String) The propety type. Only STRING is supported for this specific rule.
- `value` (String) The value against which the operation will be performed.


<a id="nestedblock--rule_str_set"></a>
### Nested Schema for `rule_str_set`

Required:

- `matcher_type` (String) The operation to perform on the property. Only 'ANY_OF' and 'NONE_OF' values are supported for this specific rule
- `property_name` (String) The property name subject of the rule
- `property_type` (String) The propety type. Only STRING is supported for this specific rule.
- `value` (List of String) The value against which the operation will be performed.


<a id="nestedblock--rule_str_state"></a>
### Nested Schema for `rule_str_state`

Required:

- `matcher_type` (String) The operation to perform on the property. Only 'EXISTS' value is supported for this specific rule
- `property_name` (String) The property name subject of the rule
- `property_type` (String) The propety type. Only STRING is supported for this specific rule.
- `value` (Boolean) The value against which the operation will be performed.

## Import

Import is supported using the following syntax:

```shell
# In order for the import to work, you should provide a ID composed of the following:
#  {ORG_ID}/{ENV_ID}/{REGION_ID}/{EXCHANGE_ID}/{QUEUE_ID}

terraform import \
  -var-file params.tfvars.json \    #variable file
  anypoint_ame_binding.ame_b \      #resource name
  aa1f55d6-213d-4f60-845c-201282484cd1/7074fcdd-9b23-4ab6-97r8-5db5f4adf17d/us-east-1/MY-AWESOME-EXCHANGE   #resource id
```