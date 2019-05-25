package azure

import (
	"fmt"
	"time"

	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/utils"

	"github.com/Azure/go-autorest/autorest/date"

	"github.com/Azure/azure-sdk-for-go/services/web/mgmt/2018-02-01/web"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/suppress"
	"github.com/terraform-providers/terraform-provider-azurerm/azurerm/helpers/validate"
)

func SchemaAppServiceScheduleBackup() *schema.Schema {
	return &schema.Schema{
		Type:     schema.TypeList,
		Optional: true,
		MaxItems: 1,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"frequency_interval": {
					Type:         schema.TypeInt,
					Required:     true,
					ValidateFunc: validateFrequencyInterval,
				},

				"frequency_unit": {
					Type:     schema.TypeString,
					Optional: true,
					Default:  "Day",
					ValidateFunc: validation.StringInSlice([]string{
						"Day",
						"Hour",
					}, false),
				},

				"keep_at_least_one_backup": {
					Type:     schema.TypeBool,
					Optional: true,
					Default:  false,
				},

				"retention_period_in_days": {
					Type:         schema.TypeInt,
					Optional:     true,
					Default:      30,
					ValidateFunc: validateRetentionPeriod,
				},

				"start_time": {
					Type:             schema.TypeString,
					Optional:         true,
					DiffSuppressFunc: suppress.RFC3339Time,
					ValidateFunc:     validate.RFC3339Time,
				},
			},
		},
	}
}

func validateFrequencyInterval(val interface{}, key string) (warns []string, errs []error) {
	v := val.(int)

	if v < 0 || v > 1000 {
		errs = append(errs, fmt.Errorf("%q must be between 0 and 1000 inclusive, got: %d", key, v))
	}
	return
}

func validateRetentionPeriod(val interface{}, key string) (warns []string, errs []error) {
	v := val.(int)

	if v < 0 || v > 9999999 {
		errs = append(errs, fmt.Errorf("%q must be between 0 and 9999999 inclusive, got: %d", key, v))
	}
	return
}

func ExpandAppServiceScheduleBackup(input interface{}) web.BackupSchedule {
	configs := input.([]interface{})
	backupSchedule := web.BackupSchedule{}

	if len(configs) == 0 {
		return backupSchedule
	}

	config := configs[0].(map[string]interface{})

	if v, ok := config["frequency_interval"].(int); ok {
		backupSchedule.FrequencyInterval = utils.Int32(int32(v))
	}

	if v, ok := config["frequency_unit"]; ok {
		backupSchedule.FrequencyUnit = web.FrequencyUnit(v.(string))
	}

	if v, ok := config["keep_at_least_one_backup"]; ok {
		backupSchedule.KeepAtLeastOneBackup = utils.Bool(v.(bool))
	}

	if v, ok := config["retention_period_in_days"].(int); ok {
		backupSchedule.RetentionPeriodInDays = utils.Int32(int32(v))
	}

	if v, ok := config["start_time"].(string); ok {
		dateTimeToStart, _ := time.Parse(time.RFC3339, v) //validated by schema
		backupSchedule.StartTime = &date.Time{Time: (dateTimeToStart)}
	}

	return backupSchedule
}