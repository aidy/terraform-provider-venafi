# Venafi Certificate and Key Policy Specification

The _Venafi Certificate and Key Policy Specification_ is a standard for defining 
constraints and recommendations that govern key generation and certificate issuance. 
The specification is consumable by the VCert CLI and VCert-based integrations like 
the [Venafi Collection for Ansible](https://github.com/Venafi/ansible-collection-venafi) 
and the [Venafi Provider for HashiCorp Terraform](https://github.com/Venafi/terraform-provider-venafi) 
that support _Certificate Policy Management_ for Trust Protection Platform (TPP) and Venafi
as a Service (VaaS).

## Policy-as-Code Structure (JSON)

The structure of the _Venafi Certificate and Key Policy Specification_ is shown 
below and is the same starter policy that can be output by executing the `vcert 
getpolicy --starter` command. The specification has two sections, "policy" and 
"defaults". The "policy" section specifies values with which new certificate requests 
must comply and the "defaults" section specifies values that are recommended for use 
in certificate requests when those values are not specified or overridden. 
VCert also supports YAML formatted input specifications.

```json
{
  "policy": {
    "domains": [ "" ],
    "wildcardAllowed": false,
    "autoInstalled": false,
    "maxValidDays": 0,
    "certificateAuthority": "",
    "subject": {
      "orgs": [ "" ],
      "orgUnits": [ "" ],
      "localities": [ "" ],
      "states": [ "" ],
      "countries": [ "" ]
    },
    "keyPair": {
      "keyTypes": [ "" ],
      "rsaKeySizes": [ 0 ],
      "ellipticCurves": [ "" ],
      "serviceGenerated": false,
      "reuseAllowed": false
    },
    "subjectAltNames": {
      "dnsAllowed": false,
      "ipAllowed": false,
      "emailAllowed": false,
      "uriAllowed": false,
      "upnAllowed": false
    }
  },
  "defaults": {
    "domain": "",
    "subject": {
      "org": "",
      "orgUnits": [ "" ],
      "locality": "",
      "state": "",
      "country": ""
    },
    "keyPair": {
      "keyType": "",
      "rsaKeySize": 0,
      "ellipticCurve": "",
      "serviceGenerated": false
    }
  }
}
```

## Policy-as-Code Parameters

All parameters in a specification are optional thus `{}` is the most simple valid 
specification and results in a policy that uses TPP or VaaS defaults.

| Parameter | Data Type | Description |
| ---- | ---- | ---- |
| `policy`&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; |||
| &emsp;`domains` |string&nbsp;array| Specifies domain suffixes that are permitted in Common Name (CN) and DNS Subject Alternative Name (SAN) values |
| &emsp;`wildcardAllowed` |boolean| Indicates whether CN and DNS SAN values may specify wildcards like "*.example.com" |
| &emsp;`autoInstalled` |boolean| ![TPP Only](https://img.shields.io/badge/TPP%20Only-orange.svg) Indicates whether the requested certificate will be automatically installed (i.e. provisioned) |
| &emsp;`maxValidDays` |integer| Number of days for which the requested certificate will be valid.  May be ignored if the integration with the issuing CA does not support specific end dates. |
| &emsp;`certificateAuthority` |string| **TPP**: the distinguished name of a CA Template object.<br />For example, "\VED\Policy\Certificate Authorites\Entrust Advantage"<br /><br />**VaaS**: CA Account Type ("DIGICERT", "ENTRUST", "GLOBALSIGN", or "BUILTIN"), CA Account Name (as it appears in the web console), and CA Product Type delimited by backslash characters.<br />For example, "DIGICERT\My DigiCert Account\ssl_plus" |
| &emsp;`subject` |||
| &emsp;&emsp;`orgs` | string&nbsp;array | Organization (O) values that are permitted |
| &emsp;&emsp;`orgUnits` | string&nbsp;array | Organizational Unit (OU) values that are permitted |
| &emsp;&emsp;`localities` | string&nbsp;array | City/Locality (L) values that are permitted |
| &emsp;&emsp;`states` | string&nbsp;array | State/Province (ST) values that are permitted |
| &emsp;&emsp;`countries` | string&nbsp;array | [ISO 3166 2-Alpha](https://www.iso.org/obp/ui/#search/code/) Country (C) code values that are permitted |
| &emsp; `keyPair` |||
| &emsp;&emsp;`keyTypes` | string&nbsp;array | Key algorithm: "RSA" and/or _"ECDSA"_ ![TPP Only](https://img.shields.io/badge/TPP%20Only-orange.svg) |
| &emsp;&emsp;`rsaKeySizes` | integer&nbsp;array | Permitted number of bits for RSA keys: 512, 1024, 2048, 3072, and/or 4096 |
| &emsp;&emsp;`ellipticCurves` | string&nbsp;array | ![TPP Only](https://img.shields.io/badge/TPP%20Only-orange.svg) Permitted elliptic curves: "P256", "P384", and/or "P521" |
| &emsp;&emsp;`serviceGenerated` | boolean | ![TPP Only](https://img.shields.io/badge/TPP%20Only-orange.svg) Indicates whether key pair and CSR must be generated by the Venafi machine identity service |
| &emsp;&emsp;`reuseAllowed` | boolean | Indicates whether new certificate requests are permitted to reuse a key pair of a known certificate |
| &emsp;`subjectAltNames` |||
| &emsp;&emsp;`dnsAllowed` | boolean | Indicates whether DNS Subject Alternative Names are permitted|
| &emsp;&emsp;`ipAllowed` | boolean | ![TPP Only](https://img.shields.io/badge/TPP%20Only-orange.svg) Indicates whether IP Address Subject Alternative Names are permitted |
| &emsp;&emsp;`emailAllowed` | boolean | ![TPP Only](https://img.shields.io/badge/TPP%20Only-orange.svg) Indicates whether Email Address (RFC822) Subject Alternative Names are permitted |
| &emsp;&emsp;`uriAllowed` | boolean | ![TPP Only](https://img.shields.io/badge/TPP%20Only-orange.svg) Indicates whether Uniform Resource Indicator (URI) Subject Alternative Names are permitted |
| &emsp;&emsp;`upnAllowed` | boolean | ![TPP Only](https://img.shields.io/badge/TPP%20Only-orange.svg) Indicates whether User Principal Name (UPN) Subject Alternative Names are permitted |
| `defaults` |||
| &emsp;`domain` |string| Domain suffix that should be used by default (e.g. "example.com")|
| &emsp;`subject` |||
| &emsp;&emsp;`org` | string | Organization (O) value that should be used by default (e.g. "Example, Inc.")|
| &emsp;&emsp;`orgUnits` | string&nbsp;array | Organizational Unit (OU) values that should be used by default (e.g. "Quality Assurance")|
| &emsp;&emsp;`locality` | string | City/Locality (L) value that should be used by default (e.g. "Salt Lake City")|
| &emsp;&emsp;`state` | string | State/Province (ST) value that should be used by default (e.g. "Utah")|
| &emsp;&emsp;`country` | string |[ISO 3166 2-Alpha](https://www.iso.org/obp/ui/#search/code/) Country (C) code value that should be used by default (e.g. "US")|
| &emsp;`keyPair` |||
| &emsp;&emsp;`keyType` | string | Key algorithm that should be used by default, "RSA" or _"ECDSA"_ ![TPP Only](https://img.shields.io/badge/TPP%20Only-orange.svg)|
| &emsp;&emsp;`rsaKeySize` | integer | Number of bits that should be used by default for RSA keys: 512, 1024, 2048, 3072, or 4096|
| &emsp;&emsp;`ellipticCurve` | string | ![TPP Only](https://img.shields.io/badge/TPP%20Only-orange.svg) The elliptic curve that should be used by default: "P256", "P384", or "P521"|
| &emsp;&emsp;`serviceGenerated` | boolean | ![TPP Only](https://img.shields.io/badge/TPP%20Only-orange.svg) Indicates whether keys should be generated by the Venafi machine identity service by default|