package ncloud

import (
	"fmt"

	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/ncloud"
	"github.com/NaverCloudPlatform/ncloud-sdk-go-v2/services/server"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceNcloudPublicIp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNcloudPublicIpRead,

		Schema: map[string]*schema.Schema{
			"most_recent": {
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				ForceNew:    true,
				Description: "If more than one result is returned, get the most recent created Public IP.",
			},
			"internet_line_type_code": {
				Type:         schema.TypeString,
				Optional:     true,
				ValidateFunc: validateInternetLineTypeCode,
				Description:  "Internet line type code. `PUBLC` (Public), `GLBL` (Global)",
			},
			"is_associated": {
				Type:        schema.TypeBool,
				Optional:    true,
				Description: "Indicates whether the public IP address is associated or not.",
			},
			"instance_no_list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of public IP instance numbers to get.",
			},
			"list": {
				Type:        schema.TypeList,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Description: "List of public IP addresses to get.",
			},
			"search_filter_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "`publicIp` (Public IP) | `associatedServerName` (Associated server name)",
			},
			"search_filter_value": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Filter value to search",
			},
			"region_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region code. Get available values using the `data ncloud_regions`.",
				ConflictsWith: []string{"region_no"},
			},
			"region_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Region number. Get available values using the `data ncloud_regions`.",
				ConflictsWith: []string{"region_code"},
			},
			"zone_code": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Zone code. You can filter the list of public IP instances by zones. All the public IP addresses in the zone of the region will be selected if the filter is not specified.",
				ConflictsWith: []string{"zone_no"},
			},
			"zone_no": {
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Zone number. You can filter the list of public IP instances by zones. All the public IP addresses in the zone of the region will be selected if the filter is not specified.",
				ConflictsWith: []string{"zone_code"},
			},
			"sorted_by": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The column based on which you want to sort the list.",
			},
			"sorting_order": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Sorting order of the list. `ascending` (Ascending) | `descending` (Descending) [case insensitive]. Default: `ascending` Ascending",
			},

			"instance_no": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP instance number",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP",
			},
			"description": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP description",
			},
			"internet_line_type": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Internet line type",
			},
			"instance_status_name": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Public IP instance status name",
			},
			"instance_status": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Public IP instance status",
			},
			"instance_operation": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Public IP instance operation",
			},
			"kind_type": {
				Type:        schema.TypeMap,
				Computed:    true,
				Elem:        commonCodeSchemaResource,
				Description: "Public IP kind type",
			},
			"server_instance": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "Associated server instance",
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"server_instance_no": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Associated server instance number",
						},
						"server_name": {
							Type:        schema.TypeString,
							Computed:    true,
							Description: "Associated server name",
						},
					},
				},
			},
		},
	}
}

func dataSourceNcloudPublicIpRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*NcloudAPIClient)

	regionNo, err := parseRegionNoParameter(client, d)
	if err != nil {
		return err
	}
	zoneNo, err := parseZoneNoParameter(client, d)
	if err != nil {
		return err
	}

	reqParams := new(server.GetPublicIpInstanceListRequest)

	if internetLineTypeCode, ok := d.GetOk("internet_line_type_code"); ok {
		reqParams.InternetLineTypeCode = ncloud.String(internetLineTypeCode.(string))
	}

	if isAssociated, ok := d.GetOk("is_associated"); ok {
		reqParams.IsAssociated = ncloud.Bool(isAssociated.(bool))
	}

	if instanceNoList, ok := d.GetOk("instance_no_list"); ok {
		reqParams.PublicIpInstanceNoList = expandStringInterfaceList(instanceNoList.([]interface{}))
	}

	if publicIPList, ok := d.GetOk("list"); ok {
		reqParams.PublicIpList = expandStringInterfaceList(publicIPList.([]interface{}))
	}

	if searchFilterName, ok := d.GetOk("search_filter_name"); ok {
		reqParams.SearchFilterName = ncloud.String(searchFilterName.(string))
	}

	if searchFilterValue, ok := d.GetOk("search_filter_value"); ok {
		reqParams.SearchFilterValue = ncloud.String(searchFilterValue.(string))
	}

	reqParams.RegionNo = regionNo
	reqParams.ZoneNo = zoneNo

	if sortedBy, ok := d.GetOk("sorted_by"); ok {
		reqParams.SortedBy = ncloud.String(sortedBy.(string))
	}

	if sortingOrder, ok := d.GetOk("sorting_order"); ok {
		reqParams.SortingOrder = ncloud.String(sortingOrder.(string))
	}

	resp, err := client.server.V2Api.GetPublicIpInstanceList(reqParams)

	if err != nil {
		logErrorResponse("Get Public IP Instance", err, reqParams)
		return err
	}
	publicIpInstanceList := resp.PublicIpInstanceList
	var publicIpInstance *server.PublicIpInstance

	if len(publicIpInstanceList) < 1 {
		return fmt.Errorf("no results. please change search criteria and try again")
	}

	var mostRecent = false
	if _, ok := d.GetOk("most_recent"); ok {
		mostRecent = d.Get("most_recent").(bool)
	}

	if len(publicIpInstanceList) > 1 && mostRecent {
		// Query returned single result.
		publicIpInstance = mostRecentPublicIp(publicIpInstanceList)
	} else {
		publicIpInstance = publicIpInstanceList[0]
	}

	return publicIPAttributes(d, publicIpInstance)
}

func publicIPAttributes(d *schema.ResourceData, instance *server.PublicIpInstance) error {

	d.SetId(ncloud.StringValue(instance.PublicIpInstanceNo))
	d.Set("instance_no", instance.PublicIpInstanceNo)
	d.Set("public_ip", instance.PublicIp)
	d.Set("description", instance.PublicIpDescription)
	d.Set("instance_status_name", instance.PublicIpInstanceStatusName)

	if err := d.Set("internet_line_type", flattenCommonCode(instance.InternetLineType)); err != nil {
		return err
	}
	if err := d.Set("instance_status", flattenCommonCode(instance.PublicIpInstanceStatus)); err != nil {
		return err
	}
	if err := d.Set("instance_operation", flattenCommonCode(instance.PublicIpInstanceOperation)); err != nil {
		return err
	}
	if err := d.Set("kind_type", flattenCommonCode(instance.PublicIpKindType)); err != nil {
		return err
	}

	if serverInstance := instance.ServerInstanceAssociatedWithPublicIp; serverInstance != nil {
		mapping := map[string]interface{}{
			"server_instance_no": ncloud.StringValue(serverInstance.ServerInstanceNo),
			"server_name":        ncloud.StringValue(serverInstance.ServerName),
		}
		d.Set("server_instance", mapping)
	}

	return nil
}
