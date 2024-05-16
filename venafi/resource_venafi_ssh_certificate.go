package venafi

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/Venafi/vcert/v5/pkg/certificate"
	"github.com/Venafi/vcert/v5/pkg/util"
)

func resourceVenafiSshCertificate() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceVenafiSshCertCreate,
		ReadContext:   resourceVenafiSshCertRead,
		DeleteContext: resourceVenafiSshCertDelete,

		Schema: map[string]*schema.Schema{
			"key_id": {
				Type:        schema.TypeString,
				Description: "The identifier of the requested certificate",
				ForceNew:    true,
				Required:    true,
			},
			"template": {
				Type:        schema.TypeString,
				Description: "The certificate issuing template",
				ForceNew:    true,
				Required:    true,
			},
			"key_passphrase": {
				Type:        schema.TypeString,
				Description: "Passphrase for encrypting the private key",
				ForceNew:    true,
				Optional:    true,
				Sensitive:   true,
			},
			"folder": {
				Type:        schema.TypeString,
				Description: "The DN of the policy folder where the certificate object will be created",
				ForceNew:    true,
				Optional:    true,
			},
			"force_command": {
				Type:        schema.TypeString,
				Description: "The requested force command.",
				ForceNew:    true,
				Optional:    true,
			},
			"key_size": {
				Type:        schema.TypeInt,
				Description: "The key size bits, they will be used for creating keypair",
				ForceNew:    true,
				Optional:    true,
			},
			"windows": {
				Type:        schema.TypeBool,
				Description: "If the line endings of service's private key will end on MS windows format",
				ForceNew:    true,
				Optional:    true,
			},
			"valid_hours": {
				Type:        schema.TypeInt,
				Description: "How much time the requester wants to have the certificate valid, the format is hours",
				ForceNew:    true,
				Optional:    true,
			},
			"object_name": {
				Type:        schema.TypeString,
				Description: "The friendly name for the certificate object.",
				ForceNew:    true,
				Optional:    true,
			},
			"public_key": {
				Type:        schema.TypeString,
				Description: "Public key that will be used to generate the certificate",
				Optional:    true,
				ForceNew:    true,
			},
			"certificate": {
				Type:        schema.TypeString,
				Description: "The SSH Certificate",
				Computed:    true,
			},
			"public_key_method": {
				Type:        schema.TypeString,
				Description: "If the public key will be: file provided or local, service generated",
				Optional:    true,
				Default:     "local",
				ForceNew:    true,
			},
			"private_key": {
				Type:        schema.TypeString,
				Description: "Private key",
				Computed:    true,
			},
			"principal": {
				Type:          schema.TypeList,
				Description:   "The requested principals.",
				ForceNew:      true,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"principals"},
				Deprecated:    "This will be removed in the future. Use \"principals\" instead",
			},
			"principals": {
				Type:          schema.TypeList,
				Description:   "The requested principals.",
				ForceNew:      true,
				Optional:      true,
				Elem:          &schema.Schema{Type: schema.TypeString},
				ConflictsWith: []string{"principal"},
			},
			"source_address": {
				Type:        schema.TypeList,
				Description: "The requested source addresses as list of IP/CIDR",
				ForceNew:    true,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"destination_address": {
				Type:        schema.TypeList,
				Description: "The address (FQDN/hostname/IP/CIDR) of the destination host where the certificate will be used to authenticate to",
				ForceNew:    true,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"extension": {
				Type:        schema.TypeList,
				Description: "The requested certificate extensions.",
				ForceNew:    true,
				Optional:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"certificate_type": {
				Type:        schema.TypeString,
				Description: "Certificate type, server or client",
				Computed:    true,
			},
			"public_key_fingerprint": {
				Type:        schema.TypeString,
				Description: "Public key fingerprint SHA256",
				Computed:    true,
			},
			"signing_ca": {
				Type:        schema.TypeString,
				Description: "CA fingerprint SHA256",
				Computed:    true,
			},
			"serial": {
				Type:        schema.TypeString,
				Description: "Serial number",
				Computed:    true,
			},
			"valid_from": {
				Type:        schema.TypeString,
				Description: "Valid from",
				Computed:    true,
			},
			"valid_to": {
				Type:        schema.TypeString,
				Description: "Valid to",
				Computed:    true,
			},
		},
	}
}

func resourceVenafiSshCertCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {

	err := validateSshCertValues(d)

	if err != nil {
		return diag.FromErr(err)
	}

	tflog.Info(ctx, "Creating SSH certificate\n")

	id := d.Get("key_id").(string)
	method := d.Get("public_key_method").(string)
	keyPassphrase := d.Get("key_passphrase").(string)

	cl, err := getConnection(ctx, meta)

	if err != nil {
		return diag.FromErr(err)
	}

	req := buildSshCertRequest(d)

	var privateKey, publicKey []byte
	sPubKey := ""
	//local generated
	if method == "local" {
		keySize := d.Get("key_size").(int)
		if keySize <= 0 {
			keySize = 3072
		}

		privateKey, publicKey, err = util.GenerateSshKeyPair(keySize, keyPassphrase, id)

		if err != nil {
			return diag.FromErr(err)
		}

		sPubKey = string(publicKey)
		req.PublicKeyData = sPubKey
	}
	if method == "file" {
		pubKeyS := d.Get("public_key").(string)

		if pubKeyS == "" {

			return buildStandardDiagError("public key is empty")

		}

		req.PublicKeyData = pubKeyS
	}

	reqData, err := cl.RequestSSHCertificate(&req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(reqData.DN)

	retReq := certificate.SshCertRequest{
		PickupID:                  reqData.DN,
		IncludeCertificateDetails: true,
	}

	if keyPassphrase != "" {
		retReq.PrivateKeyPassphrase = keyPassphrase
	}

	retReq.Timeout = time.Duration(10) * time.Second
	data, err := cl.RetrieveSSHCertificate(&retReq)
	if err != nil {
		return diag.FromErr(fmt.Errorf("failed to retrieve certificate: %s", err))
	}

	//this case is when the keypair is local generated
	if data.PrivateKeyData == "" {
		data.PrivateKeyData = string(privateKey)
	}
	if sPubKey != "" {
		data.PublicKeyData = sPubKey
	}

	if method != "file" {

		err = d.Set("public_key", data.PublicKeyData)
		if err != nil {
			return diag.FromErr(err)
		}

		if data.PrivateKeyData != "" {
			privKeyS := data.PrivateKeyData

			windows := d.Get("windows").(bool)

			if !windows && method == "service" {
				privKeyS = strings.ReplaceAll(privKeyS, "\r\n", "\n")
			}

			err = d.Set("private_key", privKeyS)
			if err != nil {
				return diag.FromErr(err)
			}
		}
	}

	err = d.Set("certificate", data.CertificateData)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("certificate_type", data.CertificateDetails.CertificateType)
	if err != nil {
		return diag.FromErr(err)
	}

	pubKey := fmt.Sprintf("%s:%s", "SHA256", data.CertificateDetails.PublicKeyFingerprintSHA256)
	err = d.Set("public_key_fingerprint", pubKey)
	if err != nil {
		return diag.FromErr(err)
	}

	signingCa := fmt.Sprintf("%s:%s", "SHA256", data.CertificateDetails.CAFingerprintSHA256)
	err = d.Set("signing_ca", signingCa)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("serial", data.CertificateDetails.SerialNumber)
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("valid_from", util.ConvertSecondsToTime(data.CertificateDetails.ValidFrom).String())
	if err != nil {
		return diag.FromErr(err)
	}

	err = d.Set("valid_to", util.ConvertSecondsToTime(data.CertificateDetails.ValidTo).String())
	if err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceVenafiSshCertRead(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	certUntyped, ok := d.GetOk("certificate")
	if !ok {
		d.SetId("")
		return nil
	}

	certStr := certUntyped.(string)
	if certStr == "" {
		d.SetId("")
		return nil
	}
	return nil
}

func resourceVenafiSshCertDelete(_ context.Context, d *schema.ResourceData, _ interface{}) diag.Diagnostics {
	d.SetId("")
	return nil
}
