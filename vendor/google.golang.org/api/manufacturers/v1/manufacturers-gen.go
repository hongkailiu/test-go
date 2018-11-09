// Package manufacturers provides access to the Manufacturer Center API.
//
// See https://developers.google.com/manufacturers/
//
// Usage example:
//
//   import "google.golang.org/api/manufacturers/v1"
//   ...
//   manufacturersService, err := manufacturers.New(oauthHttpClient)
package manufacturers // import "google.golang.org/api/manufacturers/v1"

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	gensupport "google.golang.org/api/gensupport"
	googleapi "google.golang.org/api/googleapi"
)

// Always reference these packages, just in case the auto-generated code
// below doesn't.
var _ = bytes.NewBuffer
var _ = strconv.Itoa
var _ = fmt.Sprintf
var _ = json.NewDecoder
var _ = io.Copy
var _ = url.Parse
var _ = gensupport.MarshalJSON
var _ = googleapi.Version
var _ = errors.New
var _ = strings.Replace
var _ = context.Canceled

const apiId = "manufacturers:v1"
const apiName = "manufacturers"
const apiVersion = "v1"
const basePath = "https://manufacturers.googleapis.com/"

// OAuth2 scopes used by this API.
const (
	// Manage your product listings for Google Manufacturer Center
	ManufacturercenterScope = "https://www.googleapis.com/auth/manufacturercenter"
)

func New(client *http.Client) (*Service, error) {
	if client == nil {
		return nil, errors.New("client is nil")
	}
	s := &Service{client: client, BasePath: basePath}
	s.Accounts = NewAccountsService(s)
	return s, nil
}

type Service struct {
	client    *http.Client
	BasePath  string // API endpoint base URL
	UserAgent string // optional additional User-Agent fragment

	Accounts *AccountsService
}

func (s *Service) userAgent() string {
	if s.UserAgent == "" {
		return googleapi.UserAgent
	}
	return googleapi.UserAgent + " " + s.UserAgent
}

func NewAccountsService(s *Service) *AccountsService {
	rs := &AccountsService{s: s}
	rs.Products = NewAccountsProductsService(s)
	return rs
}

type AccountsService struct {
	s *Service

	Products *AccountsProductsService
}

func NewAccountsProductsService(s *Service) *AccountsProductsService {
	rs := &AccountsProductsService{s: s}
	return rs
}

type AccountsProductsService struct {
	s *Service
}

// Attributes: Attributes of the product. For more information,
// see
// https://support.google.com/manufacturers/answer/6124116.
type Attributes struct {
	// AdditionalImageLink: The additional images of the product. For more
	// information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#addlimage.
	AdditionalImageLink []*Image `json:"additionalImageLink,omitempty"`

	// AgeGroup: The target age group of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#agegroup.
	AgeGroup string `json:"ageGroup,omitempty"`

	// Brand: The brand name of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#brand.
	Brand string `json:"brand,omitempty"`

	// Capacity: The capacity of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#capacity.
	Capacity *Capacity `json:"capacity,omitempty"`

	// Color: The color of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#color.
	Color string `json:"color,omitempty"`

	// Count: The count of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#count.
	Count *Count `json:"count,omitempty"`

	// Description: The description of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#descriptio
	// n.
	Description string `json:"description,omitempty"`

	// DisclosureDate: The disclosure date of the product. For more
	// information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#disclosure
	// .
	DisclosureDate string `json:"disclosureDate,omitempty"`

	// ExcludedDestination: A list of excluded destinations.
	ExcludedDestination []string `json:"excludedDestination,omitempty"`

	// FeatureDescription: The rich format description of the product. For
	// more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#featuredes
	// c.
	FeatureDescription []*FeatureDescription `json:"featureDescription,omitempty"`

	// Flavor: The flavor of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#flavor.
	Flavor string `json:"flavor,omitempty"`

	// Format: The format of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#format.
	Format string `json:"format,omitempty"`

	// Gender: The target gender of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#gender.
	Gender string `json:"gender,omitempty"`

	// Gtin: The Global Trade Item Number (GTIN) of the product. For more
	// information,
	// see https://support.google.com/manufacturers/answer/6124116#gtin.
	Gtin []string `json:"gtin,omitempty"`

	// ImageLink: The image of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#image.
	ImageLink *Image `json:"imageLink,omitempty"`

	// IncludedDestination: A list of included destinations.
	IncludedDestination []string `json:"includedDestination,omitempty"`

	// ItemGroupId: The item group id of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#itemgroupi
	// d.
	ItemGroupId string `json:"itemGroupId,omitempty"`

	// Material: The material of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#material.
	Material string `json:"material,omitempty"`

	// Mpn: The Manufacturer Part Number (MPN) of the product. For more
	// information,
	// see https://support.google.com/manufacturers/answer/6124116#mpn.
	Mpn string `json:"mpn,omitempty"`

	// Pattern: The pattern of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#pattern.
	Pattern string `json:"pattern,omitempty"`

	// ProductDetail: The details of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#productdet
	// ail.
	ProductDetail []*ProductDetail `json:"productDetail,omitempty"`

	// ProductLine: The name of the group of products related to the
	// product. For more
	// information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#productlin
	// e.
	ProductLine string `json:"productLine,omitempty"`

	// ProductName: The canonical name of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#productnam
	// e.
	ProductName string `json:"productName,omitempty"`

	// ProductPageUrl: The URL of the detail page of the product. For more
	// information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#productpag
	// e.
	ProductPageUrl string `json:"productPageUrl,omitempty"`

	// ProductType: The type or category of the product. For more
	// information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#producttyp
	// e.
	ProductType []string `json:"productType,omitempty"`

	// ReleaseDate: The release date of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#release.
	ReleaseDate string `json:"releaseDate,omitempty"`

	// Scent: The scent of the product. For more information, see
	//  https://support.google.com/manufacturers/answer/6124116#scent.
	Scent string `json:"scent,omitempty"`

	// Size: The size of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#size.
	Size string `json:"size,omitempty"`

	// SizeSystem: The size system of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#sizesystem
	// .
	SizeSystem string `json:"sizeSystem,omitempty"`

	// SizeType: The size type of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#sizetype.
	SizeType string `json:"sizeType,omitempty"`

	// SuggestedRetailPrice: The suggested retail price (MSRP) of the
	// product. For more information,
	// see https://support.google.com/manufacturers/answer/6124116#price.
	SuggestedRetailPrice *Price `json:"suggestedRetailPrice,omitempty"`

	// TargetClientId: The target client id. Should only be used in the
	// accounts of the data
	// partners.
	TargetClientId string `json:"targetClientId,omitempty"`

	// Theme: The theme of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#theme.
	Theme string `json:"theme,omitempty"`

	// Title: The title of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#title.
	Title string `json:"title,omitempty"`

	// VideoLink: The videos of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#video.
	VideoLink []string `json:"videoLink,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AdditionalImageLink")
	// to unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AdditionalImageLink") to
	// include in API requests with the JSON null value. By default, fields
	// with empty values are omitted from API requests. However, any field
	// with an empty value appearing in NullFields will be sent to the
	// server as null. It is an error if a field in this list has a
	// non-empty value. This may be used to include null fields in Patch
	// requests.
	NullFields []string `json:"-"`
}

func (s *Attributes) MarshalJSON() ([]byte, error) {
	type NoMethod Attributes
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Capacity: The capacity of a product. For more information,
// see
// https://support.google.com/manufacturers/answer/6124116#capacity.
type Capacity struct {
	// Unit: The unit of the capacity, i.e., MB, GB, or TB.
	Unit string `json:"unit,omitempty"`

	// Value: The numeric value of the capacity.
	Value int64 `json:"value,omitempty,string"`

	// ForceSendFields is a list of field names (e.g. "Unit") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Unit") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Capacity) MarshalJSON() ([]byte, error) {
	type NoMethod Capacity
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Count: The number of products in a single package. For more
// information,
// see
// https://support.google.com/manufacturers/answer/6124116#count.
type Count struct {
	// Unit: The unit in which these products are counted.
	Unit string `json:"unit,omitempty"`

	// Value: The numeric value of the number of products in a package.
	Value int64 `json:"value,omitempty,string"`

	// ForceSendFields is a list of field names (e.g. "Unit") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Unit") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Count) MarshalJSON() ([]byte, error) {
	type NoMethod Count
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// DestinationStatus: The destination status.
type DestinationStatus struct {
	// Destination: The name of the destination.
	Destination string `json:"destination,omitempty"`

	// Status: The status of the destination.
	//
	// Possible values:
	//   "UNKNOWN" - Unspecified status, never used.
	//   "ACTIVE" - The product is used for this destination.
	//   "PENDING" - The decision is still pending.
	//   "DISAPPROVED" - The product is disapproved. Please look at the
	// issues.
	Status string `json:"status,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Destination") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Destination") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *DestinationStatus) MarshalJSON() ([]byte, error) {
	type NoMethod DestinationStatus
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Empty: A generic empty message that you can re-use to avoid defining
// duplicated
// empty messages in your APIs. A typical example is to use it as the
// request
// or the response type of an API method. For instance:
//
//     service Foo {
//       rpc Bar(google.protobuf.Empty) returns
// (google.protobuf.Empty);
//     }
//
// The JSON representation for `Empty` is empty JSON object `{}`.
type Empty struct {
	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`
}

// FeatureDescription: A feature description of the product. For more
// information,
// see
// https://support.google.com/manufacturers/answer/6124116#featuredes
// c.
type FeatureDescription struct {
	// Headline: A short description of the feature.
	Headline string `json:"headline,omitempty"`

	// Image: An optional image describing the feature.
	Image *Image `json:"image,omitempty"`

	// Text: A detailed description of the feature.
	Text string `json:"text,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Headline") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Headline") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *FeatureDescription) MarshalJSON() ([]byte, error) {
	type NoMethod FeatureDescription
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Image: An image.
type Image struct {
	// ImageUrl: The URL of the image. For crawled images, this is the
	// provided URL. For
	// uploaded images, this is a serving URL from Google if the image has
	// been
	// processed successfully.
	ImageUrl string `json:"imageUrl,omitempty"`

	// Status: The status of the image.
	// @OutputOnly
	//
	// Possible values:
	//   "STATUS_UNSPECIFIED" - The image status is unspecified. Should not
	// be used.
	//   "PENDING_PROCESSING" - The image was uploaded and is being
	// processed.
	//   "PENDING_CRAWL" - The image crawl is still pending.
	//   "OK" - The image was processed and it meets the requirements.
	//   "ROBOTED" - The image URL is protected by robots.txt file and
	// cannot be crawled.
	//   "XROBOTED" - The image URL is protected by X-Robots-Tag and cannot
	// be crawled.
	//   "CRAWL_ERROR" - There was an error while crawling the image.
	//   "PROCESSING_ERROR" - The image cannot be processed.
	//   "DECODING_ERROR" - The image cannot be decoded.
	//   "TOO_BIG" - The image is too big.
	//   "CRAWL_SKIPPED" - The image was manually overridden and will not be
	// crawled.
	//   "HOSTLOADED" - The image crawl was postponed to avoid overloading
	// the host.
	//   "HTTP_404" - The image URL returned a "404 Not Found" error.
	Status string `json:"status,omitempty"`

	// Type: The type of the image, i.e., crawled or uploaded.
	// @OutputOnly
	//
	// Possible values:
	//   "TYPE_UNSPECIFIED" - Type is unspecified. Should not be used.
	//   "CRAWLED" - The image was crawled from a provided URL.
	//   "UPLOADED" - The image was uploaded.
	Type string `json:"type,omitempty"`

	// ForceSendFields is a list of field names (e.g. "ImageUrl") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "ImageUrl") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Image) MarshalJSON() ([]byte, error) {
	type NoMethod Image
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Issue: Product issue.
type Issue struct {
	// Attribute: If present, the attribute that triggered the issue. For
	// more information
	// about attributes,
	// see
	// https://support.google.com/manufacturers/answer/6124116.
	Attribute string `json:"attribute,omitempty"`

	// Description: Longer description of the issue focused on how to
	// resolve it.
	Description string `json:"description,omitempty"`

	// Destination: The destination this issue applies to.
	Destination string `json:"destination,omitempty"`

	// Resolution: What needs to happen to resolve the issue.
	//
	// Possible values:
	//   "RESOLUTION_UNSPECIFIED" - Unspecified resolution, never used.
	//   "USER_ACTION" - The user who provided the data must act in order to
	// resolve the issue
	// (for example by correcting some data).
	//   "PENDING_PROCESSING" - The issue will be resolved automatically
	// (for example image crawl or
	// Google review). No action is required now. Resolution might lead
	// to
	// another issue (for example if crawl fails).
	Resolution string `json:"resolution,omitempty"`

	// Severity: The severity of the issue.
	//
	// Possible values:
	//   "SEVERITY_UNSPECIFIED" - Unspecified severity, never used.
	//   "ERROR" - Error severity. The issue prevents the usage of the whole
	// item.
	//   "WARNING" - Warning severity. The issue is either one that prevents
	// the usage of the
	// attribute that triggered it or one that will soon prevent the usage
	// of
	// the whole item.
	//   "INFO" - Info severity. The issue is one that doesn't require
	// immediate attention.
	// It is, for example, used to communicate which attributes are
	// still
	// pending review.
	Severity string `json:"severity,omitempty"`

	// Timestamp: The timestamp when this issue appeared.
	Timestamp string `json:"timestamp,omitempty"`

	// Title: Short title describing the nature of the issue.
	Title string `json:"title,omitempty"`

	// Type: The server-generated type of the issue, for
	// example,
	// “INCORRECT_TEXT_FORMATTING”, “IMAGE_NOT_SERVEABLE”, etc.
	Type string `json:"type,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Attribute") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Attribute") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Issue) MarshalJSON() ([]byte, error) {
	type NoMethod Issue
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

type ListProductsResponse struct {
	// NextPageToken: The token for the retrieval of the next page of
	// product statuses.
	NextPageToken string `json:"nextPageToken,omitempty"`

	// Products: List of the products.
	Products []*Product `json:"products,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "NextPageToken") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "NextPageToken") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ListProductsResponse) MarshalJSON() ([]byte, error) {
	type NoMethod ListProductsResponse
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Price: A price.
type Price struct {
	// Amount: The numeric value of the price.
	Amount string `json:"amount,omitempty"`

	// Currency: The currency in which the price is denoted.
	Currency string `json:"currency,omitempty"`

	// ForceSendFields is a list of field names (e.g. "Amount") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Amount") to include in API
	// requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Price) MarshalJSON() ([]byte, error) {
	type NoMethod Price
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// Product: Product data.
type Product struct {
	// Attributes: Attributes of the product uploaded to the Manufacturer
	// Center. Manually
	// edited attributes are taken into account.
	Attributes *Attributes `json:"attributes,omitempty"`

	// ContentLanguage: The content language of the product as a two-letter
	// ISO 639-1 language code
	// (for example, en).
	ContentLanguage string `json:"contentLanguage,omitempty"`

	// DestinationStatuses: The status of the destinations.
	DestinationStatuses []*DestinationStatus `json:"destinationStatuses,omitempty"`

	// Issues: A server-generated list of issues associated with the
	// product.
	Issues []*Issue `json:"issues,omitempty"`

	// Name: Name in the format
	// `{target_country}:{content_language}:{product_id}`.
	//
	// `target_country`   - The target country of the product as a CLDR
	// territory
	//                      code (for example, US).
	//
	// `content_language` - The content language of the product as a
	// two-letter
	//                      ISO 639-1 language code (for example,
	// en).
	//
	// `product_id`     -   The ID of the product. For more information,
	// see
	//
	// https://support.google.com/manufacturers/answer/6124116#id.
	Name string `json:"name,omitempty"`

	// Parent: Parent ID in the format
	// `accounts/{account_id}`.
	//
	// `account_id` - The ID of the Manufacturer Center account.
	Parent string `json:"parent,omitempty"`

	// ProductId: The ID of the product. For more information,
	// see
	// https://support.google.com/manufacturers/answer/6124116#id.
	ProductId string `json:"productId,omitempty"`

	// TargetCountry: The target country of the product as a CLDR territory
	// code (for example,
	// US).
	TargetCountry string `json:"targetCountry,omitempty"`

	// ServerResponse contains the HTTP response code and headers from the
	// server.
	googleapi.ServerResponse `json:"-"`

	// ForceSendFields is a list of field names (e.g. "Attributes") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "Attributes") to include in
	// API requests with the JSON null value. By default, fields with empty
	// values are omitted from API requests. However, any field with an
	// empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *Product) MarshalJSON() ([]byte, error) {
	type NoMethod Product
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// ProductDetail: A product detail of the product. For more information,
// see
// https://support.google.com/manufacturers/answer/6124116#productdet
// ail.
type ProductDetail struct {
	// AttributeName: The name of the attribute.
	AttributeName string `json:"attributeName,omitempty"`

	// AttributeValue: The value of the attribute.
	AttributeValue string `json:"attributeValue,omitempty"`

	// SectionName: A short section name that can be reused between multiple
	// product details.
	SectionName string `json:"sectionName,omitempty"`

	// ForceSendFields is a list of field names (e.g. "AttributeName") to
	// unconditionally include in API requests. By default, fields with
	// empty values are omitted from API requests. However, any non-pointer,
	// non-interface field appearing in ForceSendFields will be sent to the
	// server regardless of whether the field is empty or not. This may be
	// used to include empty fields in Patch requests.
	ForceSendFields []string `json:"-"`

	// NullFields is a list of field names (e.g. "AttributeName") to include
	// in API requests with the JSON null value. By default, fields with
	// empty values are omitted from API requests. However, any field with
	// an empty value appearing in NullFields will be sent to the server as
	// null. It is an error if a field in this list has a non-empty value.
	// This may be used to include null fields in Patch requests.
	NullFields []string `json:"-"`
}

func (s *ProductDetail) MarshalJSON() ([]byte, error) {
	type NoMethod ProductDetail
	raw := NoMethod(*s)
	return gensupport.MarshalJSON(raw, s.ForceSendFields, s.NullFields)
}

// method id "manufacturers.accounts.products.delete":

type AccountsProductsDeleteCall struct {
	s          *Service
	parent     string
	name       string
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Delete: Deletes the product from a Manufacturer Center account.
func (r *AccountsProductsService) Delete(parent string, name string) *AccountsProductsDeleteCall {
	c := &AccountsProductsDeleteCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	c.name = name
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AccountsProductsDeleteCall) Fields(s ...googleapi.Field) *AccountsProductsDeleteCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AccountsProductsDeleteCall) Context(ctx context.Context) *AccountsProductsDeleteCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AccountsProductsDeleteCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AccountsProductsDeleteCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	c.urlParams_.Set("prettyPrint", "false")
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/{+parent}/products/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, err := http.NewRequest("DELETE", urls, body)
	if err != nil {
		return nil, err
	}
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"parent": c.parent,
		"name":   c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "manufacturers.accounts.products.delete" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *AccountsProductsDeleteCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Deletes the product from a Manufacturer Center account.",
	//   "flatPath": "v1/accounts/{accountsId}/products/{productsId}",
	//   "httpMethod": "DELETE",
	//   "id": "manufacturers.accounts.products.delete",
	//   "parameterOrder": [
	//     "parent",
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "Name in the format `{target_country}:{content_language}:{product_id}`.\n\n`target_country`   - The target country of the product as a CLDR territory\n                     code (for example, US).\n\n`content_language` - The content language of the product as a two-letter\n                     ISO 639-1 language code (for example, en).\n\n`product_id`     -   The ID of the product. For more information, see\n                     https://support.google.com/manufacturers/answer/6124116#id.",
	//       "location": "path",
	//       "pattern": "^[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "parent": {
	//       "description": "Parent ID in the format `accounts/{account_id}`.\n\n`account_id` - The ID of the Manufacturer Center account.",
	//       "location": "path",
	//       "pattern": "^accounts/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/{+parent}/products/{+name}",
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/manufacturercenter"
	//   ]
	// }

}

// method id "manufacturers.accounts.products.get":

type AccountsProductsGetCall struct {
	s            *Service
	parent       string
	name         string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// Get: Gets the product from a Manufacturer Center account, including
// product
// issues.
//
// A recently updated product takes around 15 minutes to process.
// Changes are
// only visible after it has been processed. While some issues may
// be
// available once the product has been processed, other issues may take
// days
// to appear.
func (r *AccountsProductsService) Get(parent string, name string) *AccountsProductsGetCall {
	c := &AccountsProductsGetCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	c.name = name
	return c
}

// Include sets the optional parameter "include": The information to be
// included in the response. Only sections listed here
// will be returned.
//
// Possible values:
//   "UNKNOWN"
//   "ATTRIBUTES"
//   "ISSUES"
//   "DESTINATION_STATUSES"
func (c *AccountsProductsGetCall) Include(include ...string) *AccountsProductsGetCall {
	c.urlParams_.SetMulti("include", append([]string{}, include...))
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AccountsProductsGetCall) Fields(s ...googleapi.Field) *AccountsProductsGetCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AccountsProductsGetCall) IfNoneMatch(entityTag string) *AccountsProductsGetCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AccountsProductsGetCall) Context(ctx context.Context) *AccountsProductsGetCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AccountsProductsGetCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AccountsProductsGetCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	c.urlParams_.Set("prettyPrint", "false")
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/{+parent}/products/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, err := http.NewRequest("GET", urls, body)
	if err != nil {
		return nil, err
	}
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"parent": c.parent,
		"name":   c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "manufacturers.accounts.products.get" call.
// Exactly one of *Product or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Product.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *AccountsProductsGetCall) Do(opts ...googleapi.CallOption) (*Product, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Product{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Gets the product from a Manufacturer Center account, including product\nissues.\n\nA recently updated product takes around 15 minutes to process. Changes are\nonly visible after it has been processed. While some issues may be\navailable once the product has been processed, other issues may take days\nto appear.",
	//   "flatPath": "v1/accounts/{accountsId}/products/{productsId}",
	//   "httpMethod": "GET",
	//   "id": "manufacturers.accounts.products.get",
	//   "parameterOrder": [
	//     "parent",
	//     "name"
	//   ],
	//   "parameters": {
	//     "include": {
	//       "description": "The information to be included in the response. Only sections listed here\nwill be returned.",
	//       "enum": [
	//         "UNKNOWN",
	//         "ATTRIBUTES",
	//         "ISSUES",
	//         "DESTINATION_STATUSES"
	//       ],
	//       "location": "query",
	//       "repeated": true,
	//       "type": "string"
	//     },
	//     "name": {
	//       "description": "Name in the format `{target_country}:{content_language}:{product_id}`.\n\n`target_country`   - The target country of the product as a CLDR territory\n                     code (for example, US).\n\n`content_language` - The content language of the product as a two-letter\n                     ISO 639-1 language code (for example, en).\n\n`product_id`     -   The ID of the product. For more information, see\n                     https://support.google.com/manufacturers/answer/6124116#id.",
	//       "location": "path",
	//       "pattern": "^[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "parent": {
	//       "description": "Parent ID in the format `accounts/{account_id}`.\n\n`account_id` - The ID of the Manufacturer Center account.",
	//       "location": "path",
	//       "pattern": "^accounts/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/{+parent}/products/{+name}",
	//   "response": {
	//     "$ref": "Product"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/manufacturercenter"
	//   ]
	// }

}

// method id "manufacturers.accounts.products.list":

type AccountsProductsListCall struct {
	s            *Service
	parent       string
	urlParams_   gensupport.URLParams
	ifNoneMatch_ string
	ctx_         context.Context
	header_      http.Header
}

// List: Lists all the products in a Manufacturer Center account.
func (r *AccountsProductsService) List(parent string) *AccountsProductsListCall {
	c := &AccountsProductsListCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	return c
}

// Include sets the optional parameter "include": The information to be
// included in the response. Only sections listed here
// will be returned.
//
// Possible values:
//   "UNKNOWN"
//   "ATTRIBUTES"
//   "ISSUES"
//   "DESTINATION_STATUSES"
func (c *AccountsProductsListCall) Include(include ...string) *AccountsProductsListCall {
	c.urlParams_.SetMulti("include", append([]string{}, include...))
	return c
}

// PageSize sets the optional parameter "pageSize": Maximum number of
// product statuses to return in the response, used for
// paging.
func (c *AccountsProductsListCall) PageSize(pageSize int64) *AccountsProductsListCall {
	c.urlParams_.Set("pageSize", fmt.Sprint(pageSize))
	return c
}

// PageToken sets the optional parameter "pageToken": The token returned
// by the previous request.
func (c *AccountsProductsListCall) PageToken(pageToken string) *AccountsProductsListCall {
	c.urlParams_.Set("pageToken", pageToken)
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AccountsProductsListCall) Fields(s ...googleapi.Field) *AccountsProductsListCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// IfNoneMatch sets the optional parameter which makes the operation
// fail if the object's ETag matches the given value. This is useful for
// getting updates only after the object has changed since the last
// request. Use googleapi.IsNotModified to check whether the response
// error from Do is the result of In-None-Match.
func (c *AccountsProductsListCall) IfNoneMatch(entityTag string) *AccountsProductsListCall {
	c.ifNoneMatch_ = entityTag
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AccountsProductsListCall) Context(ctx context.Context) *AccountsProductsListCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AccountsProductsListCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AccountsProductsListCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	if c.ifNoneMatch_ != "" {
		reqHeaders.Set("If-None-Match", c.ifNoneMatch_)
	}
	var body io.Reader = nil
	c.urlParams_.Set("alt", alt)
	c.urlParams_.Set("prettyPrint", "false")
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/{+parent}/products")
	urls += "?" + c.urlParams_.Encode()
	req, err := http.NewRequest("GET", urls, body)
	if err != nil {
		return nil, err
	}
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"parent": c.parent,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "manufacturers.accounts.products.list" call.
// Exactly one of *ListProductsResponse or error will be non-nil. Any
// non-2xx status code is an error. Response headers are in either
// *ListProductsResponse.ServerResponse.Header or (if a response was
// returned at all) in error.(*googleapi.Error).Header. Use
// googleapi.IsNotModified to check whether the returned error was
// because http.StatusNotModified was returned.
func (c *AccountsProductsListCall) Do(opts ...googleapi.CallOption) (*ListProductsResponse, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &ListProductsResponse{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Lists all the products in a Manufacturer Center account.",
	//   "flatPath": "v1/accounts/{accountsId}/products",
	//   "httpMethod": "GET",
	//   "id": "manufacturers.accounts.products.list",
	//   "parameterOrder": [
	//     "parent"
	//   ],
	//   "parameters": {
	//     "include": {
	//       "description": "The information to be included in the response. Only sections listed here\nwill be returned.",
	//       "enum": [
	//         "UNKNOWN",
	//         "ATTRIBUTES",
	//         "ISSUES",
	//         "DESTINATION_STATUSES"
	//       ],
	//       "location": "query",
	//       "repeated": true,
	//       "type": "string"
	//     },
	//     "pageSize": {
	//       "description": "Maximum number of product statuses to return in the response, used for\npaging.",
	//       "format": "int32",
	//       "location": "query",
	//       "type": "integer"
	//     },
	//     "pageToken": {
	//       "description": "The token returned by the previous request.",
	//       "location": "query",
	//       "type": "string"
	//     },
	//     "parent": {
	//       "description": "Parent ID in the format `accounts/{account_id}`.\n\n`account_id` - The ID of the Manufacturer Center account.",
	//       "location": "path",
	//       "pattern": "^accounts/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/{+parent}/products",
	//   "response": {
	//     "$ref": "ListProductsResponse"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/manufacturercenter"
	//   ]
	// }

}

// Pages invokes f for each page of results.
// A non-nil error returned from f will halt the iteration.
// The provided context supersedes any context provided to the Context method.
func (c *AccountsProductsListCall) Pages(ctx context.Context, f func(*ListProductsResponse) error) error {
	c.ctx_ = ctx
	defer c.PageToken(c.urlParams_.Get("pageToken")) // reset paging to original point
	for {
		x, err := c.Do()
		if err != nil {
			return err
		}
		if err := f(x); err != nil {
			return err
		}
		if x.NextPageToken == "" {
			return nil
		}
		c.PageToken(x.NextPageToken)
	}
}

// method id "manufacturers.accounts.products.update":

type AccountsProductsUpdateCall struct {
	s          *Service
	parent     string
	name       string
	attributes *Attributes
	urlParams_ gensupport.URLParams
	ctx_       context.Context
	header_    http.Header
}

// Update: Inserts or updates the attributes of the product in a
// Manufacturer Center
// account.
//
// Creates a product with the provided attributes. If the product
// already
// exists, then all attributes are replaced with the new ones. The
// checks at
// upload time are minimal. All required attributes need to be present
// for a
// product to be valid. Issues may show up later after the API has
// accepted a
// new upload for a product and it is possible to overwrite an existing
// valid
// product with an invalid product. To detect this, you should retrieve
// the
// product and check it for issues once the new version is
// available.
//
// Uploaded attributes first need to be processed before they can
// be
// retrieved. Until then, new products will be unavailable, and
// retrieval
// of previously uploaded products will return the original state of
// the
// product.
func (r *AccountsProductsService) Update(parent string, name string, attributes *Attributes) *AccountsProductsUpdateCall {
	c := &AccountsProductsUpdateCall{s: r.s, urlParams_: make(gensupport.URLParams)}
	c.parent = parent
	c.name = name
	c.attributes = attributes
	return c
}

// Fields allows partial responses to be retrieved. See
// https://developers.google.com/gdata/docs/2.0/basics#PartialResponse
// for more information.
func (c *AccountsProductsUpdateCall) Fields(s ...googleapi.Field) *AccountsProductsUpdateCall {
	c.urlParams_.Set("fields", googleapi.CombineFields(s))
	return c
}

// Context sets the context to be used in this call's Do method. Any
// pending HTTP request will be aborted if the provided context is
// canceled.
func (c *AccountsProductsUpdateCall) Context(ctx context.Context) *AccountsProductsUpdateCall {
	c.ctx_ = ctx
	return c
}

// Header returns an http.Header that can be modified by the caller to
// add HTTP headers to the request.
func (c *AccountsProductsUpdateCall) Header() http.Header {
	if c.header_ == nil {
		c.header_ = make(http.Header)
	}
	return c.header_
}

func (c *AccountsProductsUpdateCall) doRequest(alt string) (*http.Response, error) {
	reqHeaders := make(http.Header)
	for k, v := range c.header_ {
		reqHeaders[k] = v
	}
	reqHeaders.Set("User-Agent", c.s.userAgent())
	var body io.Reader = nil
	body, err := googleapi.WithoutDataWrapper.JSONReader(c.attributes)
	if err != nil {
		return nil, err
	}
	reqHeaders.Set("Content-Type", "application/json")
	c.urlParams_.Set("alt", alt)
	c.urlParams_.Set("prettyPrint", "false")
	urls := googleapi.ResolveRelative(c.s.BasePath, "v1/{+parent}/products/{+name}")
	urls += "?" + c.urlParams_.Encode()
	req, err := http.NewRequest("PUT", urls, body)
	if err != nil {
		return nil, err
	}
	req.Header = reqHeaders
	googleapi.Expand(req.URL, map[string]string{
		"parent": c.parent,
		"name":   c.name,
	})
	return gensupport.SendRequest(c.ctx_, c.s.client, req)
}

// Do executes the "manufacturers.accounts.products.update" call.
// Exactly one of *Empty or error will be non-nil. Any non-2xx status
// code is an error. Response headers are in either
// *Empty.ServerResponse.Header or (if a response was returned at all)
// in error.(*googleapi.Error).Header. Use googleapi.IsNotModified to
// check whether the returned error was because http.StatusNotModified
// was returned.
func (c *AccountsProductsUpdateCall) Do(opts ...googleapi.CallOption) (*Empty, error) {
	gensupport.SetOptions(c.urlParams_, opts...)
	res, err := c.doRequest("json")
	if res != nil && res.StatusCode == http.StatusNotModified {
		if res.Body != nil {
			res.Body.Close()
		}
		return nil, &googleapi.Error{
			Code:   res.StatusCode,
			Header: res.Header,
		}
	}
	if err != nil {
		return nil, err
	}
	defer googleapi.CloseBody(res)
	if err := googleapi.CheckResponse(res); err != nil {
		return nil, err
	}
	ret := &Empty{
		ServerResponse: googleapi.ServerResponse{
			Header:         res.Header,
			HTTPStatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := gensupport.DecodeResponse(target, res); err != nil {
		return nil, err
	}
	return ret, nil
	// {
	//   "description": "Inserts or updates the attributes of the product in a Manufacturer Center\naccount.\n\nCreates a product with the provided attributes. If the product already\nexists, then all attributes are replaced with the new ones. The checks at\nupload time are minimal. All required attributes need to be present for a\nproduct to be valid. Issues may show up later after the API has accepted a\nnew upload for a product and it is possible to overwrite an existing valid\nproduct with an invalid product. To detect this, you should retrieve the\nproduct and check it for issues once the new version is available.\n\nUploaded attributes first need to be processed before they can be\nretrieved. Until then, new products will be unavailable, and retrieval\nof previously uploaded products will return the original state of the\nproduct.",
	//   "flatPath": "v1/accounts/{accountsId}/products/{productsId}",
	//   "httpMethod": "PUT",
	//   "id": "manufacturers.accounts.products.update",
	//   "parameterOrder": [
	//     "parent",
	//     "name"
	//   ],
	//   "parameters": {
	//     "name": {
	//       "description": "Name in the format `{target_country}:{content_language}:{product_id}`.\n\n`target_country`   - The target country of the product as a CLDR territory\n                     code (for example, US).\n\n`content_language` - The content language of the product as a two-letter\n                     ISO 639-1 language code (for example, en).\n\n`product_id`     -   The ID of the product. For more information, see\n                     https://support.google.com/manufacturers/answer/6124116#id.",
	//       "location": "path",
	//       "pattern": "^[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     },
	//     "parent": {
	//       "description": "Parent ID in the format `accounts/{account_id}`.\n\n`account_id` - The ID of the Manufacturer Center account.",
	//       "location": "path",
	//       "pattern": "^accounts/[^/]+$",
	//       "required": true,
	//       "type": "string"
	//     }
	//   },
	//   "path": "v1/{+parent}/products/{+name}",
	//   "request": {
	//     "$ref": "Attributes"
	//   },
	//   "response": {
	//     "$ref": "Empty"
	//   },
	//   "scopes": [
	//     "https://www.googleapis.com/auth/manufacturercenter"
	//   ]
	// }

}
